package library

import (
	"context"
	"izumi/internal/handler"
	"izumi/internal/service"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/middleware"
)

func (h *Handler) Ls(ctx context.Context, req *handler.LsReq, op *middleware.PreHandleOptions) (*core.ListResponseData[service.PathInfo], error) {
	resp := new(core.ListResponseData[service.PathInfo])
	if req.Path == "" {
		req.Path = "/"
	}
	infos, err := h.srv.Library().Ls(ctx, req.Path, req.OnlyNoScrap)
	if err != nil {
		return nil, err
	}
	resp.List = infos
	return resp, nil
}
