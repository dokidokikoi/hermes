package series

import (
	"context"
	"errors"
	"hermes/db/data"
	"hermes/model"

	comm_errors "github.com/dokidokikoi/go-common/errors"
)

func (h Handler) Update(ctx context.Context, input *model.Series) (any, *comm_errors.APIError) {
	if err := data.GetDataFactory().Series().Update(ctx, input, nil); err != nil {
		if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
			return nil, comm_errors.Wrap(comm_errors.ApiErrSystemErr, err)
		}
	}

	return nil, nil
}
