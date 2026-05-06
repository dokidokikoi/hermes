package inittask

import (
	"context"
	"izumi/config"
	"izumi/db/data"
	"izumi/internal/service"
	"izumi/scraper/event"
	systemtask "izumi/internal/system_task"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"go.uber.org/zap"
)

func Init() {
	config.SetConfig()
	zaplog.SetLogger(config.GetConfig().LogConfig)
	zaplog.L().With(zap.Any("config", config.GetConfig())).Info("")
	db := data.GetDataFactory()

	// Initialize Lua plugins
	if err := event.InitLuaPlugins(); err != nil {
		zaplog.L().Warn("failed to initialize Lua plugins", zap.Error(err))
	}

	err := service.NewPolicy(db).PolicyEffect(context.Background())
	if err != nil {
		panic(err)
	}
	service.StartLibraryCache()

	systemtask.StartDownload()
	systemtask.StartLoad()
	systemtask.StartAutoScrap()
}
