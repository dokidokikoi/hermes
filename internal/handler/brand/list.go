package brand

import (
	"context"
	"izumi/db"
	"izumi/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/middleware"
)

type ListResponse struct {
	List  []*model.Brand `json:"list"`
	Total int64          `json:"total"`
}

func (h Handler) List(ctx context.Context, req *struct{}, op *middleware.PreHandleOptions) (any, error) {
	list, err := db.GetStore().Brand().List(ctx, &model.Brand{}, &meta.ListOption{Order: "created_at desc"})
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		List:  list,
		Total: int64(len(list)),
	}, nil
}
