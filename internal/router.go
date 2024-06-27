package internal

import (
	"hermes/internal/handler/category"
	"hermes/internal/handler/character"
	"hermes/internal/handler/developer"
	"hermes/internal/handler/file"
	"hermes/internal/handler/game"
	"hermes/internal/handler/person"
	"hermes/internal/handler/publisher"
	"hermes/internal/handler/scraper"
	"hermes/internal/handler/series"
	"hermes/internal/handler/tag"

	"github.com/dokidokikoi/go-common/middleware"
	"github.com/gin-gonic/gin"
)

func Install(r gin.IRouter) {
	r.Use(middleware.Cors())

	gH := game.NewHandler()
	gG := r.Group("/game")
	{
		gG.PUT("", gH.Create)
		gG.POST("/search", gH.Search)
		gG.GET("/:id", gH.Get)
		gG.PATCH("", gH.Update)
	}

	sH := scraper.NewHandler()
	sG := r.Group("/scraper")
	{
		sG.POST("", sH.Search)
		sG.GET("", sH.Get)
		sG.GET("/ref", sH.GetRef)
		sG.POST("/scrape", sH.Scrape)
		sG.POST("/bind", sH.Bind)
	}

	fH := file.NewHandler()
	fG := r.Group("/file")
	{
		fG.GET("/:name", fH.Get)
	}

	tH := tag.NewHandler()
	tG := r.Group("/tags")
	{
		tG.GET("", tH.List)
		tG.POST("", tH.Create)
	}

	cH := category.NewHandler()
	cG := r.Group("/categories")
	{
		cG.GET("", cH.List)
		cG.POST("", cH.Create)
	}

	seriesH := series.NewHandler()
	seriesG := r.Group("/series")
	{
		seriesG.GET("", seriesH.List)
		seriesG.POST("", seriesH.Create)
	}

	devH := developer.NewHandler()
	devG := r.Group("/developer")
	{
		devG.GET("", devH.List)
		devG.POST("", devH.Create)
	}

	pubH := publisher.NewHandler()
	pubG := r.Group("/publisher")
	{
		pubG.GET("", pubH.List)
		pubG.POST("", pubH.Create)
	}

	characterH := character.NewHandler()
	characterG := r.Group("/character")
	{
		characterG.POST("/search", characterH.Search)
	}

	personH := person.NewHandler()
	personG := r.Group("/person")
	{
		personG.POST("/search", personH.Search)
		personG.POST("", personH.Upsert)
	}
}
