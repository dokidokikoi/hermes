package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"izumi/config"
	"izumi/db/data"
	"izumi/internal/handler"
	"izumi/model"
	"izumi/utils"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/tools"
	"go.uber.org/zap"
)

var GameBasicSearchNode = []GameWhereNodeFunc{
	GameWhereNodeKeyword,
	GameWhereNodeTag,
	GameWhereNodeCharacter,
	GameWhereNodeStaff,
	GameWhereNodeSeries,
	GameWhereNodeCategory,
	GameWhereNodeBrand,
	GameWhereNodeSizeRange,
	GameWhereNodeIssueDateRange,
	GameWhereNodeCreatedAtRange,
}

func GameWhereNodeKeyword(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	keyword := tools.TrimBlankChar(param.Keyword)
	if keyword != "" {
		node.Next = &meta.WhereNode{
			Conditions: []*meta.Condition{
				{
					Field:    "name",
					Operator: meta.LIKE,
					Value:    fmt.Sprintf("%%%s%%", keyword),
				},
				{
					Field:    "text",
					Operator: meta.LIKE,
					Value:    fmt.Sprintf("%%%s%%", keyword),
				},
			},
		}
		if param.FullText {
			node.Next.Conditions = append(node.Next.Conditions, []*meta.Condition{
				{
					Field:    "story",
					Operator: meta.LIKE,
					Value:    fmt.Sprintf("%%%s%%", keyword),
				},
				{
					Field:    "other_info",
					Operator: meta.LIKE,
					Value:    fmt.Sprintf("%%%s%%", keyword),
				},
			}...)
		}
		node = node.Next
	}
	return node, opt
}
func GameWhereNodeTag(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	if len(param.Tags) < 1 {
		return node, opt
	}

	defer func() {
		o = opt
		n = node
	}()
	tmpdb := data.GetDataFactory().Tag().ListComplexDB(ctx, &model.Tag{}, &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "name",
				Operator: meta.IN,
				Value:    param.Tags,
			},
		},
	}, nil)
	db := data.GetDataFactory().GameTag().ListDB(ctx, &model.GameTag{}, &meta.ListOption{
		GetOption: meta.GetOption{
			Join: []*meta.Join{
				{
					Method:          meta.INNER_JOIN,
					Table:           model.GameTag{}.TableName(),
					InnerQuery:      tmpdb,
					InnerQueryAlias: "t1",
					TableField:      "tag_id",
					JoinTableField:  "id",
				},
			},
			Select: []string{"DISTINCT \"game_tag\".\"game_id\""},
		},
	})
	opt.GetOption.Join = append(opt.GetOption.Join, &meta.Join{
		Method:          meta.INNER_JOIN,
		Table:           model.Game{}.TableName(),
		InnerQuery:      db,
		InnerQueryAlias: "game_tag",
		TableField:      "id",
		JoinTableField:  "game_id",
	})
	return
}
func GameWhereNodeCharacter(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	if param.Character == 0 {
		return node, opt
	}
	db := data.GetDataFactory().GameCharacter().ListDB(ctx, &model.GameCharacter{CharacterID: param.Character}, nil)
	opt.GetOption.Join = append(opt.GetOption.Join, &meta.Join{
		Method:          meta.INNER_JOIN,
		Table:           model.Game{}.TableName(),
		InnerQuery:      db,
		InnerQueryAlias: "game_character",
		TableField:      "id",
		JoinTableField:  "game_id",
	})
	return node, opt
}
func GameWhereNodeStaff(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption) {
	if param.Staff == 0 {
		return node, opt
	}
	db := data.GetDataFactory().GameStaff().ListDB(ctx, &model.GameStaff{PersonID: param.Staff}, nil)
	opt.GetOption.Join = append(opt.GetOption.Join, &meta.Join{
		Method:          meta.INNER_JOIN,
		Table:           model.Game{}.TableName(),
		InnerQuery:      db,
		InnerQueryAlias: "game_staff",
		TableField:      "id",
		JoinTableField:  "game_id",
	})
	return node, opt
}
func GameWhereNodeSeries(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption) {
	if param.Series == 0 {
		return node, opt
	}
	db := data.GetDataFactory().GameSeries().ListDB(ctx, &model.GameSeries{SeriesID: param.Series}, nil)
	opt.GetOption.Join = append(opt.GetOption.Join, &meta.Join{
		Method:          meta.INNER_JOIN,
		Table:           model.Game{}.TableName(),
		InnerQuery:      db,
		InnerQueryAlias: "game_series",
		TableField:      "id",
		JoinTableField:  "game_id",
	})
	return node, opt
}
func GameWhereNodeCategory(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	if param.Category == 0 {
		return node, opt
	}
	defer func() {
		n = node.Next
		o = opt
	}()

	node.Next = &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "category_id",
				Operator: meta.EQUAL,
				Value:    param.Category,
			},
		},
	}
	return
}
func GameWhereNodeBrand(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	if param.Brand == 0 {
		return node, opt
	}
	defer func() {
		n = node.Next
		o = opt
	}()

	node.Next = &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "brand_id",
				Operator: meta.EQUAL,
				Value:    param.Brand,
			},
		},
	}
	opt.GetOption.Join = append(opt.GetOption.Join, &meta.Join{
		Method:         meta.INNER_JOIN,
		Table:          model.Game{}.TableName(),
		JoinTable:      model.GameBrands{}.TableName(),
		TableField:     "id",
		JoinTableField: "game_id",
	})
	return
}
func GameWhereNodeSizeRange(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	if len(param.SizeRange) > 0 {
		node.Next = &meta.WhereNode{
			Conditions: []*meta.Condition{
				{
					Field:    "size",
					Operator: meta.GTE,
					Value:    param.SizeRange[0],
				},
			},
		}

		node = node.Next
		if len(param.SizeRange) > 1 {
			node.Next = &meta.WhereNode{
				Conditions: []*meta.Condition{
					{
						Field:    "size",
						Operator: meta.LTE,
						Value:    param.SizeRange[1],
					},
				},
			}
			node = node.Next
		}
	}

	return node, opt
}
func GameWhereNodeIssueDateRange(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption) {
	if len(param.IssueDateRange) > 0 {
		node.Next = &meta.WhereNode{
			Conditions: []*meta.Condition{
				{
					Field:    "issue_date",
					Operator: meta.GTE,
					Value:    param.IssueDateRange[0],
				},
			},
		}

		node = node.Next
		if len(param.IssueDateRange) > 1 {
			node.Next = &meta.WhereNode{
				Conditions: []*meta.Condition{
					{
						Field:    "issue_date",
						Operator: meta.LTE,
						Value:    param.IssueDateRange[1],
					},
				},
			}
			node = node.Next
		}
	}

	return node, opt
}
func GameWhereNodeCreatedAtRange(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption) {
	if len(param.CreatedAtRange) > 0 {
		node.Next = &meta.WhereNode{
			Conditions: []*meta.Condition{
				{
					Field:    "created_at",
					Operator: meta.GTE,
					Value:    param.CreatedAtRange[0],
				},
			},
		}

		node = node.Next
		if len(param.CreatedAtRange) > 1 {
			node.Next = &meta.WhereNode{
				Conditions: []*meta.Condition{
					{
						Field:    "created_at",
						Operator: meta.LTE,
						Value:    param.CreatedAtRange[1],
					},
				},
			}
			node = node.Next
		}
	}

	return node, opt
}

func cpGameAllImages(logger *zap.Logger, path string, gVo *handler.GameVo) error {
	path = filepath.Join(path, "images")
	_, err := os.ReadDir(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	root, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(gVo.Story)))
	if err != nil {
		logger.With(zap.Error(err)).Error("NewDocumentFromReader")
	} else {
		root.Find("img").Each(func(i int, s *goquery.Selection) {
			src, ok := s.Attr("src")
			if ok {
				src = strings.TrimPrefix(src, "/api/file/")
				err := tools.Cp(filepath.Join(config.DataDir, src), filepath.Join(path, src))
				if err != nil {
					logger.With(zap.Error(err)).Sugar().Errorf("tools.Cp(%s)", filepath.Join(path, src))
				}
			}
		})
	}
	if gVo.Cover != "" {
		err := tools.Cp(filepath.Join(config.DataDir, gVo.Cover), filepath.Join(path, gVo.Cover))
		if err != nil {
			logger.With(zap.Error(err)).Sugar().Errorf("tools.Cp(%s)", filepath.Join(path, gVo.Cover))
		}
	}
	for _, image := range gVo.Images {
		if image != "" {
			err := tools.Cp(filepath.Join(config.DataDir, image), filepath.Join(path, image))
			if err != nil {
				logger.With(zap.Error(err)).Sugar().Errorf("tools.Cp(%s)", filepath.Join(path, image))
			}
		}
	}
	for _, c := range gVo.Characters {
		if c.Cover != "" {
			err := tools.Cp(filepath.Join(config.DataDir, c.Cover), filepath.Join(path, c.Cover))
			if err != nil {
				logger.With(zap.Error(err)).Sugar().Errorf("tools.Cp(%s)", filepath.Join(path, c.Cover))
			}
		}
		for _, image := range c.Images {
			if image != "" {
				err := tools.Cp(filepath.Join(config.DataDir, image), filepath.Join(path, image))
				if err != nil {
					logger.With(zap.Error(err)).Sugar().Errorf("tools.Cp(%s)", filepath.Join(path, image))
				}
			}
		}
	}
	for _, s := range gVo.Staff {
		if s.Cover != "" {
			err := tools.Cp(filepath.Join(config.DataDir, s.Cover), filepath.Join(path, s.Cover))
			if err != nil {
				logger.With(zap.Error(err)).Sugar().Errorf("tools.Cp(%s)", filepath.Join(path, s.Cover))
			}
		}
		for _, image := range s.Images {
			if image != "" {
				err := tools.Cp(filepath.Join(config.DataDir, image), filepath.Join(path, image))
				if err != nil {
					logger.With(zap.Error(err)).Sugar().Errorf("tools.Cp(%s)", filepath.Join(path, image))
				}
			}
		}
	}
	return nil
}

func downloadInfo(obj any, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	return enc.Encode(obj)
}

func SaveHtmlImg(ctx context.Context, html string) (string, error) {
	imagesM := make(map[string]string)
	root, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(html)))
	if err != nil {
		return html, err
	}
	images := []string{}
	root.Find("img").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if ok {
			if strings.HasPrefix(src, "http") {
				imagesM[src] = ""
				images = append(images, src)
			} else {
				err := tools.Cp(filepath.Join(config.TmpDir, strings.TrimPrefix("/api/file/", src)), filepath.Join(config.DataDir, utils.GetFileName(src)))
				if err != nil && !errors.Is(err, os.ErrExist) {
					zaplog.L().Error("rename file error", zap.String("file path", src), zap.Error(err))
				}
			}
		}
	})

	lock := sync.Mutex{}
	wait := sync.WaitGroup{}
	for _, image := range images {
		url := image
		wait.Add(1)
		gopool.CtxGo(ctx, func() {
			defer wait.Done()

			rsp, err := tools.Req(http.MethodGet, url, nil)
			if err != nil {
				zaplog.L().Error("fetch file error", zap.String("file url", url), zap.Error(err))
				return
			}
			if rsp.StatusCode() != http.StatusOK {
				zaplog.L().Error("fetch file status code not 200", zap.Int("status code", rsp.StatusCode()))
				return
			}
			path, err := tools.SaveFileWithMd5Name(bytes.NewBuffer(rsp.Bytes()), config.DataDir, filepath.Ext(utils.GetFileName(url)))
			if err != nil {
				zaplog.L().Error("save file error", zap.Error(err))
				return
			}
			p := utils.GetFileName(path)
			lock.Lock()
			imagesM[url] = p
			lock.Unlock()
		})

	}
	wait.Wait()

	for k, v := range imagesM {
		html = strings.ReplaceAll(html, k, v)
	}

	return html, nil
}
