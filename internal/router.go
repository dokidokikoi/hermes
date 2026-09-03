package internal

import (
	"izumi/internal/handler/brand"
	"izumi/internal/handler/category"
	"izumi/internal/handler/character"
	"izumi/internal/handler/file"
	"izumi/internal/handler/game"
	"izumi/internal/handler/library"
	"izumi/internal/handler/person"
	"izumi/internal/handler/policy"
	"izumi/internal/handler/scraper"
	"izumi/internal/handler/series"
	"izumi/internal/handler/system_task"
	"izumi/internal/handler/tag"

	"github.com/dokidokikoi/go-common/middleware"
	"github.com/dokidokikoi/go-common/notice"
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
		gG.POST("/ins", middleware.PreHandle(gH.CreateIns))
		gG.GET("", middleware.PreHandle(gH.Get))
		gG.PATCH("", middleware.PreHandle(gH.Update))
		gG.PATCH("/ins", middleware.PreHandle(gH.UpdateIns))
		gG.POST("/download", middleware.PreHandle(gH.DownloadInfo))
		gG.POST("/download/all", middleware.PreHandle(gH.DownloadAllInfo))
		gG.POST("/load/all", middleware.PreHandle(gH.LoadInfo))
		gG.GET("/brief", middleware.PreHandle(gH.GetBrief))
		gG.GET("/panel", middleware.PreHandle(gH.Panel))
	}

	sH := scraper.NewHandler()
	sG := r.Group("/scraper")
	{
		sG.POST("", middleware.PreHandle(sH.Search))
		sG.GET("", middleware.PreHandle(sH.Get))
		sG.POST("/scrape", middleware.PreHandle(sH.Scrape))
		sG.POST("/auto", middleware.PreHandle(sH.AutoScrape))
		sG.POST("/auto-detect", middleware.PreHandle(sH.AutoDetectScrape))
	}

	fH := file.NewHandler()
	fG := r.Group("/file")
	{
		fG.GET("/:name", fH.Get)
		fG.POST("/upload", fH.Upload)
	}

	tH := tag.NewHandler()
	tG := r.Group("/tags")
	{
		tG.GET("", middleware.PreHandle(tH.List))
		tG.POST("", middleware.PreHandle(tH.Create))
		tG.DELETE("", middleware.PreHandle(tH.Del))
		tG.PATCH("", middleware.PreHandle(tH.Update))
		tG.GET("/ehtag", middleware.PreHandle(tH.EtagSearch))
		tG.POST("/ehtag/sync", middleware.PreHandle(tH.EtagSync))
		tG.GET("/decided", middleware.PreHandle(tH.ListDecided))
		tG.PATCH("/decided", middleware.PreHandle(tH.Decide))
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

	brandH := brand.NewHandler()
	brandG := r.Group("/brand")
	{
		brandG.GET("", middleware.PreHandle(brandH.List))
		brandG.POST("", middleware.PreHandle(brandH.Create))
		brandG.DELETE("", middleware.PreHandle(brandH.Del))
		brandG.PATCH("", middleware.PreHandle(brandH.Update))
	}

	characterH := character.NewHandler()
	characterG := r.Group("/character")
	{
		characterG.POST("/search", middleware.PreHandle(characterH.Search))
		characterG.GET("/:id", middleware.PreHandle(characterH.Get))
		characterG.DELETE("", middleware.PreHandle(characterH.Del))
		characterG.PATCH("", middleware.PreHandle(characterH.Update))
		characterG.POST("", middleware.PreHandle(characterH.Create))
	}

	personH := person.NewHandler()
	personG := r.Group("/person")
	{
		personG.GET("/:id", middleware.PreHandle(personH.Get))
		personG.POST("/search", middleware.PreHandle(personH.Search))
		personG.POST("", middleware.PreHandle(personH.Upsert))
		personG.DELETE("", middleware.PreHandle(personH.Del))

	}

	policyH := policy.NewHandler()
	policyG := r.Group("/policy")
	{
		policyG.GET("", middleware.PreHandle(policyH.List))
		policyG.PATCH("", middleware.PreHandle(policyH.Update))
	}

	notifyG := r.Group("/notify")
	{
		notifyG.GET("", notice.ServeWs)
	}

	libraryH := library.NewHandler()
	libraryG := r.Group("/library")
	{
		libraryG.GET("", middleware.PreHandle(libraryH.Ls))
	}

	systemTaskH := system_task.NewHandler()
	systemTaskG := r.Group("/system_task")
	{
		systemTaskG.GET("", middleware.PreHandle(systemTaskH.List))
	}
}
