package internal

import (
	"hermes/internal/handler/game"

	"github.com/gin-gonic/gin"
)

func Install(r gin.IRouter) {
	gH := game.NewHandler()
	gG := r.Group("/game")
	{
		gG.PUT("", gH.Create)
	}
}
