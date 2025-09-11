package main

import (
	"context"
	"errors"
	"fmt"
	"hermes/config"
	"hermes/db/data"
	"hermes/internal"
	"hermes/internal/inittask"
	"log"
	"net/http"
	"time"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/dokidokikoi/go-common/middleware"
	"github.com/dokidokikoi/go-common/shutdown"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	inittask.Init()

	zaplog.L().Debug("config", zap.Any("config", config.GetConfig()))

	e := gin.New()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Cors())

	internal.Install(e)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.GetConfig().AppConfig.Host, config.GetConfig().AppConfig.Port),
		Handler: e,
	}

	go func() {
		log.Printf("serve on :%d\n", config.GetConfig().AppConfig.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zaplog.L().Panic("server panic", zap.Error(err))
		}
	}()

	shutdown.Close(
		// 关闭服务器
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				zaplog.L().Error("server shutdown error", zap.Error(err))
			}
		},
		// 关闭数据库连接
		func() {
			if err := data.Close(); err != nil {
				zaplog.L().Error("db close error", zap.Error(err))
			}
		},
	)
}
