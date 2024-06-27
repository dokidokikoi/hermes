package person

import (
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h Handler) Upsert(ctx *gin.Context) {
	var input handler.StaffVo
	if err := ctx.ShouldBindJSON(&input); err != nil {
		zaplog.L().Error(errors.ApiErrValidation.Message, zap.Error(err))
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}

	p := &model.Person{
		ID:      input.ID,
		Name:    input.Name,
		Alias:   input.Alias,
		Cover:   input.Cover,
		Images:  input.Images,
		Tags:    input.Tags,
		Summary: input.Summary,
		Gender:  model.GenderMap[input.Gender],
	}

	var err error
	if input.ID == 0 {
		err = data.GetDataFactory().Person().Create(ctx, p, nil)
	} else {
		err = data.GetDataFactory().Person().Update(ctx, p, nil)
	}
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	core.WriteResponse(ctx, nil, input.ID)
}
