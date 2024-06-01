package tag

import (
	"hermes/db/data"
	"hermes/model"
	"strconv"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

func (h Handler) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}
	cate, err := data.GetDataFactory().Tag().Get(ctx, &model.Tag{ID: uint(id)}, nil)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	core.WriteResponse(ctx, nil, cate)
}
