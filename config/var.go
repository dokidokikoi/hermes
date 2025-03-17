package config

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	DIR_ENV = "HERMES_DATA_DIR"
)

var (
	DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
	ZhLanguage       = "zh-CN,zh-HK;q=0.9,zh;q=0.8"
	TmpDir           = filepath.Join(os.TempDir(), "hermes")
	Dir              = ""
	DefaultRetryCnt  = 5
)

func init() {
	err := os.Mkdir(TmpDir, os.ModePerm)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			panic(err)
		}
	}
	dir := os.Getenv(DIR_ENV)
	if dir == "" {
		dir, err = os.UserConfigDir()
		if err != nil {
			panic(err)
		}
		dir = filepath.Join(dir, "hermes")
	}
	Dir = dir

	err = os.MkdirAll(Dir, os.ModePerm)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return
		}
		panic(err)
	}
}
