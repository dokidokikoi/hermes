package person

import (
	"context"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"
	"strings"

	"github.com/google/uuid"
)

func (h Handler) Upsert(ctx context.Context, input *handler.StaffVo) (uint, error) {
	p := &model.Person{
		ID:      input.ID,
		UUID:    uuid.NewString(),
		Name:    strings.TrimSpace(input.Name),
		Alias:   input.Alias,
		Cover:   input.Cover,
		Images:  input.Images,
		Tags:    input.Tags,
		Summary: input.Summary,
		Gender:  input.Gender,
	}

	var err error
	if input.ID == 0 {
		err = data.GetDataFactory().Person().Create(ctx, p, nil)
	} else {
		err = data.GetDataFactory().Person().Update(ctx, p, nil)
	}
	if err != nil {
		return 0, err
	}

	return p.ID, nil
}
