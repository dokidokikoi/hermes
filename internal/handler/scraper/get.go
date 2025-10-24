package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"izumi/db/data"
	"izumi/internal/handler"
	"izumi/model"
	"izumi/scraper"

	"github.com/dokidokikoi/go-common/middleware"
	"github.com/tidwall/gjson"
)

func (h Handler) Get(ctx context.Context, input *handler.ScraperGetReq, op *middleware.PreHandleOptions) (any, error) {
	list, err := data.GetDataFactory().Task().List(ctx, &model.Task{RequestID: input.RequestID, Status: model.TaskStatusSucceed}, nil)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		if list[0].Type == model.TaskTypeDetail {
			var res []scraper.GameItem
			for _, l := range list {
				item := scraper.GameItem{}
				err := json.Unmarshal([]byte(l.Result), &item)
				if err != nil {
					return nil, err
				}
				res = append(res, item)
			}
			return res, nil
		} else if list[0].Type == model.TaskTypeSearch {
			res := map[string][]scraper.SearchItem{}
			for _, l := range list {
				items := []scraper.SearchItem{}
				err := json.Unmarshal([]byte(l.Result), &items)
				if err != nil {
					return nil, err
				}

				res[fmt.Sprintf("%s - %s - %d", l.ScraperName, gjson.Get(l.Param, "keyword").String(), gjson.Get(l.Param, "page").Int())] = items
			}
			return res, nil
		}
	}

	return nil, nil
}
