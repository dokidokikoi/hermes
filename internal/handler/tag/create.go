package tag

import (
	"context"
	"errors"
	"hermes/db/data"
	"hermes/model"

	comm_errors "github.com/dokidokikoi/go-common/errors"
)

func (h Handler) Create(ctx context.Context, input *model.Tag) (uint, *comm_errors.APIError) {
	if err := data.GetDataFactory().Tag().Create(ctx, input, nil); err != nil {
		if !errors.Is(err, comm_errors.ErrNameDuplicate) {
			return 0, comm_errors.Wrap(comm_errors.ApiErrSystemErr, err)
		}
	}

	return input.ID, nil
}
