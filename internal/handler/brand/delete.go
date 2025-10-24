package brand

import (
	"context"
	"errors"
	"izumi/db/data"
	"izumi/model"

	"github.com/dokidokikoi/go-common/middleware"
)

type DelRequest struct {
	IDs []uint `json:"ids"`
}

func (h Handler) Del(ctx context.Context, req *DelRequest, op *middleware.PreHandleOptions) (any, error) {
	brandIDs := []*model.Brand{}
	for _, id := range req.IDs {
		brandIDs = append(brandIDs, &model.Brand{
			ID: id,
		})
	}
	errs := data.GetDataFactory().Brand().DeleteCollection(ctx, brandIDs, nil)
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return nil, nil
}
