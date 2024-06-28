package inittask

import (
	"context"
	"hermes/config"
	"hermes/db/data"
	"hermes/internal/service"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
)

func Init() {
	zaplog.SetLogger(config.GetConfig().LogConfig)
	db := data.GetDataFactory()

	err := service.NewPolicy(db).PolicyEffect(context.Background())
	if err != nil {
		panic(err)
	}
}
