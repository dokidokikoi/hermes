package tag

import (
	"context"
	"izumi/db"
	"izumi/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/middleware"
)

type ListRequest struct {
	Ns       string `form:"ns"`
	Key      string `form:"key"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type ListResponse struct {
	List  []*model.Tag `json:"list"`
	Total int64        `json:"total"`
}

func (h Handler) List(ctx context.Context, req *ListRequest, op *middleware.PreHandleOptions) (*ListResponse, error) {
	list, err := db.GetStore().Tag().List(
		ctx,
		&model.Tag{NS: req.Ns, Key: req.Key},
		meta.NewListOption(
			[]string{},
			meta.WithOrderBy("created_at desc"),
			meta.WithPage(req.Page),
			meta.WithPageSize(req.PageSize),
		))
	if err != nil {
		return nil, err
	}

	total, err := db.GetStore().Tag().Count(
		ctx,
		&model.Tag{NS: req.Ns, Key: req.Key},
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		List:  list,
		Total: total,
	}, nil
}
