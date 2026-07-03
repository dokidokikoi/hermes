package character

import (
	"context"
	"errors"
	"izumi/db"
	"izumi/model"

	"github.com/dokidokikoi/go-common/middleware"
)

type DelRequest struct {
	IDs []uint `json:"ids"`
}

func (h Handler) Del(ctx context.Context, req *DelRequest, op *middleware.PreHandleOptions) (any, error) {
	CharacterIDs := []*model.Character{}
	for _, id := range req.IDs {
		CharacterIDs = append(CharacterIDs, &model.Character{
			ID: id,
		})
	}

	errs := db.GetStore().Character().DeleteCollection(ctx, CharacterIDs, nil)
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return nil, nil
}
