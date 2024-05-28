package main

import (
	"flag"
	"fmt"
	"hermes/config"
	"hermes/inittask"
	"log"

	"github.com/dokidokikoi/go-common/middleware"
	"github.com/dokidokikoi/go-common/shutdown"
	"github.com/gin-gonic/gin"
)

var (
	configFile = flag.String("f", "/data/conf/application.yaml", "the config file")
)

func main() {
	flag.Parse()

	config.SetConfig(*configFile)
	inittask.Init()

	e := gin.New()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	e.Use(middleware.Logger())

	err := e.Run(fmt.Sprintf("%s:%d", config.GetConfig().AppConfig.Host, config.GetConfig().AppConfig.Port))
	if err != nil {
		panic(err)
	}

	shutdown.Close(
		func() {

		},
	)
}
