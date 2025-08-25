package tag

import (
	"context"
	"hermes/db/data"
	"hermes/model"

	"github.com/dokidokikoi/go-common/errors"
	meta "github.com/dokidokikoi/go-common/meta/option"
)

type ListResponse struct {
	List  []*model.Tag `json:"list"`
	Total int64        `json:"total"`
}

func (h Handler) List(ctx context.Context, req any) (*ListResponse, *errors.APIError) {
	list, err := data.GetDataFactory().Tag().List(ctx, &model.Tag{}, &meta.ListOption{Order: "created_at desc"})
	if err != nil {
		return nil, errors.Wrap(errors.ApiErrSystemErr, err)
	}

	return &ListResponse{
		List:  list,
		Total: int64(len(list)),
	}, nil
}
