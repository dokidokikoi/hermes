package game

import (
	"context"
	"hermes/db/data"
	"hermes/model"

	"github.com/dokidokikoi/go-common/errors"
	meta "github.com/dokidokikoi/go-common/meta/option"
)

type GetInsRequest struct {
	ID uint `form:"id" query:"id" json:"id"`
}

func (h Handler) GetIns(ctx context.Context, req GetInsRequest) ([]*model.GameInstance, *errors.APIError) {
	gIns, err := data.GetDataFactory().GameInstance().List(ctx, &model.GameInstance{GameID: req.ID}, nil)
	if err != nil {
		return nil, errors.ApiErrSystemErr
	}

	return gIns, nil
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
