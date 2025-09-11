package config

import (
	"bytes"
	_ "embed"

	c "github.com/dokidokikoi/go-common/config"
)

//go:embed application.yaml
var defaultConfig []byte

var configIns *config

type config struct {
	PGConfig     c.PGConfig     `mapstructure:"postgresql"`
	SqliteConfig c.SqliteConfig `mapstructure:"sqlite"`
	LogConfig    c.LogConfig    `mapstructure:"log"`
	AppConfig    c.AppConfig    `mapstructure:"app"`
}

func SetConfig() {
	c.Parse(bytes.NewBuffer(defaultConfig), &configIns, "HERMES")
}

func GetConfig() config {
	return *configIns
}
