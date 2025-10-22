package tag

import (
	"context"
	"errors"
	"izumi/db/data"
	"izumi/model"

	"gorm.io/gorm"
)

func (h Handler) Create(ctx context.Context, input *model.Tag) (uint, error) {
	if err := data.GetDataFactory().Tag().Create(ctx, input, nil); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c, err := data.GetDataFactory().Tag().Get(ctx, &model.Tag{Name: input.Name}, nil)
			if err != nil {
				return 0, err
			}
			return c.ID, nil
		}
		return 0, err

	}

	return input.ID, nil
}
