package series

import (
	"context"
	"hermes/db/data"
	"hermes/model"
)

func (h Handler) Create(ctx context.Context, input *model.Series) (uint, error) {
	if err := data.GetDataFactory().Series().Create(ctx, input, nil); err != nil {
		return 0, err
	}

	return input.ID, nil
}
