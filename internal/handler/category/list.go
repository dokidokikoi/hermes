package category

import (
	"context"
	"hermes/db/data"
	"hermes/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

func (h Handler) List(ctx context.Context, req *struct{}) ([]*model.Category, error) {
	list, err := data.GetDataFactory().Category().List(ctx, &model.Category{}, &meta.ListOption{Order: "created_at desc"})
	if err != nil {
		return nil, err
	}
	return list, nil
}
