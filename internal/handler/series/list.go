package series

import (
	"context"
	"izumi/db"
	"izumi/internal/handler"
	"izumi/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/middleware"
)

type ListRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type ListResponse struct {
	List  []*handler.SeriesVo `json:"list"`
	Total int64               `json:"total"`
}

func (h Handler) List(ctx context.Context, req *ListRequest, op *middleware.PreHandleOptions) (*ListResponse, error) {
	opt := &meta.ListOption{
		Order: "id desc",
		GetOption: meta.GetOption{
			Select: []string{"id", "name", "created_at", "count(game_id) games"},
			Join: []*meta.Join{
				{
					Method:         meta.LEFT_JOIN,
					Table:          model.Series{}.TableName(),
					TableField:     "id",
					JoinTable:      model.GameSeries{}.TableName(),
					JoinTableField: "series_id",
				},
			},
			Group: meta.Group{
				Fields: []string{"id"},
			},
		},
	}
	if req.Page > 0 && req.PageSize > 0 {
		opt.Page = req.Page
		opt.PageSize = req.PageSize
	}

	list, err := db.GetStore().Series().List(ctx, &model.Series{}, opt)
	if err != nil {
		return nil, err
	}

	var vos []*handler.SeriesVo
	for _, series := range list {
		vos = append(vos, &handler.SeriesVo{
			ID:        series.ID,
			Name:      series.Name,
			Games:     series.Games,
			CreatedAt: series.CreatedAt,
		})
	}

	return &ListResponse{
		List:  vos,
		Total: int64(len(list)),
	}, nil
}
