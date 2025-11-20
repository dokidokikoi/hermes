package inittask

import (
	"context"
	"izumi/config"
	"izumi/db/data"
	"izumi/internal/service"
	systemtask "izumi/internal/system_task"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
)

func Init() {
	config.SetConfig()
	zaplog.SetLogger(config.GetConfig().LogConfig)
	db := data.GetDataFactory()

	err := service.NewPolicy(db).PolicyEffect(context.Background())
	if err != nil {
		panic(err)
	}
	service.StartLibraryCache()

	systemtask.StartDownload()
	systemtask.StartLoad()
}
