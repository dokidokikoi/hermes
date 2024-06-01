package game

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

	gvo, err := h.srv.Game().GetVOByID(ctx, uint(id))
	if err != nil {
		core.WriteResponse(ctx, nil, errors.ErrSystemErr)
		return
	}

	ver := ctx.Query("version")
	if ver == "" {
		// veriosn
		gins, err := data.GetDataFactory().GameInstance().List(ctx, &model.GameInstance{GameID: uint(id)}, &meta.ListOption{GetOption: meta.GetOption{Select: []string{"ID", "GameID", "Version"}}})
		if err != nil {
			core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
			return
		}
		var version []string
		for _, v := range gins {
			version = append(version, v.Version)
		}
		gvo.Versions = version

		core.WriteResponse(ctx, nil, gvo)
		return
	}

	gIns, err := data.GetDataFactory().GameInstance().Get(ctx, &model.GameInstance{GameID: uint(id), Version: ver}, nil)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	gInsVo := handler.GameInsVo{
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

		Version:   gIns.Version,
		Path:      gIns.Path,
		Language:  gIns.Language,
		Size:      gIns.Size,
		CreatedAt: gIns.CreatedAt,
	}
	core.WriteResponse(ctx, nil, gInsVo)
}
