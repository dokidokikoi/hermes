package internal

import (
	"hermes/internal/handler/category"
	"hermes/internal/handler/character"
	"hermes/internal/handler/developer"
	"hermes/internal/handler/file"
	"hermes/internal/handler/game"
	"hermes/internal/handler/library"
	"hermes/internal/handler/notice"
	"hermes/internal/handler/person"
	"hermes/internal/handler/policy"
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
		gG.PUT("", middleware.PreHandle(gH.Create))
		gG.POST("/search", middleware.PreHandle(gH.Search))
		gG.GET("/ins", middleware.PreHandle(gH.GetIns))
		gG.PATCH("", middleware.PreHandle(gH.Update))
		gG.GET("/ver", middleware.PreHandle(gH.GetVer))
	}

	sH := scraper.NewHandler()
	sG := r.Group("/scraper")
	{
		sG.POST("", middleware.PreHandle(sH.Search))
		sG.GET("", middleware.PreHandle(sH.Get))
		sG.POST("/scrape", middleware.PreHandle(sH.Scrape))
	}

	fH := file.NewHandler()
	fG := r.Group("/file")
	{
		fG.GET("/:name", fH.Get)
	}

	tH := tag.NewHandler()
	tG := r.Group("/tags")
	{
		tG.GET("", middleware.PreHandle(tH.List))
		tG.POST("", middleware.PreHandle(tH.Create))
		tG.DELETE("", middleware.PreHandle(tH.Del))
		tG.PATCH("", middleware.PreHandle(tH.Update))
	}

	cH := category.NewHandler()
	cG := r.Group("/categories")
	{
		cG.GET("", middleware.PreHandle(cH.List))
		cG.POST("", middleware.PreHandle(cH.Create))
		cG.DELETE("", middleware.PreHandle(cH.Del))
		cG.PATCH("", middleware.PreHandle(cH.Update))
	}

	seriesH := series.NewHandler()
	seriesG := r.Group("/series")
	{
		seriesG.GET("", middleware.PreHandle(seriesH.List))
		seriesG.POST("", middleware.PreHandle(seriesH.Create))
		seriesG.DELETE("", middleware.PreHandle(seriesH.Del))
		seriesG.PATCH("", middleware.PreHandle(seriesH.Update))
	}

	devH := developer.NewHandler()
	devG := r.Group("/developer")
	{
		devG.GET("", middleware.PreHandle(devH.List))
		devG.POST("", middleware.PreHandle(devH.Create))
		devG.DELETE("", middleware.PreHandle(devH.Del))
		devG.PATCH("", middleware.PreHandle(devH.Update))
	}

	characterH := character.NewHandler()
	characterG := r.Group("/character")
	{
		characterG.POST("/search", middleware.PreHandle(characterH.Search))
		characterG.GET("/:id", middleware.PreHandle(characterH.Get))
		characterG.DELETE("", middleware.PreHandle(characterH.Del))
		characterG.PATCH("", middleware.PreHandle(characterH.Update))
	}

	personH := person.NewHandler()
	personG := r.Group("/person")
	{
		personG.POST("/search", middleware.PreHandle(personH.Search))
		personG.POST("", middleware.PreHandle(personH.Upsert))
	}

	policyH := policy.NewHandler()
	policyG := r.Group("/policy")
	{
		policyG.GET("", middleware.PreHandle(policyH.List))
		policyG.PATCH("", middleware.PreHandle(policyH.Update))
	}

	notifyG := r.Group("/notify")
	{
		notifyG.GET("scrap", notice.ServeWs)
	}

	libraryH := library.NewHandler()
	libraryG := r.Group("/library")
	{
		libraryG.GET("", middleware.PreHandle(libraryH.Ls))
	}
}
