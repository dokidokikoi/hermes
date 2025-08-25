package category

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
	categoryIDs := []*model.Category{}
	for _, id := range req.IDs {
		categoryIDs = append(categoryIDs, &model.Category{
			ID: id,
		})
	}

	errs := data.GetDataFactory().Category().DeleteCollection(ctx, categoryIDs, nil)
	if len(errs) > 0 {
		return nil, comm_errs.Wrap(comm_errs.ApiErrSystemErr, errors.Join(errs...))
	}
	return nil, nil
}
