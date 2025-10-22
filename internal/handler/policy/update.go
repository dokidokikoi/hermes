package policy

import (
	"context"
	"encoding/json"
	"izumi/db/data"
	"izumi/internal/handler"
	"izumi/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

func (h Handler) Update(ctx context.Context, input *handler.UpdateProxyReq) (any, error) {
	var policy any

	switch input.Key {
	case model.SystemPolicy{}.Key():
		t := new(model.SystemPolicy)
		err := json.Unmarshal([]byte(input.Policy), t)
		if err != nil {
			return nil, err
		}
		policy = t
		h.srv.Policy().SystemPolicyEffect(ctx, t)
	case model.ScraperPolicy{}.Key():
		t := new(model.ScraperPolicy)
		err := json.Unmarshal([]byte(input.Policy), t)
		if err != nil {
			return nil, err
		}
		policy = t
		h.srv.Policy().ScraperPolicyEffect(ctx, t)
	case model.LanguagePolicy{}.Key():
		t := new(model.LanguagePolicy)
		err := json.Unmarshal([]byte(input.Policy), t)
		if err != nil {
			return nil, err
		}
		policy = t
	case model.PlatformPolicy{}.Key():
		t := new(model.PlatformPolicy)
		err := json.Unmarshal([]byte(input.Policy), t)
		if err != nil {
			return nil, err
		}
		policy = t
	}
	d, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	err = data.GetDataFactory().Policy().UpdateByWhere(ctx, &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "key",
				Operator: meta.EQUAL,
				Value:    input.Key,
			},
		},
	}, &model.Policy{Policy: string(d)}, nil)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
