package game

import (
	"hermes/model"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

func (h Handler) Create(ctx *gin.Context) {
	var input model.Game
	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}
	input.ID = 0
	err := h.srv.Game().CreateL(ctx, &input)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	core.WriteResponse(ctx, nil, nil)
}
