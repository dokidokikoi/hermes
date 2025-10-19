package tag

import (
	"context"
	"errors"
	"hermes/db/data"
	"hermes/model"

	comm_errors "github.com/dokidokikoi/go-common/errors"
	meta "github.com/dokidokikoi/go-common/meta/option"
)

func (h Handler) Create(ctx context.Context, input *model.Tag) (uint, error) {
	if err := data.GetDataFactory().Tag().Create(ctx, input, nil); err != nil {
		if errors.Is(err, comm_errors.ErrNameDuplicate) {
			tag, err := data.GetDataFactory().Tag().Get(ctx, &model.Tag{Name: input.Name}, &meta.GetOption{})
			if err != nil {
				return 0, err
			}
			if tag == nil {
				return 0, nil
			}
			return tag.ID, nil
		}
		return 0, err
	}

	return input.ID, nil
}
