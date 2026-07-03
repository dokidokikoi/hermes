package brand

import (
	"context"
	"errors"
	"izumi/db"
	"izumi/model"
	"strings"

	"github.com/dokidokikoi/go-common/middleware"
	"gorm.io/gorm"
)

func (h Handler) Create(ctx context.Context, input *model.Brand, op *middleware.PreHandleOptions) (uint, error) {
	input.Name = strings.TrimSpace(input.Name)
	if err := db.GetStore().Brand().Create(ctx, input, nil); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c, err := db.GetStore().Brand().Get(ctx, &model.Brand{Name: input.Name}, nil)
			if err != nil {
				return 0, err
			}
			return c.ID, nil
		}
		return 0, err
	}

	return input.ID, nil
}
