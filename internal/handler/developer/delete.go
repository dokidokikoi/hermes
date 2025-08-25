package developer

import (
	"context"
	"errors"
	"hermes/db/data"
	"hermes/model"

	comm_errs "github.com/dokidokikoi/go-common/errors"
)

type DelRequest struct {
	IDs []uint `json:"ids"`
}

func (h Handler) Del(ctx context.Context, req *DelRequest) (any, *comm_errs.APIError) {
	developerIDs := []*model.Developer{}
	for _, id := range req.IDs {
		developerIDs = append(developerIDs, &model.Developer{
			ID: id,
		})
	}
	errs := data.GetDataFactory().Developer().DeleteCollection(ctx, developerIDs, nil)
	if len(errs) > 0 {
		return nil, comm_errs.Wrap(comm_errs.ApiErrSystemErr, errors.Join(errs...))
	}

	return nil, nil
}
