package category

import (
	"context"
	"errors"
	"hermes/db/data"
	"hermes/model"
	"strings"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"gorm.io/gorm"
)

func (h Handler) Create(ctx context.Context, input *model.Category) (uint, error) {
	input.Name = strings.ToUpper(input.Name)
	if err := data.GetDataFactory().Category().Create(ctx, input, nil); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			category, err := data.GetDataFactory().Category().Get(ctx, &model.Category{Name: input.Name}, &meta.GetOption{})
			if err != nil {
				return 0, err
			}
			if category == nil {
				return 0, nil
			}
			return category.ID, nil
		}
		return 0, err
	}

	return input.ID, nil
}
