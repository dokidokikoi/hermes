package category

import (
	"context"
	"errors"
	"hermes/db/data"
	"hermes/model"
	"strings"

	"gorm.io/gorm"
)

func (h Handler) Create(ctx context.Context, input *model.Category) (uint, error) {
	input.Name = strings.TrimSpace(strings.ToUpper(input.Name))
	if err := data.GetDataFactory().Category().Create(ctx, input, nil); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c, err := data.GetDataFactory().Category().Get(ctx, &model.Category{Name: input.Name}, nil)
			if err != nil {
				return 0, err
			}
			return c.ID, nil
		}
		return 0, err
	}

	return input.ID, nil
}
