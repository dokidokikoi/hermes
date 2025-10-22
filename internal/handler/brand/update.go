package brand

import (
	"context"
	"errors"
	"izumi/db/data"
	"izumi/model"
	"strings"

	comm_errors "github.com/dokidokikoi/go-common/errors"
)

func (h Handler) Update(ctx context.Context, input *model.Brand) (any, error) {
	input.Name = strings.TrimSpace(input.Name)
	if err := data.GetDataFactory().Brand().Update(ctx, input, nil); err != nil {
		if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
			return nil, err
		}
	}

	return nil, nil
}
