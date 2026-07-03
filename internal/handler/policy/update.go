package policy

import (
	"context"
	"encoding/json"
	"izumi/db"
	"izumi/internal/handler"
	"izumi/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/middleware"
)

func (h Handler) Update(ctx context.Context, input *handler.UpdateProxyReq, op *middleware.PreHandleOptions) (any, error) {

	switch input.Key {
	case model.SystemPolicy{}.Key():
		t := new(model.SystemPolicy)
		err := json.Unmarshal([]byte(input.Policy), t)
		if err != nil {
			return nil, err
		}
		h.srv.Policy().SystemPolicyEffect(ctx, t)
	case model.ScraperPolicy{}.Key():
		t := new(model.ScraperPolicy)
		err := json.Unmarshal([]byte(input.Policy), t)
		if err != nil {
			return nil, err
		}
		h.srv.Policy().ScraperPolicyEffect(ctx, t)
	}

	err := db.GetStore().Policy().UpdateByWhere(ctx, &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "key",
				Operator: meta.EQUAL,
				Value:    input.Key,
			},
		},
	}, &model.Policy{Policy: input.Policy}, nil)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
