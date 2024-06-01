package inittask

import (
	"hermes/config"
	"hermes/db/data"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
)

func Init() {
	zaplog.SetLogger(config.GetConfig().LogConfig)
	_ = data.GetDataFactory()
}
