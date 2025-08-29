package game

import (
	"context"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"

	"github.com/dokidokikoi/go-common/errors"
)

func (h Handler) Update(ctx context.Context, input *handler.GameVo) (any, *errors.APIError) {
	g := &model.Game{
		ID:        input.ID,
		JanCode:   input.JanCode,
		Code:      input.Code,
		Name:      input.Name,
		Cover:     input.Cover,
		Images:    input.Images,
		Alias:     input.Alias,
		Category:  input.Category,
		Developer: input.Developer,
		Price:     input.Price,
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
			Relation:    c.Rlation,
			Character: &model.Character{
				ID:       c.ID,
				Name:     c.Name,
				Alias:    c.Alias,
				Gender:   c.Gender,
				Summary:  c.Summary,
				Images:   c.Images,
				Cover:    c.Cover,
				Tags:     c.Tags,
				PersonID: c.CV.ID,
			},
		})
	}
	for _, s := range input.Staff {
		ss = append(ss, &model.GameStaff{
			GameID:   g.ID,
			PersonID: s.ID,
			Person: &model.Person{
				ID:      s.ID,
				Name:    s.Name,
				Alias:   s.Alias,
				Gender:  s.Gender,
				Summary: s.Summary,
				Cover:   s.Cover,
				Images:  s.Images,
				Tags:    s.Tags,
			},
			Relations: s.Relation,
		})
	}
	if err := h.srv.Game().UpdateL(ctx, g, cs, ss); err != nil {
		return nil, errors.Wrap(errors.ApiErrSystemErr, err)
	}

	return nil, nil
}

func (h Handler) UpdateIns(ctx context.Context, input *model.GameInstance) (any, *errors.APIError) {
	err := data.GetDataFactory().GameInstance().Update(ctx, input, nil)
	if err != nil {
		return nil, errors.Wrap(errors.ApiErrSystemErr, err)
	}
	return nil, nil
}
