package config

import (
	"sync"

	c "github.com/dokidokikoi/go-common/config"
)

var configIns *config

var proxyLock sync.RWMutex

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

func GetProxyConfig() ProxyConfig {
	proxyLock.RLock()
	defer proxyLock.RUnlock()

	return configIns.ProxyConfig
}
func SetProxyConfig(p ProxyConfig) {
	proxyLock.Lock()
	defer proxyLock.Unlock()

	configIns.ProxyConfig = p
}
