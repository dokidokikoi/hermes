package series

import (
	"context"
	"izumi/db/data"
	"izumi/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/middleware"
)

type ListResponse struct {
	List  []*model.Series `json:"list"`
	Total int64           `json:"total"`
}

func (h Handler) List(ctx context.Context, req *struct{}, op *middleware.PreHandleOptions) (*ListResponse, error) {
	list, err := data.GetDataFactory().Series().List(ctx, &model.Series{}, &meta.ListOption{Order: "created_at desc"})
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		List:  list,
		Total: int64(len(list)),
	}, nil
}
