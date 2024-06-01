package game

import (
	"hermes/internal/handler"
	"hermes/model"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

func (h Handler) Update(ctx *gin.Context) {
	var input handler.GameVo
	if err := ctx.ShouldBindJSON(&input); err != nil || input.ID == 0 {
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}

	g := &model.Game{
		ID:        input.ID,
		Name:      input.Name,
		Cover:     input.Cover,
		Images:    input.Images,
		Alias:     input.Alias,
		Category:  input.Category,
		Publisher: input.Publisher,
		Developer: input.Developer,
		Price:     input.Price,
		Platform:  input.Platform,
		Story:     input.Platform,
		Series:    input.Series,
		IssueDate: input.IssueDate,
		Links:     input.Links,
		OtherInfo: input.OtherInfo,
		Tags:      input.Tags,
	}
	cs := []*model.GameCharacter{}
	ss := []*model.GameStaff{}
	for _, c := range input.Characters {
		cs = append(cs, &model.GameCharacter{
			GameID:      input.ID,
			CharacterID: c.ID,
			Relation:    model.CRelationMap[c.Rlation],
		})
	}
	for _, s := range input.Staff {
		for _, sr := range s.Relation {
			ss = append(ss, &model.GameStaff{
				GameID:   input.ID,
				PersonID: s.ID,
				Relation: model.PRelationMap[sr],
			})
		}
	}
	if err := h.srv.Game().UpdateL(ctx, g, cs, ss); err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}

	core.WriteResponse(ctx, nil, nil)
}
