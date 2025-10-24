package policy

import (
	"context"
	"izumi/db/data"
	"izumi/model"

	"github.com/dokidokikoi/go-common/middleware"
	"github.com/gin-gonic/gin"
)

type ListRequest struct{}

func (h Handler) List(ctx context.Context, req *ListRequest, op *middleware.PreHandleOptions) (gin.H, error) {
	ps, err := data.GetDataFactory().Policy().List(ctx, &model.Policy{}, nil)
	if err != nil {
		return nil, err
	}
	res := gin.H{}
	for _, p := range ps {
		res[p.Key] = p.Policy
	}

	return res, nil
}
