package character

import (
	"context"
	"hermes/internal/handler"
	"hermes/internal/service"

	"github.com/dokidokikoi/go-common/query"
)

type SearchResponse struct {
	List  []handler.CharacterVo `json:"list"`
	Total int64                 `json:"total"`
}

func (h Handler) Search(ctx context.Context, input *handler.CharacterSearchReq) (*SearchResponse, error) {
	q := &query.PageQuery{
		Page:     input.Page,
		PageSize: input.PageSize,
		Order:    input.OrderBy,
	}
	total, list, err := h.srv.Character().Search(ctx, *input, q.GetListOption(), service.CharacterBasicSearchNode...)
	if err != nil {
		return nil, err
	}
	vos := make([]handler.CharacterVo, len(list))
	for i := range list {
		vos[i] = handler.CharacterVo{
			ID:      list[i].ID,
			Name:    list[i].Name,
			Alias:   list[i].Alias,
			Gender:  list[i].Gender,
			Summary: list[i].Summary,
			Cover:   list[i].Cover,
			Images:  list[i].Images,
			CV: handler.StaffVo{
				ID:        list[i].CV.ID,
				Name:      list[i].CV.Name,
				Alias:     list[i].CV.Alias,
				Cover:     list[i].CV.Cover,
				Images:    list[i].CV.Images,
				Tags:      list[i].CV.Tags,
				Summary:   list[i].CV.Summary,
				Gender:    list[i].CV.Gender,
				CreatedAt: list[i].CV.CreatedAt,
			},
			Tags:      list[i].Tags,
			CreatedAt: list[i].CreatedAt,
		}
	}

	return &SearchResponse{
		List:  vos,
		Total: total,
	}, nil
}
