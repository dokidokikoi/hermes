package character

import (
	"context"
	"izumi/db"
	"izumi/model"
	"strings"

	"github.com/dokidokikoi/go-common/middleware"
)

func (h Handler) Create(ctx context.Context, input *model.Character, op *middleware.PreHandleOptions) (uint, error) {
	input.Name = strings.TrimSpace(input.Name)
	if err := db.GetStore().Character().Create(ctx, input, nil); err != nil {
		return 0, err
	}

	return input.ID, nil
}
