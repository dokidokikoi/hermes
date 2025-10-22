package config

import (
	"bytes"
	_ "embed"
	"izumi/constant"
	"strings"

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
	c.Parse(bytes.NewBuffer(defaultConfig), &configIns, strings.ToUpper(constant.PROJECT_NAME))
}

func GetConfig() config {
	return *configIns
}
