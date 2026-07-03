package person

import (
	"context"
	"izumi/db"
	"izumi/internal/handler"
	"izumi/model"
	"strings"

	"github.com/dokidokikoi/go-common/middleware"
)

func (h Handler) Upsert(ctx context.Context, input *handler.StaffVo, op *middleware.PreHandleOptions) (uint, error) {
	p := &model.Person{
		ID:      input.ID,
		Name:    strings.TrimSpace(input.Name),
		Alias:   input.Alias,
		Cover:   input.Cover,
		Images:  input.Images,
		Summary: input.Summary,
		Gender:  input.Gender,
	}

	var err error
	if input.ID == 0 {
		err = db.GetStore().Person().Create(ctx, p, nil)
	} else {
		err = db.GetStore().Person().Update(ctx, p, nil)
	}
	if err != nil {
		return 0, err
	}

	return p.ID, nil
}
