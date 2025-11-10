package tools

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"izumi/config"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"go.uber.org/zap"
)

func SaveFile(ext string, data io.Reader, path string) (string, error) {
	tmpPath := filepath.Join(config.TmpDir, fmt.Sprintf("%s_%d%s", time.Now().Format("20060102150405"), rand.Intn(100000), ext))
	f, err := os.Create(tmpPath)
	if err != nil {
		return "", err
	}
	r := io.TeeReader(data, f)
	h := sha256.New()
	_, err = io.Copy(h, r)
	if err != nil {
		return "", err
	}
	newPath := filepath.Join(path, fmt.Sprintf("%X%s", h.Sum(nil), ext))
	_, err = os.Stat(newPath)
	if err == nil {
		return newPath, nil
	}
	err = Move(tmpPath, newPath)
	if err != nil {
		return "", err
	}
	return newPath, nil
}

func Move(source, destination string) error {
	err := os.Rename(source, destination)
	if err != nil && strings.Contains(err.Error(), "invalid cross-device link") {
		return moveCrossDevice(source, destination)
	}
	return err
}

func moveCrossDevice(source, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return errors.Wrap(err, "Open(source)")
	}
	dst, err := os.Create(destination)
	if err != nil {
		src.Close()
		return errors.Wrap(err, "Create(destination)")
	}
	_, err = io.Copy(dst, src)
	src.Close()
	dst.Close()
	if err != nil {
		return errors.Wrap(err, "Copy")
	}
	fi, err := os.Stat(source)
	if err != nil {
		os.Remove(destination)
		return errors.Wrap(err, "Stat")
	}
	err = os.Chmod(destination, fi.Mode())
	if err != nil {
		os.Remove(destination)
		return errors.Wrap(err, "Stat")
	}
	os.Remove(source)
	return nil
}

func SaveTmpFile(ext string, data io.Reader) (string, error) {
	return SaveFile(ext, data, config.TmpDir)
}

func SaveBunchTmpFile(fn func(uri string) ([]byte, error), urls []string) map[string]string {
	res := map[string]string{}

	var lock sync.Mutex
	wait := sync.WaitGroup{}
	for _, uri := range urls {
		uri := uri
		wait.Add(1)
		gopool.Go(func() {
			defer wait.Done()

			u, err := url.Parse(uri)
			if err != nil {
				zaplog.L().Error("parse url error", zap.String("url", uri), zap.Error(err))
				return
			}
			cnt := 0
			var data []byte
			err = errors.New("fetch file")
			for err != nil && cnt < 10 {
				cnt++
				data, err = fn(uri)
				if err != nil {
					zaplog.L().Error("fetch file error", zap.Int("retry", cnt), zap.String("url", uri), zap.Error(err))
				}
			}
			if err != nil {
				zaplog.L().Error("fetch file failed", zap.String("url", uri))
			}

			lock.Lock()
			res[uri], err = SaveTmpFile(filepath.Ext(u.Path), bytes.NewBuffer(data))
			lock.Unlock()
			if err != nil {
				zaplog.L().Error("fetch file failed", zap.String("url", uri), zap.Error(err))
			}
		})
	}
	wait.Wait()

	return res
}

func Cp(src, dst string) error {
	_, err := os.Stat(dst)
	if err == nil {
		return os.ErrExist
	}
	if !os.IsNotExist(err) {
		return err
	}
	dstF, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstF.Close()

	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	_, err = io.Copy(dstF, srcF)
	return err
}
