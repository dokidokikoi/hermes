package game

import (
	"context"
	"izumi/db/data"
	"izumi/internal/handler"
	"izumi/model"

	"github.com/dokidokikoi/go-common/middleware"
)

func (h Handler) Update(ctx context.Context, input *handler.GameVo, op *middleware.PreHandleOptions) (any, error) {
	g := &model.Game{
		ID:        input.ID,
		JanCode:   input.JanCode,
		Code:      input.Code,
		Name:      input.Name,
		Cover:     input.Cover,
		Images:    input.Images,
		Alias:     input.Alias,
		Category:  input.Category,
		Brand:     input.Brand,
		Price:     input.Price,
		Story:     input.Story,
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
				Weight:   c.Weight,
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
		return nil, err
	}

	return nil, nil
}

func (h Handler) UpdateIns(ctx context.Context, input *model.GameInstance) (any, error) {
	err := data.GetDataFactory().GameInstance().Update(ctx, input, nil)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
