package utils

import (
	"bytes"
	"io"
	"izumi/config"
	"net/url"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"

	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/dokidokikoi/go-common/tools"
	"go.uber.org/zap"
)

func SaveTmpFile(ext string, data io.Reader) (string, error) {
	return tools.SaveFileWithMd5Name(data, config.TmpDir, ext)
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
