package category

import (
	"context"
	"izumi/db"
	"izumi/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/middleware"
)

func (h Handler) List(ctx context.Context, req *struct{}, op *middleware.PreHandleOptions) ([]*model.Category, error) {
	list, err := db.GetStore().Category().List(ctx, &model.Category{}, &meta.ListOption{Order: "created_at desc"})
	if err != nil {
		return nil, err
	}
	return list, nil
}
