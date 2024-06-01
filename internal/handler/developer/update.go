package developer

import (
	"errors"
	"hermes/db/data"
	"hermes/model"

	"github.com/dokidokikoi/go-common/core"
	comm_errors "github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

func (h Handler) Update(ctx *gin.Context) {
	var input model.Developer
	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.WriteResponse(ctx, comm_errors.ApiErrValidation, nil)
		return
	}

	if err := data.GetDataFactory().Developer().Update(ctx, &input, nil); err != nil {
		if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
			core.WriteResponse(ctx, comm_errors.ApiErrSystemErr, nil)
			return
		}
	}

	core.WriteResponse(ctx, nil, nil)
}
