package character

import (
	"hermes/db/data"
	"hermes/model"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

func (h Handler) Create(ctx *gin.Context) {
	var input model.Character
	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}

	if err := data.GetDataFactory().Character().Create(ctx, &input, nil); err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}

	core.WriteResponse(ctx, nil, nil)
}
