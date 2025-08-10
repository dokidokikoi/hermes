package game

import (
	"context"
	"hermes/internal/handler"
	"hermes/internal/service"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/query"
)

func (h Handler) Search(ctx context.Context, req handler.GameListReq) (*core.ListResponseData[handler.GameVo], *errors.APIError) {
	var q query.PageQuery
	q.Page = req.Page
	q.PageSize = req.PageSize

	total, vos, err := h.srv.Game().Search(ctx, req, q.GetListOption(), service.GameBasicSearchNode...)
	if err != nil {
		return nil, errors.Wrap(errors.ApiErrSystemErr, err)
	}
	return &core.ListResponseData[handler.GameVo]{
		List:  vos,
		Total: total,
	}, nil
}
