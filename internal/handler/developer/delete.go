package developer

import (
	"context"
	"errors"
	"hermes/db/data"
	"hermes/model"
)

type DelRequest struct {
	IDs []uint `json:"ids"`
}

func (h Handler) Del(ctx context.Context, req *DelRequest) (any, error) {
	developerIDs := []*model.Developer{}
	for _, id := range req.IDs {
		developerIDs = append(developerIDs, &model.Developer{
			ID: id,
		})
	}
	errs := data.GetDataFactory().Developer().DeleteCollection(ctx, developerIDs, nil)
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return nil, nil
}
