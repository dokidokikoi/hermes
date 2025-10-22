package series

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
	seriesIDs := []*model.Series{}
	for _, id := range req.Ids {
		seriesIDs = append(seriesIDs, &model.Series{
			ID: id,
		})
	}
	errs := data.GetDataFactory().Series().DeleteCollection(ctx, seriesIDs, nil)
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return nil, nil
}
