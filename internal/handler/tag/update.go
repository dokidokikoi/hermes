package tag

import (
	"context"
	"errors"
	"izumi/db"
	"izumi/model"

	"github.com/abadojack/whatlanggo"
	comm_errors "github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/middleware"
)

func (h Handler) Update(ctx context.Context, input *model.Tag, op *middleware.PreHandleOptions) (any, error) {
	lang := whatlanggo.DetectLang(input.Name)
	input.Lang = lang.Iso6391()
	if err := db.GetStore().Tag().Update(ctx, input, nil); err != nil {
		if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
			return nil, err
		}
	}

	return nil, nil
}
