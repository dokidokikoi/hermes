package character

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
	CharacterIDs := []*model.Character{}
	for _, id := range req.IDs {
		CharacterIDs = append(CharacterIDs, &model.Character{
			ID: id,
		})
	}

	errs := data.GetDataFactory().Character().DeleteCollection(ctx, CharacterIDs, nil)
	if len(errs) > 0 {
		return nil, comm_errs.Wrap(comm_errs.ApiErrDatabaseOp, errors.Join(errs...))
	}
	return nil, nil
}
