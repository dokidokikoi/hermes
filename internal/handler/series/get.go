package series

import (
	"context"
	"izumi/db/data"
	"izumi/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h Handler) Get(ctx context.Context, input struct{}) (*model.Series, error) {
	id, err := strconv.ParseUint(ctx.(*gin.Context).Param("id"), 10, 32)
	if err != nil {
		return nil, err
	}
	series, err := data.GetDataFactory().Series().Get(ctx, &model.Series{ID: uint(id)}, nil)
	if err != nil {
		return nil, err
	}

	return series, nil
}
