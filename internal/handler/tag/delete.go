package tag

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
	tagIDs := []*model.Tag{}
	for _, id := range req.Ids {
		tagIDs = append(tagIDs, &model.Tag{
			ID: id,
		})
	}
	errs := data.GetDataFactory().Tag().DeleteCollection(ctx, tagIDs, nil)
	if len(errs) > 0 {
		return nil, comm_errs.Wrap(comm_errs.ApiErrSystemErr, errors.Join(errs...))
	}

	return nil, nil
}
