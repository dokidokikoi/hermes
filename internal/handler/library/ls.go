package library

import (
	"context"
	"hermes/internal/handler"
	"hermes/internal/service"

	"github.com/dokidokikoi/go-common/core"
)

func (h *Handler) Ls(ctx context.Context, req *handler.LsReq) (*core.ListResponseData[service.PathInfo], error) {
	resp := new(core.ListResponseData[service.PathInfo])
	if req.Path == "" {
		req.Path = "/"
	}
	infos, err := h.srv.Library().Ls(ctx, req.Path)
	if err != nil {
		return nil, err
	}
	resp.List = infos
	return resp, nil
}
