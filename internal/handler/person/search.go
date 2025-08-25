package person

import (
	"context"
	"hermes/internal/handler"
	"hermes/internal/service"

	"github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/query"
)

type SearchResponse struct {
	List  []handler.StaffVo `json:"list"`
	Total int64             `json:"total"`
}

func (h Handler) Search(ctx context.Context, input *handler.PersonListReq) (*SearchResponse, *errors.APIError) {
	var q = &query.PageQuery{
		Page:     input.Page,
		PageSize: input.PageSize,
		Order:    input.OrderBy,
	}
	total, list, err := h.srv.Person().Search(ctx, *input, q.GetListOption(), service.PersonBasicSearchNode...)
	if err != nil {
		return nil, errors.Wrap(errors.ApiErrSystemErr, err)
	}

	vos := make([]handler.StaffVo, len(list))
	for i := range list {
		vos[i] = handler.StaffVo{
			ID:        list[i].ID,
			Name:      list[i].Name,
			Alias:     list[i].Alias,
			Cover:     list[i].Cover,
			Images:    list[i].Images,
			Tags:      list[i].Tags,
			Summary:   list[i].Summary,
			Gender:    list[i].Gender,
			CreatedAt: list[i].CreatedAt,
		}
	}

	return &SearchResponse{
		List:  vos,
		Total: total,
	}, nil
}
