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
	vos := make([]handler.CharacterVo, len(list))
	for i := range list {
		vos[i] = handler.CharacterVo{
			ID:      list[i].ID,
			Name:    list[i].Name,
			Alias:   list[i].Alias,
			Gender:  list[i].Gender.String(),
			Summary: list[i].Summary,
			Cover:   list[i].Cover,
			Images:  list[i].Images,
			CV: handler.StaffVo{
				ID:        list[i].CV.ID,
				Name:      list[i].CV.Name,
				Alias:     list[i].CV.Alias,
				Cover:     list[i].CV.Cover,
				Images:    list[i].CV.Images,
				Tags:      list[i].CV.Tags,
				Summary:   list[i].CV.Summary,
				Gender:    list[i].CV.Gender.String(),
				CreatedAt: list[i].CV.CreatedAt,
			},
			Tags:      list[i].Tags,
			CreatedAt: list[i].CreatedAt,
		}
	}
	core.WriteListResponse(ctx, nil, total, vos)
}
