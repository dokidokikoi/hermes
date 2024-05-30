package game

import (
	"hermes/internal/handler"
	"hermes/internal/service"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

func (h Handler) Search(ctx *gin.Context) {
	var input handler.GameListReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}

	total, vos, err := h.srv.Game().Search(ctx, input, service.BasicSearchNode...)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}

	core.WriteListResponse(ctx, nil, total, vos)

	// data.GetDataFactory().Game().ListComplex(ctx, &model.Game{})
}
