package character

import (
	"errors"
	"hermes/db/data"
	"hermes/model"

	"github.com/dokidokikoi/go-common/core"
	comm_errors "github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

func (h Handler) Update(ctx *gin.Context) {
	var input model.Character
	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.WriteResponse(ctx, comm_errors.ApiErrValidation, nil)
		return
	}

	tx := data.GetDataFactory().Transaction().Begin()
	err := tx.CharacterTag().Delete(ctx, &model.CharacterTag{CharacterID: input.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		core.WriteResponse(ctx, comm_errors.ApiErrSystemErr, nil)
		return
	}

	if err := tx.Character().Update(ctx, &input, nil); err != nil {
		if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
			tx.Transaction().Rollback()
			core.WriteResponse(ctx, comm_errors.ApiErrSystemErr, nil)
			return
		}
	}

	core.WriteResponse(ctx, nil, nil)
}
