package config

import c "github.com/dokidokikoi/go-common/config"

var configIns *config

type config struct {
	PGConfig    c.PGConfig  `mapstructure:"postgresql"`
	LogConfig   c.LogConfig `mapstructure:"log"`
	AppConfig   c.AppConfig `mapstructure:"app"`
	ProxyConfig ProxyConfig `mapstructure:"proxy"`
}

func SetConfig(path string) {
	c.Parse(path, &configIns)
}

func GetConfig() config {
	return *configIns
}
