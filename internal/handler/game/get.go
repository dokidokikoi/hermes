package game

import (
	"context"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"

	"github.com/dokidokikoi/go-common/errors"
	meta "github.com/dokidokikoi/go-common/meta/option"
)

type GetInsRequest struct {
	ID      uint   `form:"id" query:"id" json:"id"`
	Version string `form:"version" query:"version" json:"version"`
}

func (h Handler) GetIns(ctx context.Context, req GetInsRequest) (*handler.GameInsVo, *errors.APIError) {
	gvo, err := h.srv.Game().GetVOByID(ctx, req.ID)
	if err != nil {
		return nil, errors.ApiErrSystemErr
	}
	gInsVo := &handler.GameInsVo{
		ID:         gvo.ID,
		Name:       gvo.Name,
		Alias:      gvo.Alias,
		Cover:      gvo.Cover,
		Images:     gvo.Images,
		Category:   gvo.Category,
		Series:     gvo.Series,
		Developer:  gvo.Developer,
		Publisher:  gvo.Publisher,
		Price:      gvo.Price,
		IssueDate:  gvo.IssueDate,
		Story:      gvo.Story,
		Platform:   gvo.Platform,
		Tags:       gvo.Tags,
		Characters: gvo.Characters,
		Staff:      gvo.Staff,
		Links:      gvo.Links,
		OtherInfo:  gvo.OtherInfo,
	}

	ver := req.Version
	if ver != "" {
		gIns, err := data.GetDataFactory().GameInstance().Get(ctx, &model.GameInstance{GameID: req.ID, Version: ver}, nil)
		if err != nil {
			return nil, errors.ApiErrSystemErr
		}
		gInsVo.Version = gIns.Version
		gInsVo.Path = gIns.Path
		gInsVo.Language = gIns.Language
		gInsVo.Size = gIns.Size
		gInsVo.CreatedAt = gIns.CreatedAt
	}

	return gInsVo, nil
}

type GetVerRequest struct {
	ID uint `form:"id" query:"id" json:"id"`
}

func (h Handler) GetVer(ctx context.Context, req GetVerRequest) ([]string, *errors.APIError) {
	// veriosn
	gins, err := data.GetDataFactory().GameInstance().List(ctx, &model.GameInstance{GameID: req.ID}, &meta.ListOption{GetOption: meta.GetOption{Select: []string{"ID", "GameID", "Version"}}})
	if err != nil {
		return nil, errors.ApiErrSystemErr
	}
	var versions []string
	for _, v := range gins {
		versions = append(versions, v.Version)
	}

	return versions, nil
}
