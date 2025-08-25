package person

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
	PersonIDs := []*model.Person{}
	for _, id := range req.IDs {
		PersonIDs = append(PersonIDs, &model.Person{
			ID: id,
		})
	}

	errs := data.GetDataFactory().Person().DeleteCollection(ctx, PersonIDs, nil)
	if len(errs) > 0 {
		return nil, comm_errs.Wrap(comm_errs.ApiErrSystemErr, errors.Join(errs...))
	}

	return nil, nil
}
