package category

import (
	"context"
	"hermes/db/data"
	"hermes/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h Handler) Get(ctx context.Context, input struct{}) (*model.Category, error) {
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
