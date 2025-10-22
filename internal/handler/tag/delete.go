package tag

import (
	"context"
	"errors"
	"izumi/db/data"
	"izumi/model"
)

type DelRequest struct {
	Ids []uint `json:"ids"`
}

func (h Handler) Del(ctx context.Context, req *DelRequest) (any, error) {
	tagIDs := []*model.Tag{}
	for _, id := range req.Ids {
		tagIDs = append(tagIDs, &model.Tag{
			ID: id,
		})
	}
	errs := data.GetDataFactory().Tag().DeleteCollection(ctx, tagIDs, nil)
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return nil, nil
}
