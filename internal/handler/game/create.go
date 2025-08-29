package game

import (
	"context"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"

	"github.com/dokidokikoi/go-common/errors"
)

type CreateGameInstanceRequst struct {
	GameID   uint     `json:"game_id"`
	Version  string   `json:"version"`
	Path     string   `json:"path"`
	Size     int64    `json:"size"`
	Language []string `json:"language"`
	Comment  string   `json:"comment"`
}

type CreateGameRequest struct {
	Game    handler.GameVo      `json:"game"`
	GameIns *model.GameInstance `json:"game_ins"`
}

func (h Handler) Create(ctx context.Context, input *CreateGameRequest) (uint, *errors.APIError) {
	g := &model.Game{
		Code:        input.Game.Code,
		JanCode:     input.Game.JanCode,
		Name:        input.Game.Name,
		Cover:       input.Game.Cover,
		Alias:       input.Game.Alias,
		Images:      input.Game.Images,
		CategoryID:  input.Game.Category.ID,
		Series:      input.Game.Series,
		DeveloperID: input.Game.Developer.ID,
		Price:       input.Game.Price,
		IssueDate:   input.Game.IssueDate,
		Story:       input.Game.Story,
		Tags:        input.Game.Tags,
		Links:       input.Game.Links,
		OtherInfo:   input.Game.OtherInfo,
	}
	if g.Developer != nil && g.Developer.Name == "" && g.Developer.ID == 0 {
		g.Developer = nil
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
		return 0, errors.Wrap(errors.ApiErrSystemErr, err)
	}

	// 游戏实体
	var gameIns *model.GameInstance
	if input.GameIns != nil {
		gameIns = &model.GameInstance{
			Path:     input.GameIns.Path,
			Version:  input.GameIns.Version,
			Language: input.GameIns.Language,
			Size:     input.GameIns.Size,
			Comment:  input.GameIns.Comment,
		}
	}

	err = h.srv.Game().CreateL(ctx, g, cs, ss, gameIns)
	if err != nil {
		return 0, errors.Wrap(errors.ApiErrSystemErr, err)
	}
	return g.ID, nil
}

func (h Handler) CreateIns(ctx context.Context, input *CreateGameInstanceRequst) (uint, *errors.APIError) {
	gameIns := &model.GameInstance{
		GameID:   input.GameID,
		Version:  input.Version,
		Path:     input.Path,
		Size:     input.Size,
		Language: input.Language,
		Comment:  input.Comment,
	}
	if err := data.GetDataFactory().GameInstance().Create(ctx, gameIns, nil); err != nil {
		return 0, errors.Wrap(errors.ApiErrSystemErr, err)
	}

	return gameIns.ID, nil
}
