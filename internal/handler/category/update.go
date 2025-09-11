package category

import (
	"context"
	"errors"
	"hermes/db/data"
	"hermes/model"

	comm_errors "github.com/dokidokikoi/go-common/errors"
)

func (h Handler) Update(ctx context.Context, req *model.Category) (any, error) {
	if err := data.GetDataFactory().Category().Update(ctx, req, nil); err != nil {
		if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
			return nil, err
		}
	}

	return nil, nil
}
