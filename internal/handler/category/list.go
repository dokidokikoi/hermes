package category

import (
	"context"
	"hermes/db/data"
	"hermes/model"

	"github.com/dokidokikoi/go-common/errors"
	meta "github.com/dokidokikoi/go-common/meta/option"
)

func (h Handler) List(ctx context.Context, req any) ([]*model.Category, *errors.APIError) {
	list, err := data.GetDataFactory().Category().List(ctx, &model.Category{}, &meta.ListOption{Order: "created_at desc"})
	if err != nil {
		return nil, errors.Wrap(errors.ApiErrSystemErr, err)
	}
	return list, nil
}
