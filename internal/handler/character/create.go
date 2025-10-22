package character

import (
	"context"
	"izumi/db/data"
	"izumi/model"
	"strings"

	"github.com/google/uuid"
)

func (h Handler) Create(ctx context.Context, input *model.Character) (uint, error) {
	input.UUID = uuid.NewString()
	input.Name = strings.TrimSpace(input.Name)
	if err := data.GetDataFactory().Character().Create(ctx, input, nil); err != nil {
		return 0, err
	}

	return input.ID, nil
}
