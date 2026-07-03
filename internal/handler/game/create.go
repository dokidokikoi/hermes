package game

import (
	"context"
	"izumi/db"
	"izumi/internal/handler"
	"izumi/model"

	"github.com/dokidokikoi/go-common/middleware"
	"github.com/dokidokikoi/go-common/tools"
)

type CreateGameInstanceRequst struct {
	GameID   uint     `json:"game_id"`
	Version  string   `json:"version"`
	Path     string   `json:"path"`
	Size     int64    `json:"size"`
	Language []string `json:"language"`
	Platform []string `json:"platform"`
	Comment  string   `json:"comment"`
}

type CreateGameRequest struct {
	Game    handler.GameVo      `json:"game"`
	GameIns *model.GameInstance `json:"game_ins"`
}

func (h Handler) Create(ctx context.Context, input *CreateGameRequest, op *middleware.PreHandleOptions) (uint, error) {
	g := tools.GetPtr(handler.Vo2Game(input.Game))

	err := h.srv.Game().SaveFiles(ctx, &input.Game)
	if err != nil {
		return 0, err
	}
	// 角色
	var cs []*model.GameCharacter
	for _, c := range input.Game.Characters {
		if c.Name == "" {
			continue
		}
		cs = append(cs, &model.GameCharacter{
			Character: tools.GetPtr(handler.Vo2Character(c)),
			Relation:  c.Rlation,
		})
	}

	// 参与人员
	var ss []*model.GameStaff
	for _, s := range input.Game.Staff {
		if s.Name == "" {
			continue
		}
		ss = append(ss, &model.GameStaff{
			Person:    tools.GetPtr(handler.Vo2Person(s)),
			Relations: s.Relation,
		})
	}

	// 游戏实体
	var gameIns *model.GameInstance
	if input.GameIns != nil && input.GameIns.Path != "" {
		gameIns = input.GameIns
	}

	err = h.srv.Game().CreateL(ctx, g, cs, ss, gameIns)
	if err != nil {
		return 0, err
	}
	return g.ID, nil
}

func (h Handler) CreateIns(ctx context.Context, input *CreateGameInstanceRequst, op *middleware.PreHandleOptions) (uint, error) {
	gameIns := &model.GameInstance{
		GameID:  input.GameID,
		Version: input.Version,
		Path:    input.Path,
		Size:    input.Size,
		Comment: input.Comment,
	}
	if err := db.GetStore().GameInstance().Create(ctx, gameIns, nil); err != nil {
		return 0, err
	}

	return gameIns.ID, nil
}
