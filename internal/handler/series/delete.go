package series

import (
	"context"
	"errors"
	"hermes/db/data"
	"hermes/model"

	comm_errs "github.com/dokidokikoi/go-common/errors"
)

type DelRequest struct {
	Ids []uint `json:"ids"`
}

func (h Handler) Del(ctx context.Context, req *DelRequest) (any, *comm_errs.APIError) {
	seriesIDs := []*model.Series{}
	for _, id := range req.Ids {
		seriesIDs = append(seriesIDs, &model.Series{
			ID: id,
		})
	}
	errs := data.GetDataFactory().Series().DeleteCollection(ctx, seriesIDs, nil)
	if len(errs) > 0 {
		return nil, comm_errs.Wrap(comm_errs.ApiErrSystemErr, errors.Join(errs...))
	}

	return nil, nil
}
