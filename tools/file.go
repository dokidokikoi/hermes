package tools

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"hermes/config"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

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
	err = os.Rename(tmpPath, newPath)
	if err != nil {
		return "", err
	}
	return newPath, nil
}

func SaveTmpFile(ext string, data io.Reader) (string, error) {
	return SaveFile(ext, data, config.TmpDir)
}

func SaveBunchTmpFile(fn func(url string) ([]byte, error), urls []string) map[string]string {
	res := map[string]string{}

	wait := sync.WaitGroup{}
	for _, url := range urls {
		url := url
		wait.Add(1)
		gopool.Go(func() {
			defer wait.Done()

			cnt := 0
			var data []byte
			var err error = errors.New("fetch file")
			for err != nil && cnt < 10 {
				cnt++
				data, err = fn(url)
				if err != nil {
					zaplog.L().Error("fetch file error", zap.Int("retry", cnt), zap.String("url", url), zap.Error(err))
				}
			}
			if err != nil {
				zaplog.L().Error("fetch file failed", zap.String("url", url))
			}

			res[url], err = SaveTmpFile(filepath.Ext(url), bytes.NewBuffer(data))
			if err != nil {
				zaplog.L().Error("fetch file failed", zap.String("url", url), zap.Error(err))
			}
		})
	}
	wait.Wait()

	return res
}
