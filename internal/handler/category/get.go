package category

import (
	"context"
	"izumi/db/data"
	"izumi/model"
	"strconv"

	"github.com/dokidokikoi/go-common/middleware"
	"github.com/gin-gonic/gin"
)

func (h Handler) Get(ctx context.Context, input struct{}, op *middleware.PreHandleOptions) (*model.Category, error) {
	id, err := strconv.ParseUint(ctx.(*gin.Context).Param("id"), 10, 32)
	if err != nil {
		return nil, err
	}
	cate, err := data.GetDataFactory().Category().Get(ctx, &model.Category{ID: uint(id)}, nil)
	if err != nil {
		return nil, err
	}
	return cate, nil
}
