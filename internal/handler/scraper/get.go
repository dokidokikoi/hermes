package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"
	"hermes/scraper"

	"github.com/dokidokikoi/go-common/errors"
	"github.com/tidwall/gjson"
)

func (h Handler) Get(ctx context.Context, input *handler.ScraperGetReq) (any, *errors.APIError) {
	list, err := data.GetDataFactory().Task().List(ctx, &model.Task{RequestID: input.RequestID, Status: model.TaskStatusSucceed}, nil)
	if err != nil {
		return nil, errors.Wrap(errors.ApiErrSystemErr, err)
	}

	if len(list) > 0 {
		if list[0].Type == model.TaskTypeDetail {
			var res []scraper.GameItem
			for _, l := range list {
				item := scraper.GameItem{}
				err := json.Unmarshal([]byte(l.Result), &item)
				if err != nil {
					return nil, errors.Wrap(errors.ApiErrSystemErr, err)
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
					return nil, errors.Wrap(errors.ApiErrSystemErr, err)
				}

				res[fmt.Sprintf("%s - %s - %d", l.ScraperName, gjson.Get(l.Param, "keyword").String(), gjson.Get(l.Param, "page").Int())] = items
			}
			return res, nil
		}
	}

	return nil, nil
}
