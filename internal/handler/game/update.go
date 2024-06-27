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
		JanCode:   input.JanCode,
		Code:      input.Code,
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
			Character: &model.Character{
				ID:       c.ID,
				Name:     c.Name,
				Alias:    c.Alias,
				Gender:   model.GenderMap[c.Gender],
				Summary:  c.Summary,
				Images:   c.Images,
				Cover:    c.Cover,
				Tags:     c.Tags,
				PersonID: c.CV.ID,
			},
		})
	}
	for _, s := range input.Staff {
		relations := []model.PersonRelation{}
		for _, r := range s.Relation {
			relations = append(relations, model.PRelationMap[r])
		}
		ss = append(ss, &model.GameStaff{
			GameID:   g.ID,
			PersonID: s.ID,
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
	if err := h.srv.Game().UpdateL(ctx, g, cs, ss); err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}

	core.WriteResponse(ctx, nil, nil)
}
