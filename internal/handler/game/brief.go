package game

import (
	"context"
	"izumi/db/data"
	"izumi/internal/handler"
	"izumi/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/middleware"
	"github.com/dokidokikoi/go-common/tools"
)

type GetBriefRequest struct {
	IDs uint `json:"ids"`
}

func (h Handler) GetBrief(ctx context.Context, req *GetBriefRequest, op *middleware.PreHandleOptions) (resp []*handler.GameVo, e error) {
	res, err := data.GetDataFactory().Game().ListComplex(ctx, &model.Game{}, &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "id",
				Operator: meta.IN,
				Value:    req.IDs,
			},
		},
	}, &meta.ListOption{GetOption: meta.GetOption{Select: []string{"ID", "Name", "Cover"}}})
	if err != nil {
		return nil, err
	}

	for _, v := range res {
		resp = append(resp, tools.GetPtr(handler.Game2Vo(*v, nil, nil)))
	}

	return
}
