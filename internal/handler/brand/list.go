package brand

import (
	"context"
	"izumi/db/data"
	"izumi/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

type ListResponse struct {
	List  []*model.Brand `json:"list"`
	Total int64          `json:"total"`
}

func (h Handler) List(ctx context.Context, req *struct{}) (any, error) {
	list, err := data.GetDataFactory().Brand().List(ctx, &model.Brand{}, &meta.ListOption{Order: "created_at desc"})
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		List:  list,
		Total: int64(len(list)),
	}, nil
}
