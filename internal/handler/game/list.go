package game

import (
	"hermes/internal/handler"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

func (h Handler) List(ctx *gin.Context) {
	var input handler.GameListReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}

	// data.GetDataFactory().Game().ListComplex(ctx, &model.Game{})
}
