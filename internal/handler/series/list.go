package series

import (
	"context"
	"hermes/db/data"
	"hermes/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

type ListResponse struct {
	List  []*model.Series `json:"list"`
	Total int64           `json:"total"`
}

func (h Handler) List(ctx context.Context, req *struct{}) (*ListResponse, error) {
	list, err := data.GetDataFactory().Series().List(ctx, &model.Series{}, &meta.ListOption{Order: "created_at desc"})
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		List:  list,
		Total: int64(len(list)),
	}, nil
}
