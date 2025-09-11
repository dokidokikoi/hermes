package developer

import (
	"context"
	"hermes/db/data"
	"hermes/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h Handler) Get(ctx context.Context, input struct{}) (*model.Developer, error) {
	id, err := strconv.ParseUint(ctx.(*gin.Context).Param("id"), 10, 32)
	if err != nil {
		return nil, err
	}
	developer, err := data.GetDataFactory().Developer().Get(ctx, &model.Developer{ID: uint(id)}, nil)
	if err != nil {
		return nil, err
	}

	return developer, nil
}
