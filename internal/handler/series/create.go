package series

import (
	"context"
	"errors"
	"hermes/db/data"
	"hermes/model"

	comm_errors "github.com/dokidokikoi/go-common/errors"
	meta "github.com/dokidokikoi/go-common/meta/option"
)

func (h Handler) Create(ctx context.Context, input *model.Series) (uint, error) {
	if err := data.GetDataFactory().Series().Create(ctx, input, &meta.CreateOption{}); err != nil {
		if errors.Is(err, comm_errors.ErrNameDuplicate) {
			series, err := data.GetDataFactory().Series().Get(ctx, &model.Series{Name: input.Name}, &meta.GetOption{})
			if err != nil {
				return 0, err
			}
			if series == nil {
				return 0, nil
			}
			return series.ID, nil
		}
		return 0, err
	}

	return input.ID, nil
}
