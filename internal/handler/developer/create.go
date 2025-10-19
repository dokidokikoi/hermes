package developer

import (
	"context"
	"errors"
	"hermes/db/data"
	"hermes/model"

	comm_errors "github.com/dokidokikoi/go-common/errors"
	meta "github.com/dokidokikoi/go-common/meta/option"
)

func (h Handler) Create(ctx context.Context, input *model.Developer) (uint, error) {
	if err := data.GetDataFactory().Developer().Create(ctx, input, nil); err != nil {
		if errors.Is(err, comm_errors.ErrNameDuplicate) {
			developer, err := data.GetDataFactory().Developer().Get(ctx, &model.Developer{Name: input.Name}, &meta.GetOption{})
			if err != nil {
				return 0, err
			}
			if developer == nil {
				return 0, nil
			}
			return developer.ID, nil
		}
		return 0, err
	}

	return input.ID, nil
}
