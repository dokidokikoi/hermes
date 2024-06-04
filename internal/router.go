package internal

import (
	"hermes/internal/handler/file"
	"hermes/internal/handler/game"
	"hermes/internal/handler/scraper"

	"github.com/gin-gonic/gin"
)

func Install(r gin.IRouter) {
	gH := game.NewHandler()
	gG := r.Group("/game")
	{
		gG.PUT("", gH.Create)
	}

	sH := scraper.NewHandler()
	sG := r.Group("/scraper")
	{
		sG.POST("", sH.Search)
		sG.POST("/detail", sH.Detail)
		sG.GET("", sH.Get)
	}

	fH := file.NewHandler()
	fG := r.Group("/file")
	{
		fG.GET("/:name", fH.Get)
	}
}
