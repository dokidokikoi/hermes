package character

import (
	"context"
	"izumi/db"
	"izumi/internal/handler"
	"izumi/model"
	"strconv"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/middleware"
	"github.com/gin-gonic/gin"
)

func (h Handler) Get(ctx context.Context, req *struct{}, op *middleware.PreHandleOptions) (*handler.CharacterVo, error) {
	id, err := strconv.ParseUint(ctx.(*gin.Context).Param("id"), 10, 32)
	if err != nil {
		return nil, err
	}
	c, err := db.GetStore().Character().Get(ctx, &model.Character{ID: uint(id)}, nil)
	if err != nil {
		return nil, err
	}

	gs, err := db.GetStore().Game().List(ctx, &model.Game{}, &meta.ListOption{
		GetOption: meta.GetOption{
			Join: []*meta.Join{
				{
					Method:         meta.INNER_JOIN,
					Table:          model.Game{}.TableName(),
					JoinTable:      model.GameCharacter{}.TableName(),
					TableField:     "id",
					JoinTableField: "game_id",
					JoinTableCondition: []meta.Condition{
						{
							Field:    "character_id",
							Operator: meta.EQUAL,
							Value:    c.ID,
						},
					},
				},
			},
			Select: []string{"ID", "Name"},
		},
	})
	if err != nil {
		return nil, err
	}

	cgvos := make([]handler.CharacterGameVo, len(gs))
	for i, g := range gs {
		cgvos[i] = handler.CharacterGameVo{
			ID:   g.ID,
			Name: g.Name,
		}
	}

	vo := handler.CharacterVo{
		ID:      c.ID,
		Name:    c.Name,
		Alias:   c.Alias,
		Gender:  c.Gender,
		Summary: c.Summary,
		Cover:   c.Cover,
		Images:  c.Images,
		CV: handler.StaffVo{
			ID:        c.CV.ID,
			Name:      c.CV.Name,
			Alias:     c.CV.Alias,
			Cover:     c.CV.Cover,
			Images:    c.CV.Images,
			Summary:   c.CV.Summary,
			Gender:    c.CV.Gender,
			CreatedAt: c.CV.CreatedAt,
		},
		CreatedAt: c.CreatedAt,
		Games:     cgvos,
	}

	return &vo, nil
}
