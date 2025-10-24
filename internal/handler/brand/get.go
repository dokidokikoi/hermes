package brand

import (
	"context"
	"izumi/db/data"
	"izumi/model"
	"strconv"

	"github.com/dokidokikoi/go-common/middleware"
	"github.com/gin-gonic/gin"
)

func (h Handler) Get(ctx context.Context, input struct{}, op *middleware.PreHandleOptions) (*model.Brand, error) {
	id, err := strconv.ParseUint(ctx.(*gin.Context).Param("id"), 10, 32)
	if err != nil {
		return nil, err
	}
	brand, err := data.GetDataFactory().Brand().Get(ctx, &model.Brand{ID: uint(id)}, nil)
	if err != nil {
		return nil, err
	}

	return brand, nil
}
