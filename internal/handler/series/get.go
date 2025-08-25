package series

import (
	"context"
	"hermes/db/data"
	"hermes/model"
	"strconv"

	"github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

func (h Handler) Get(ctx context.Context, id uint64) (*model.Series, *errors.APIError) {
	id, err := strconv.ParseUint(ctx.(*gin.Context).Param("id"), 10, 32)
	if err != nil {
		return nil, errors.Wrap(errors.ApiErrValidation, err)
	}
	series, err := data.GetDataFactory().Series().Get(ctx, &model.Series{ID: uint(id)}, nil)
	if err != nil {
		return nil, errors.Wrap(errors.ApiErrSystemErr, err)
	}

	return series, nil
}
