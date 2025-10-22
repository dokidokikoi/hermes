package tag

import (
	"context"
	"izumi/db/data"
	"izumi/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

type ListResponse struct {
	List  []*model.Tag `json:"list"`
	Total int64        `json:"total"`
}

func (h Handler) List(ctx context.Context, req *struct{}) (*ListResponse, error) {
	list, err := data.GetDataFactory().Tag().List(ctx, &model.Tag{}, &meta.ListOption{Order: "created_at desc"})
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		List:  list,
		Total: int64(len(list)),
	}, nil
}
