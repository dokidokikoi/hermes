package system_task

import (
	"context"
	"izumi/db/data"
	"izumi/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/middleware"
)

type ListReq struct {
	Type model.SystemTaskType `form:"type"`
}

type ListResponse struct {
	List []*model.SystemTask `json:"list"`
}

func (h *Handler) List(ctx context.Context, input *ListReq, op *middleware.PreHandleOptions) (*ListResponse, error) {
	ts, err := data.GetDataFactory().SystemTask().List(ctx, &model.SystemTask{
		Type: input.Type,
	}, &meta.ListOption{Order: "id desc", Page: 1, PageSize: 5})
	if err != nil {
		return nil, err
	}
	return &ListResponse{
		List: ts,
	}, nil
}
