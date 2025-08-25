package developer

import (
	"context"
	"hermes/db/data"
	"hermes/model"

	"github.com/dokidokikoi/go-common/errors"
)

func (h Handler) Create(ctx context.Context, input *model.Developer) (uint, *errors.APIError) {
	if err := data.GetDataFactory().Developer().Create(ctx, input, nil); err != nil {
		return 0, errors.Wrap(errors.ApiErrSystemErr, err)
	}

	return input.ID, nil
}
