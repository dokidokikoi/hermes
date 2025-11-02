package game

import (
	"context"
	"izumi/db/data"
	"izumi/internal/handler"
	"izumi/model"
	"strings"

	"github.com/dokidokikoi/go-common/middleware"
	"github.com/google/uuid"
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
	g := &model.Game{
		UUID:    uuid.NewString(),
		Code:    strings.TrimSpace(input.Game.Code),
		JanCode: strings.TrimSpace(input.Game.JanCode),
		Name:    strings.TrimSpace(input.Game.Name),
		Cover:   input.Game.Cover,
		Alias:   input.Game.Alias,
		Images:  input.Game.Images,
		CategoryID: func() uint {
			if input.Game.Category != nil {
				return input.Game.Category.ID
			}
			return 0
		}(),
		Series: input.Game.Series,
		BrandID: func() uint {
			if input.Game.Brand != nil {
				return input.Game.Brand.ID
			}
			return 0
		}(),
		Price:     input.Game.Price,
		IssueDate: input.Game.IssueDate,
		Story:     strings.TrimSpace(input.Game.Story),
		Tags:      input.Game.Tags,
		Links:     input.Game.Links,
		OtherInfo: input.Game.OtherInfo,
	}
	if g.Brand != nil && g.Brand.Name == "" && g.Brand.ID == 0 {
		g.Brand = nil
	}

	// 角色
	var cs []*model.GameCharacter
	for _, c := range input.Game.Characters {
		if c.Name == "" {
			continue
		}
		cs = append(cs, &model.GameCharacter{
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
			Relation: c.Rlation,
		})
	}

	// 参与人员
	var ss []*model.GameStaff
	for _, s := range input.Game.Staff {
		if s.Name == "" {
			continue
		}
		ss = append(ss, &model.GameStaff{
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
	err := h.srv.Game().SaveFiles(ctx, g, cs, ss)
	if err != nil {
		return 0, err
	}

	// 游戏实体
	var gameIns *model.GameInstance
	if input.GameIns != nil && input.GameIns.Path != "" {
		gameIns = &model.GameInstance{
			Path:     input.GameIns.Path,
			Version:  input.GameIns.Version,
			Language: input.GameIns.Language,
			Size:     input.GameIns.Size,
			Comment:  input.GameIns.Comment,
			Platform: input.GameIns.Platform,
		}
	}

	err = h.srv.Game().CreateL(ctx, g, cs, ss, gameIns)
	if err != nil {
		return 0, err
	}
	return g.ID, nil
}

func (h Handler) CreateIns(ctx context.Context, input *CreateGameInstanceRequst, op *middleware.PreHandleOptions) (uint, error) {
	gameIns := &model.GameInstance{
		GameID:   input.GameID,
		Version:  input.Version,
		Path:     input.Path,
		Size:     input.Size,
		Language: input.Language,
		Comment:  input.Comment,
		Platform: input.Platform,
	}
	if err := data.GetDataFactory().GameInstance().Create(ctx, gameIns, nil); err != nil {
		return 0, err
	}

	return gameIns.ID, nil
}
