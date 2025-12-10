package system_task

import (
	"context"
	"izumi/db/data"
	"izumi/model"

	"github.com/dokidokikoi/go-common/middleware"
	"github.com/dokidokikoi/go-common/query"
)

type ListReq struct {
	query.PageQuery
	Type  model.SystemTaskType  `form:"type"`
	State model.SystemTaskState `form:"state"`
}

type ListResponse struct {
	List []*model.SystemTask `json:"list"`
}

func (h *Handler) List(ctx context.Context, input *ListReq, op *middleware.PreHandleOptions) (*ListResponse, error) {
	input.Order = "id desc"
	ts, err := data.GetDataFactory().SystemTask().List(ctx, &model.SystemTask{
		Type:  input.Type,
		State: input.State,
	}, input.GetListOption())
	if err != nil {
		return nil, err
	}
	return &ListResponse{
		List: ts,
	}, nil
}
