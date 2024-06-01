package character

import (
	"hermes/internal/handler"
	"hermes/internal/service"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/query"
	"github.com/gin-gonic/gin"
)

func (h Handler) Search(ctx *gin.Context) {
	var input handler.CharacterListReq
	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}
	var q query.PageQuery
	ctx.ShouldBindQuery(&q)

	total, list, err := h.srv.Character().Search(ctx, input, q.GetListOption(), service.CharacterBasicSearchNode...)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	core.WriteListResponse(ctx, nil, total, list)
}
