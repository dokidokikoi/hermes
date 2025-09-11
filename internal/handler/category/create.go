package category

import (
	"context"
	"hermes/db/data"
	"hermes/model"
)

func (h Handler) Create(ctx context.Context, input *model.Category) (uint, error) {
	if err := data.GetDataFactory().Category().Create(ctx, input, nil); err != nil {
		return 0, err
	}

	return input.ID, nil
}
