package game

import (
	"context"
	"izumi/internal/handler"
	"izumi/internal/service"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/middleware"
	"github.com/dokidokikoi/go-common/query"
)

func (h Handler) Search(ctx context.Context, req *handler.GameListReq, op *middleware.PreHandleOptions) (*core.ListResponseData[handler.GameVo], error) {
	var q query.PageQuery
	q.Page = req.Page
	q.PageSize = req.PageSize
	q.Order = req.OrderBy

	total, vos, err := h.srv.Game().Search(ctx, *req, q.GetListOption(), service.GameBasicSearchNode...)
	if err != nil {
		return nil, err
	}
	return &core.ListResponseData[handler.GameVo]{
		List:  vos,
		Total: total,
	}, nil
}
