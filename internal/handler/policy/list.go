package policy

import (
	"context"
	"hermes/db/data"
	"hermes/model"

	"github.com/dokidokikoi/go-common/errors"
	"github.com/gin-gonic/gin"
)

type ListRequest struct{}

func (h Handler) List(ctx context.Context, req *ListRequest) (gin.H, *errors.APIError) {
	ps, err := data.GetDataFactory().Policy().List(ctx, &model.Policy{}, nil)
	if err != nil {
		return nil, errors.Wrap(errors.ApiErrSystemErr, err)
	}
	res := gin.H{}
	for _, p := range ps {
		res[p.Key] = p.Policy
	}

	return res, nil
}
