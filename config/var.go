package config

import (
	"errors"
	"izumi/constant"
	"os"
	"path/filepath"
)

const (
	DIR_ENV = "IZUMI_DATA_DIR"
)

var (
	DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36"
	ZhLanguage       = "zh-CN;q=0.8,zh;q=0.7"
	TmpDir           = filepath.Join(os.TempDir(), constant.PROJECT_NAME)
	DataDir          = ""
	DefaultRetryCnt  = 5
	Cfg              *config
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
		dir = filepath.Join(dir, constant.PROJECT_NAME, "data")
	}
	DataDir = dir

	err = os.MkdirAll(DataDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
