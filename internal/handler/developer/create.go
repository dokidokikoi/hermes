package developer

import (
	"context"
	"hermes/db/data"
	"hermes/model"

	"github.com/google/uuid"
)

func (h Handler) Create(ctx context.Context, input *model.Developer) (uint, error) {
	input.UUID = uuid.NewString()
	if err := data.GetDataFactory().Developer().Create(ctx, input, nil); err != nil {
		return 0, err
	}

	return input.ID, nil
}
