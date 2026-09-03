package inittask

import (
	"context"
	"izumi/config"
	"izumi/db"
	"izumi/internal/service"
	systemtask "izumi/internal/system_task"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"go.uber.org/zap"
)

func Init() {
	config.SetConfig()
	zaplog.SetLogger(config.GetConfig().LogConfig)
	zaplog.L().With(zap.Any("config", config.GetConfig())).Info("")
	db := db.GetStore()

	err := service.NewPolicy(db).PolicyEffect(context.Background())
	if err != nil {
		panic(err)
	}
	service.StartLibraryCache()

	systemtask.StartDownload()
	systemtask.StartLoad()
	systemtask.StartAutoScrap()
	systemtask.StartAutoDetectScrap()
}
