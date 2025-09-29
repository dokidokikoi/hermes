package series

import (
	"context"
	"errors"
	"hermes/db/data"
	"hermes/model"
	"strings"

	"gorm.io/gorm"
)

func (h Handler) Create(ctx context.Context, input *model.Series) (uint, error) {
	input.Name = strings.TrimSpace(input.Name)
	if err := data.GetDataFactory().Series().Create(ctx, input, nil); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c, err := data.GetDataFactory().Series().Get(ctx, &model.Series{Name: input.Name}, nil)
			if err != nil {
				return 0, err
			}
			return c.ID, nil
		}
		return 0, err
	}

	return input.ID, nil
}
