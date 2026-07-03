package tag

import (
	"context"
	"izumi/db"
	"izumi/model"
	"sort"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/middleware"
)

type ListResponse struct {
	List  []*model.Tag `json:"list"`
	Total int64        `json:"total"`
}

func (h Handler) List(ctx context.Context, req *struct{}, op *middleware.PreHandleOptions) (*ListResponse, error) {
	list, err := db.GetStore().Tag().List(ctx, &model.Tag{}, &meta.ListOption{Order: "created_at desc"})
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool {
		if list[j].Lang == "zh" {
			return false
		} else if list[j].Lang == "ja" && list[i].Lang != "zh" {
			return false
		} else {
			return true
		}
	})

	return &ListResponse{
		List:  list,
		Total: int64(len(list)),
	}, nil
}
