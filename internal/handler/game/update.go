package game

import (
	"context"
	"izumi/db"
	"izumi/internal/handler"
	"izumi/model"

	"github.com/dokidokikoi/go-common/middleware"
	"github.com/dokidokikoi/go-common/tools"
)

func (h Handler) Update(ctx context.Context, input *handler.GameVo, op *middleware.PreHandleOptions) (any, error) {
	g := tools.Ptr(handler.Vo2Game(*input))
	cs := []*model.GameCharacter{}
	ss := []*model.GameStaff{}
	for _, c := range input.Characters {
		cs = append(cs, &model.GameCharacter{
			GameID:      input.ID,
			CharacterID: c.ID,
			Relation:    c.Rlation,
			Character:   tools.Ptr(handler.Vo2Character(c)),
		})
	}
	for _, s := range input.Staff {
		ss = append(ss, &model.GameStaff{
			GameID:    g.ID,
			PersonID:  s.ID,
			Person:    tools.Ptr(handler.Vo2Person(s)),
			Relations: s.Relation,
		})
	}
	if err := h.srv.Game().UpdateL(ctx, g, cs, ss); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h Handler) UpdateIns(ctx context.Context, input *model.GameInstance, op *middleware.PreHandleOptions) (any, error) {
	err := db.GetStore().GameInstance().Update(ctx, input, nil)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
