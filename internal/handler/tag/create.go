package tag

import (
	"context"
	"errors"
	"izumi/db/data"
	"izumi/model"

	"github.com/abadojack/whatlanggo"
	"github.com/dokidokikoi/go-common/middleware"
	"gorm.io/gorm"
)

func (h Handler) Create(ctx context.Context, input *model.Tag, op *middleware.PreHandleOptions) (uint, error) {
	lang := whatlanggo.DetectLang(input.Name)
	if lang.Iso6391() == "zh" || lang.Iso6391() == "ja" {
		input.Lang = lang.Iso6391()
	} else {
		input.Lang = "en"
	}
	if err := data.GetDataFactory().Tag().Create(ctx, input, nil); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c, err := data.GetDataFactory().Tag().Get(ctx, &model.Tag{Name: input.Name}, nil)
			if err != nil {
				return 0, err
			}
			return c.ID, nil
		}
		return 0, err

	}

	return input.ID, nil
}
