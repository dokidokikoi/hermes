package tag

import (
	"context"
	"errors"
	"hermes/db/data"
	"hermes/model"

	comm_errors "github.com/dokidokikoi/go-common/errors"
)

func (h Handler) Create(ctx context.Context, input *model.Tag) (uint, error) {
	if err := data.GetDataFactory().Tag().Create(ctx, input, nil); err != nil {
		if !errors.Is(err, comm_errors.ErrNameDuplicate) {
			return 0, err
		}
	}

	return input.ID, nil
}
