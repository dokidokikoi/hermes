package developer

import (
	"context"
	"errors"
	"hermes/db/data"
	"hermes/model"
	"strings"

	"gorm.io/gorm"
)

func (h Handler) Create(ctx context.Context, input *model.Developer) (uint, error) {
	input.Name = strings.TrimSpace(input.Name)
	if err := data.GetDataFactory().Developer().Create(ctx, input, nil); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c, err := data.GetDataFactory().Developer().Get(ctx, &model.Developer{Name: input.Name}, nil)
			if err != nil {
				return 0, err
			}
			return c.ID, nil
		}
		return 0, err
	}

	return input.ID, nil
}
