package developer

import (
	"context"
	"hermes/db/data"
	"hermes/model"

	"github.com/dokidokikoi/go-common/errors"
	meta "github.com/dokidokikoi/go-common/meta/option"
)

type ListResponse struct {
	List  []*model.Developer `json:"list"`
	Total int64              `json:"total"`
}

func (h Handler) List(ctx context.Context, req *struct{}) (any, *errors.APIError) {
	list, err := data.GetDataFactory().Developer().List(ctx, &model.Developer{}, &meta.ListOption{Order: "created_at desc"})
	if err != nil {
		return nil, errors.ApiErrSystemErr
	}

	return &ListResponse{
		List:  list,
		Total: int64(len(list)),
	}, nil
}
