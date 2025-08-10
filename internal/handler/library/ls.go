package library

import (
	"context"
	"hermes/internal/handler"
	"hermes/internal/service"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
)

func (h *Handler) Ls(ctx context.Context, req handler.LsReq) (*core.ListResponseData[service.PathInfo], *errors.APIError) {
	resp := new(core.ListResponseData[service.PathInfo])
	if req.Path == "" {
		req.Path = "/"
	}
	infos, err := h.srv.Library().Ls(ctx, req.Path)
	if err != nil {
		return nil, errors.Wrap(errors.ApiErrSystemErr, err)
	}
	resp.List = infos
	return resp, nil
}
