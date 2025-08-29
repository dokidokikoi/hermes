package config

import (
	"sync"

	c "github.com/dokidokikoi/go-common/config"
)

var configIns *config

var proxyLock sync.RWMutex

type config struct {
	PGConfig     c.PGConfig     `mapstructure:"postgresql"`
	SqliteConfig c.SqliteConfig `mapstructure:"sqlite"`
	LogConfig    c.LogConfig    `mapstructure:"log"`
	AppConfig    c.AppConfig    `mapstructure:"app"`
}

func SetConfig(path string) {
	c.Parse(path, &configIns, "HERMES")
}

func GetConfig() config {
	return *configIns
}
