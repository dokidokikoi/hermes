package character

import (
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"
	"strconv"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/gin-gonic/gin"
)

func (h Handler) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}
	c, err := data.GetDataFactory().Character().Get(ctx, &model.Character{ID: uint(id)}, nil)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}

	gs, err := data.GetDataFactory().Game().List(ctx, &model.Game{}, &meta.ListOption{
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
		Gender:  c.Gender.String(),
		Summary: c.Summary,
		Cover:   c.Cover,
		Images:  c.Images,
		CV: handler.StaffVo{
			ID:        c.CV.ID,
			Name:      c.CV.Name,
			Alias:     c.CV.Alias,
			Cover:     c.CV.Cover,
			Images:    c.CV.Images,
			Tags:      c.CV.Tags,
			Summary:   c.CV.Summary,
			Gender:    c.CV.Gender.String(),
			CreatedAt: c.CV.CreatedAt,
		},
		Tags:      c.Tags,
		CreatedAt: c.CreatedAt,
		Games:     cgvos,
	}
	core.WriteResponse(ctx, nil, vo)
}
