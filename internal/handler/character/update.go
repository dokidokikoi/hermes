package character

import (
	"context"
	"errors"
	"izumi/db"
	"izumi/model"
	"strings"

	comm_errors "github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/middleware"
)

func (h Handler) Update(ctx context.Context, input *model.Character, op *middleware.PreHandleOptions) (any, error) {
	tx := db.GetStore().Transaction().Begin()

	input.Name = strings.TrimSpace(input.Name)
	if err := tx.Character().Update(ctx, input, nil); err != nil {
		if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
			tx.Transaction().Rollback()
			return nil, err
		}
	}

	return nil, nil
}
