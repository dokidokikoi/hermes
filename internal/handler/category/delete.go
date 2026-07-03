package category

import (
	"context"
	"errors"
	"izumi/db"
	"izumi/model"

	"github.com/dokidokikoi/go-common/middleware"
)

type DelRequest struct {
	IDs []uint `json:"ids"`
}

func (h Handler) Del(ctx context.Context, req *DelRequest, op *middleware.PreHandleOptions) (any, error) {
	categoryIDs := []*model.Category{}
	for _, id := range req.IDs {
		categoryIDs = append(categoryIDs, &model.Category{
			ID: id,
		})
	}

	errs := db.GetStore().Category().DeleteCollection(ctx, categoryIDs, nil)
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return nil, nil
}
