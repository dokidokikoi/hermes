package inittask

import (
	"hermes/config"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
)

func Init() {
	zaplog.SetLogger(config.GetConfig().LogConfig)

}
