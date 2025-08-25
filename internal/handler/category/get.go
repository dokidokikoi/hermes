package category

import (
	"context"
	"hermes/db/data"
	"hermes/model"
	"strconv"

	"github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

func (h Handler) Get(ctx context.Context, id uint64) (*model.Category, *errors.APIError) {
	id, err := strconv.ParseUint(ctx.(*gin.Context).Param("id"), 10, 32)
	if err != nil {
		return nil, errors.Wrap(errors.ApiErrValidation, err)
	}
	cate, err := data.GetDataFactory().Category().Get(ctx, &model.Category{ID: uint(id)}, nil)
	if err != nil {
		return nil, errors.Wrap(errors.ApiErrSystemErr, err)
	}
	return cate, nil
}
