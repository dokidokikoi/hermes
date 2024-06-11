package game

import (
	"hermes/internal/handler"
	"hermes/model"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

func (h Handler) Create(ctx *gin.Context) {
	requestID := ctx.Query("request_id")
	var input handler.GameVo
	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}
	g := &model.Game{
		Name:      input.Name,
		Cover:     input.Cover,
		Alias:     input.Alias,
		Images:    input.Images,
		Category:  input.Category,
		Series:    input.Series,
		Developer: input.Developer,
		Publisher: input.Publisher,
		Price:     input.Price,
		IssueDate: input.IssueDate,
		Story:     input.Story,
		Platform:  input.Platform,
		Tags:      input.Tags,
		Links:     input.Links,
		OtherInfo: input.OtherInfo,
	}
	var cs []*model.GameCharacter
	for _, c := range input.Characters {
		cs = append(cs, &model.GameCharacter{
			Character: &model.Character{
				ID:      c.ID,
				Name:    c.Name,
				Alias:   c.Alias,
				Gender:  model.GenderMap[c.Gender],
				Summary: c.Summary,
				Images:  c.Images,
				Cover:   c.Cover,
				Tags:    c.Tags,
			},
			Relation: model.CRelationMap[c.Rlation],
		})
	}
	var ss []*model.GameStaff
	for _, s := range input.Staff {
		relations := []model.PersonRelation{}
		for _, r := range s.Relation {
			relations = append(relations, model.PRelationMap[r])
		}
		ss = append(ss, &model.GameStaff{
			Person: &model.Person{
				ID:      s.ID,
				Name:    s.Name,
				Alias:   s.Alias,
				Gender:  model.GenderMap[s.Gender],
				Summary: s.Summary,
				Cover:   s.Cover,
				Images:  s.Images,
				Tags:    s.Tags,
			},
			Relations: relations,
		})
	}
	err := h.srv.Game().CreateL(ctx, g, cs, ss, requestID)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	core.WriteResponse(ctx, nil, nil)
}
