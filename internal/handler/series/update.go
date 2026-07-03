package series

import (
	"context"
	"errors"
	"izumi/db"
	"izumi/model"
	"strings"

	comm_errors "github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/middleware"
)

func (h Handler) Update(ctx context.Context, input *model.Series, op *middleware.PreHandleOptions) (any, error) {
	input.Name = strings.TrimSpace(input.Name)
	if err := db.GetStore().Series().Update(ctx, input, nil); err != nil {
		if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
			return nil, err
		}
	}

	return nil, nil
}
