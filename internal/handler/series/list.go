package series

import (
	"hermes/db/data"
	"hermes/model"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

func (h Handler) List(ctx *gin.Context) {
	list, err := data.GetDataFactory().Series().List(ctx, &model.Series{}, nil)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	core.WriteListResponse(ctx, nil, int64(len(list)), list)
}
