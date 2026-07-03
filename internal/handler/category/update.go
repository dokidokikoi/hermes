package category

import (
	"context"
	"errors"
	"izumi/db"
	"izumi/model"
	"strings"

	comm_errors "github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/middleware"
)

func (h Handler) Update(ctx context.Context, req *model.Category, op *middleware.PreHandleOptions) (any, error) {
	req.Name = strings.TrimSpace(strings.ToUpper(req.Name))
	if err := db.GetStore().Category().Update(ctx, req, nil); err != nil {
		if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
			return nil, err
		}
	}

	return nil, nil
}
