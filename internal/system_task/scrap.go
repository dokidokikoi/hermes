package systemtask

import (
	"context"
	"encoding/json"
	"izumi/constant"
	"izumi/db/data"
	"izumi/internal/handler"
	"izumi/internal/service"
	"izumi/model"
	"izumi/scraper"
	"izumi/scraper/bangumi"
	"izumi/scraper/twodfan"
	"izumi/scraper/vndb"
	"izumi/tools"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/notice"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func StartAutoScrap() {
	ts, err := data.GetDataFactory().SystemTask().List(context.Background(), &model.SystemTask{
		Type:  model.SystemTaskTypeScrap,
		State: model.SystemTaskStateRunning,
	}, &meta.ListOption{Order: "id desc"})
	if err != nil {
		zaplog.L().Error("system download task error", zap.Error(err))
		return
	}
	if len(ts) == 0 {
		return
	}
	err = data.GetDataFactory().SystemTask().UpdateByWhere(
		context.Background(),
		&meta.WhereNode{
			Conditions: []*meta.Condition{
				{
					Field:    "id",
					Operator: meta.NOTEQUAL,
					Value:    ts[0].ID,
				},
			},
		},
		&model.SystemTask{
			State: model.SystemTaskStateCanceled,
		},
		nil,
	)
	if err != nil {
		zaplog.L().Error("system update error", zap.Error(err))
	}

	AutoScrap(ts[0])

}

func AutoScrap(t *model.SystemTask) {
	go func() {
		var result string
		var e error
		defer func() {
			var state = model.SystemTaskStateDone
			if e != nil {
				zaplog.L().Error("system scrap task error", zap.Error(e))
				state = model.SystemTaskStateFailed
			}
			err := data.GetDataFactory().SystemTask().Update(context.Background(), &model.SystemTask{
				ID:     t.ID,
				State:  state,
				Result: result,
			}, nil)
			if err != nil {
				zaplog.L().Error("system update error", zap.Error(err))
			}

			err = notice.HubIns.SendBroadcast(constant.TOPIC_SCRAPER, notice.NoticeResponse{
				Rid:     uuid.NewString(),
				Event:   constant.EVENT_SCRAPER_AUTOSCRAP,
				Success: e == nil,
				Message: result,
			})
			if err != nil {
				zaplog.L().Error("send notify error", zap.Error(err))
			}
		}()

		flag := false
		for _, obj := range t.Param.ScrapObjs {
			if obj.Name == vndb.Name {
				flag = true
			}
		}
		if !flag {
			result = "auto scrap need vndb item"
			e = errors.New(result)
			return
		}

		requestID := uuid.NewString()
		wait, err := service.NewScrap(data.GetDataFactory()).Scrap(context.Background(), requestID, t.Param.ScrapObjs)
		if err != nil {
			e = err
			result = err.Error()
			return
		}
		wait.Wait()

		list, err := data.GetDataFactory().Task().List(context.Background(), &model.Task{
			RequestID: requestID,
			Status:    model.TaskStatusSucceed,
		}, nil)
		if err != nil {
			e = err
			result = err.Error()
			return
		}
		itemM := map[string][]scraper.GameItem{}
		for _, l := range list {
			item := scraper.GameItem{}
			err := json.Unmarshal([]byte(l.Result), &item)
			if err != nil {
				zaplog.L().Error("system scrap task error", zap.Error(err))
				continue
			}
			itemM[item.ScraperName] = append(itemM[item.ScraperName], item)
		}

		items, ok := itemM[vndb.Name]
		if !ok || len(items) == 0 {
			result = "vndb item scrap failed"
			return
		}
		delete(itemM, vndb.Name)

		gVo := items[0].GameVo
		aliasM := map[string]struct{}{}
		tagM := map[string]struct{}{}
		characterM := map[string]map[int]struct{}{}
		staffM := map[string]map[int]struct{}{}
		for _, a := range gVo.Alias {
			aliasM[a] = struct{}{}
		}
		for _, t := range gVo.Tags {
			tagM[t.Name] = struct{}{}
		}
		for i, c := range gVo.Characters {
			name := tools.ToLowerNoSpace(c.Name)
			if _, ok := characterM[name]; !ok {
				characterM[name] = map[int]struct{}{}
			}
			characterM[name][i] = struct{}{}
			for _, n := range c.Alias {
				n := tools.ToLowerNoSpace(n)
				if _, ok := characterM[n]; !ok {
					characterM[n] = map[int]struct{}{}
				}
				characterM[n][i] = struct{}{}
			}
		}
		for i, s := range gVo.Staff {
			name := tools.ToLowerNoSpace(s.Name)
			if _, ok := staffM[name]; !ok {
				staffM[name] = map[int]struct{}{}
			}
			staffM[name][i] = struct{}{}
			for _, n := range s.Alias {
				n := tools.ToLowerNoSpace(n)
				if _, ok := staffM[n]; !ok {
					staffM[n] = map[int]struct{}{}
				}
				staffM[n][i] = struct{}{}

			}
		}
		for name, item := range itemM {
			for _, i := range item {
				for _, a := range i.Alias {
					if _, ok := aliasM[a]; !ok {
						gVo.Alias = append(gVo.Alias, a)
					}
				}
				for _, t := range i.Tags {
					if _, ok := tagM[t.Name]; !ok {
						tagM[t.Name] = struct{}{}
						gVo.Tags = append(gVo.Tags, t)
					}
				}
				gVo.Images = append(gVo.Images, i.Images...)
				gVo.Links = append(gVo.Links, i.Links...)
				if len(gVo.Brands) == 0 {
					gVo.Brands = append(gVo.Brands, i.Brands...)
				}
				if gVo.Category == nil || gVo.Category.Name == "" {
					gVo.Category = i.Category
				}
				if gVo.IssueDate.IsZero() {
					gVo.IssueDate = i.IssueDate
				}
				if gVo.Story == "" {
					gVo.Story = i.Story
				}
				if gVo.DlCode == "" {
					gVo.DlCode = i.DlCode
				}
				if gVo.JanCode == "" {
					gVo.JanCode = i.JanCode
				}
				if gVo.Price == "" {
					gVo.Price = i.Price
				}
			}
			switch name {
			case bangumi.Name:
				for _, i := range item {
					for _, c := range i.Characters {
						name := tools.ToLowerNoSpace(c.Name)
						if m, ok := characterM[name]; ok && len(m) == 1 {
							for k := range m {
								gVo.Characters[k].Images = append(gVo.Characters[k].Images, c.Images...)
								if c.Cover != "" {
									gVo.Characters[k].Images = append(gVo.Characters[k].Images, c.Cover)
								}
								if gVo.Characters[k].Gender == model.UnKnown {
									gVo.Characters[k].Gender = c.Gender
								}
								if gVo.Characters[k].Rlation == "" {
									gVo.Characters[k].Rlation = c.Rlation
								}
								if c.Summary != "" {
									gVo.Characters[k].Summary = c.Summary
								}
							}
						} else {
							_, ok := characterM[name]
							if !ok {
								characterM[name] = map[int]struct{}{}
							}
							characterM[name][len(gVo.Characters)] = struct{}{}
							gVo.Characters = append(gVo.Characters, c)
						}
					}
					for _, s := range i.Staff {
						name := tools.ToLowerNoSpace(s.Name)
						if m, ok := staffM[name]; ok && len(m) == 1 {
							for k := range m {
								gVo.Staff[k].Images = append(gVo.Staff[k].Images, s.Images...)
								if s.Cover != "" {
									gVo.Staff[k].Images = append(gVo.Staff[k].Images, s.Cover)
								}
								gVo.Staff[k].Links = append(gVo.Staff[k].Links, s.Links...)
								if gVo.Staff[k].Gender == model.UnKnown {
									gVo.Staff[k].Gender = s.Gender
								}
								if len(gVo.Staff[k].Relation) == 0 {
									gVo.Staff[k].Relation = s.Relation
								}
								if s.Summary != "" {
									gVo.Staff[k].Summary = s.Summary
								}

							}
						} else {
							_, ok := staffM[name]
							if !ok {
								staffM[name] = map[int]struct{}{}
							}
							staffM[name][len(gVo.Staff)] = struct{}{}
							gVo.Staff = append(gVo.Staff, s)
						}
					}
				}
			case twodfan.Name:
				for _, i := range item {
					if i.Story != "" {
						gVo.Story = i.Story
					}
				}
				fallthrough
			default:
				for _, i := range item {
					for _, c := range i.Characters {
						name := tools.ToLowerNoSpace(c.Name)
						if m, ok := characterM[name]; ok && len(m) == 1 {
							for k := range m {
								gVo.Characters[k].Images = append(gVo.Characters[k].Images, c.Images...)
								if c.Cover != "" {
									gVo.Characters[k].Images = append(gVo.Characters[k].Images, c.Cover)
								}
								if gVo.Characters[k].Gender == model.UnKnown {
									gVo.Characters[k].Gender = c.Gender
								}
								if gVo.Characters[k].Rlation == "" {
									gVo.Characters[k].Rlation = c.Rlation
								}
								if gVo.Characters[k].Summary != "" {
									gVo.Characters[k].Summary = c.Summary
								}
							}
						} else {
							_, ok := characterM[name]
							if !ok {
								characterM[name] = map[int]struct{}{}
							}
							characterM[name][len(gVo.Characters)] = struct{}{}
							gVo.Characters = append(gVo.Characters, c)
						}
					}
					for _, s := range i.Staff {
						name := tools.ToLowerNoSpace(s.Name)
						if m, ok := staffM[name]; ok && len(m) == 1 {
							for k := range m {
								gVo.Staff[k].Images = append(gVo.Staff[k].Images, s.Images...)
								if s.Cover != "" {
									gVo.Staff[k].Images = append(gVo.Staff[k].Images, s.Cover)
								}
								gVo.Staff[k].Links = append(gVo.Staff[k].Links, s.Links...)
								if gVo.Staff[k].Gender == model.UnKnown {
									gVo.Staff[k].Gender = s.Gender
								}
								if len(gVo.Staff[k].Relation) == 0 {
									gVo.Staff[k].Relation = s.Relation
								}
								if gVo.Staff[k].Summary != "" {
									gVo.Staff[k].Summary = s.Summary
								}

							}
						} else {
							_, ok := staffM[name]
							if !ok {
								staffM[name] = map[int]struct{}{}
							}
							staffM[name][len(gVo.Staff)] = struct{}{}
							gVo.Staff = append(gVo.Staff, s)
						}
					}
				}
			}
		}

		staffs := []handler.StaffVo{}
		for _, staff := range gVo.Staff {
			for _, r := range staff.Relation {
				if r != model.PRelationUnknown {
					staffs = append(staffs, staff)
					break
				}
			}
		}
		gVo.Staff = staffs

		err = service.NewGame(data.GetDataFactory()).UpsertFull(context.Background(), &gVo, &model.GameInstance{
			Path:    t.Param.Path,
			Version: t.Param.Version,
		})
		if err != nil {
			e = err
			result = err.Error()
			return
		}
	}()
}
