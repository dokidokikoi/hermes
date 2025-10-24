package person

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
	PersonIDs := []*model.Person{}
	for _, id := range req.IDs {
		PersonIDs = append(PersonIDs, &model.Person{
			ID: id,
		})
	}

	errs := data.GetDataFactory().Person().DeleteCollection(ctx, PersonIDs, nil)
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return nil, nil
}
