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

type PluginConfig struct {
	Lua LuaConfig `mapstructure:"lua"`
}

type LuaConfig struct {
	Enabled           bool   `mapstructure:"enabled"`
	ScriptDir         string `mapstructure:"script_dir"`
	VMPoolSize        int    `mapstructure:"vm_pool_size"`
	HotReload         bool   `mapstructure:"hot_reload"`
	ExecutionTimeout  string `mapstructure:"execution_timeout"`
}

type config struct {
	PGConfig     c.PGConfig     `mapstructure:"postgresql"`
	SqliteConfig c.SqliteConfig `mapstructure:"sqlite"`
	LogConfig    c.LogConfig    `mapstructure:"log"`
	AppConfig    c.AppConfig    `mapstructure:"app"`
	PluginConfig PluginConfig   `mapstructure:"plugin"`
}

func SetConfig() {
	c.Parse(bytes.NewBuffer(defaultConfig), &configIns, strings.ToUpper(constant.PROJECT_NAME))
	Cfg = configIns
}

func GetConfig() config {
	return *configIns
}
