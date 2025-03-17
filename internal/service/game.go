package service

import (
	"bytes"
	"context"
	"errors"
	"hermes/config"
	"hermes/db"
	"hermes/internal/handler"
	"hermes/model"
	"hermes/tools"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	comm_errors "github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"go.uber.org/zap"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

type GameWhereNodeFunc func(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption)

type IGame interface {
	CreateL(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff, requestID string) error
	UpdateL(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff) error
	GetVOByID(ctx context.Context, id uint) (*handler.GameVo, error)

	Search(ctx context.Context, param handler.GameListReq, opt *meta.ListOption, gwfs ...GameWhereNodeFunc) (int64, []handler.GameVo, error)
	SaveFiles(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff) error
	// WhereNodeKeyword(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption)
	// WhereNodeTag(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption)
	// WhereNodeCharacter(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption)
	// WhereNodeStaff(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption)
	// WhereNodeSeries(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption)
	// WhereNodeCategory(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption)
	// WhereNodeDeveloper(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption)
	// WhereNodePublisher(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption)
	// WhereNodeSizeRange(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption)
	// WhereNodeIssueDateRange(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption)
	// WhereNodeCreatedAtRange(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption)
}

var _ IGame = (*game)(nil)

type game struct {
	store db.IStore
}

func (gsrv *game) CreateL(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff, requestID string) error {
	tx := gsrv.store.Transaction().Begin()
	err := tx.Game().Create(ctx, g, &meta.CreateOption{Omit: []string{"Series", "Developer", "Category", "Publisher"}})
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	ref, err := tx.RefGameInstance().Get(ctx, &model.RefGameInstance{RequestID: requestID}, nil)
	if err == nil {
		err = tx.GameInstance().Create(ctx, &model.GameInstance{GameID: g.ID, Version: ref.Version, Path: ref.Path}, nil)
		if err != nil {
			tx.Transaction().Rollback()
			return err
		}
	}

	var gSeries = []*model.GameSeries{}
	for _, s := range g.Series {
		gSeries = append(gSeries, &model.GameSeries{
			GameID:   g.ID,
			SeriesID: s.ID,
		})
	}
	err = tx.GameSeries().Creates(ctx, gSeries, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	// character
	charactersCreate := []*model.Character{}
	charactersUpdate := []*model.Character{}
	for _, c := range cs {
		if c.Character.ID != 0 {
			charactersUpdate = append(charactersUpdate, c.Character)
		} else {
			charactersCreate = append(charactersCreate, c.Character)
		}
	}
	if len(charactersCreate) > 0 {
		err = tx.Character().Creates(ctx, charactersCreate, nil)
		if err != nil {
			tx.Transaction().Rollback()
			return err
		}
	}
	errs := tx.Character().UpdateCollection(ctx, charactersUpdate, nil)
	if len(errs) > 0 {
		for _, err := range errs {
			if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
				tx.Transaction().Rollback()
				return err
			}
		}
	}
	err = tx.GameCharacter().Delete(ctx, &model.GameCharacter{GameID: g.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	for _, c := range cs {
		c.GameID = g.ID
		c.CharacterID = c.Character.ID
	}
	err = tx.GameCharacter().Creates(ctx, cs, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}

	// staff
	staffCreate := []*model.Person{}
	staffUpdate := []*model.Person{}
	for _, s := range ss {
		if s.Person.ID == 0 {
			staffCreate = append(staffCreate, s.Person)
		} else {
			staffUpdate = append(staffUpdate, s.Person)
		}
	}
	if len(staffCreate) > 0 {
		err = tx.Person().Creates(ctx, staffCreate, nil)
		if err != nil {
			tx.Transaction().Rollback()
			return err
		}
	}
	errs = tx.Person().UpdateCollection(ctx, staffUpdate, nil)
	if len(errs) > 0 {
		for _, err := range errs {
			if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
				tx.Transaction().Rollback()
				return err
			}
		}
	}
	err = tx.GameStaff().Delete(ctx, &model.GameStaff{GameID: g.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}

	var staff []*model.GameStaff
	for _, s := range ss {
		for _, r := range s.Relations {
			staff = append(staff, &model.GameStaff{
				GameID:   g.ID,
				PersonID: s.Person.ID,
				Relation: r,
			})
		}
	}
	if len(staff) > 0 {
		err = tx.GameStaff().Creates(ctx, staff, nil)
		if err != nil {
			tx.Transaction().Rollback()
			return err
		}
	}

	tx.Transaction().Commit()
	return nil
}

func (gsrv *game) UpdateL(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff) error {
	tx := gsrv.store.Transaction().Begin()
	err := tx.GameTag().Delete(ctx, &model.GameTag{GameID: g.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	err = tx.GameSeries().Delete(ctx, &model.GameSeries{GameID: g.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	var gSeries = []*model.GameSeries{}
	for _, s := range g.Series {
		gSeries = append(gSeries, &model.GameSeries{
			GameID:   g.ID,
			SeriesID: s.ID,
		})
	}
	err = tx.GameSeries().Creates(ctx, gSeries, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	err = tx.GameCharacter().Delete(ctx, &model.GameCharacter{GameID: g.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	charactersCreate := []*model.Character{}
	charactersUpdate := []*model.Character{}
	for _, c := range cs {
		if c.Character.ID != 0 {
			charactersUpdate = append(charactersUpdate, c.Character)
		} else {
			charactersCreate = append(charactersCreate, c.Character)
		}
	}
	if len(charactersCreate) > 0 {
		err = tx.Character().Creates(ctx, charactersCreate, nil)
		if err != nil {
			tx.Transaction().Rollback()
			return err
		}
	}
	errs := tx.Character().UpdateCollection(ctx, charactersUpdate, nil)
	if len(errs) > 0 {
		for _, err := range errs {
			if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
				tx.Transaction().Rollback()
				return err
			}
		}
	}
	for _, c := range cs {
		c.CharacterID = c.Character.ID
	}
	err = tx.GameCharacter().Creates(ctx, cs, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	err = tx.GameStaff().Delete(ctx, &model.GameStaff{GameID: g.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	staffCreate := []*model.Person{}
	staffUpdate := []*model.Person{}
	for _, s := range ss {
		if s.Person.ID == 0 {
			staffCreate = append(staffCreate, s.Person)
		} else {
			staffUpdate = append(staffUpdate, s.Person)
		}
	}
	if len(staffCreate) > 0 {
		err = tx.Person().Creates(ctx, staffCreate, nil)
		if err != nil {
			tx.Transaction().Rollback()
			return err
		}
	}
	errs = tx.Person().UpdateCollection(ctx, staffUpdate, nil)
	if len(errs) > 0 {
		for _, err := range errs {
			if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
				tx.Transaction().Rollback()
				return err
			}
		}
	}
	err = tx.GameStaff().Delete(ctx, &model.GameStaff{GameID: g.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}

	var staff []*model.GameStaff
	for _, s := range ss {
		if s.Relation > 0 {
			for _, r := range s.Relations {
				staff = append(staff, &model.GameStaff{
					GameID:   g.ID,
					PersonID: s.Person.ID,
					Relation: r,
				})
			}
		} else {
			staff = append(staff, &model.GameStaff{
				GameID:   g.ID,
				PersonID: s.Person.ID,
				Relation: 0,
			})
		}

	}
	err = tx.GameStaff().Creates(ctx, staff, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}

	err = tx.Game().Update(ctx, g, &meta.UpdateOption{Omit: []string{"Series"}})
	if err != nil && !errors.Is(err, comm_errors.ErrNoUpdateRows) {
		tx.Transaction().Rollback()
		return err
	}
	tx.Transaction().Commit()
	return nil
}

func (gsrv *game) GetVOByID(ctx context.Context, id uint) (*handler.GameVo, error) {
	g, err := gsrv.store.Game().Get(ctx, &model.Game{ID: uint(id)}, &meta.GetOption{Preload: []string{"Tags", "Category", "Series", "Developer", "Publisher"}})
	if err != nil {
		return nil, err
	}

	// character
	var cVos []handler.CharacterVo
	gcs, err := gsrv.store.GameCharacter().List(ctx, &model.GameCharacter{GameID: g.ID}, nil)
	if err != nil {
		return nil, err
	}
	crMap := map[uint]model.CharacterRelation{}
	cIDs := []uint{}
	for _, gc := range gcs {
		crMap[gc.CharacterID] = gc.Relation
		cIDs = append(cIDs, gc.CharacterID)
	}
	node := &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "id",
				Operator: meta.IN,
				Value:    cIDs,
			},
		},
	}
	cs, err := gsrv.store.Character().ListComplex(ctx, &model.Character{}, node, &meta.ListOption{Page: 1, PageSize: 100, GetOption: meta.GetOption{Preload: []string{"CV", "Tags"}}})
	if err != nil {
		return nil, err
	}
	for _, c := range cs {
		cVos = append(cVos, handler.CharacterVo{
			ID:      c.ID,
			Name:    c.Name,
			Alias:   c.Alias,
			Gender:  c.Gender.String(),
			Rlation: crMap[c.ID].String(),
			Summary: c.Summary,
			Cover:   c.Cover,
			Images:  c.Images,
			Tags:    c.Tags,
			CV: handler.StaffVo{
				ID:        c.CV.ID,
				Name:      c.CV.Name,
				Cover:     c.CV.Cover,
				Images:    c.CV.Images,
				Alias:     c.CV.Alias,
				CreatedAt: c.CV.CreatedAt,
				Tags:      c.CV.Tags,
				Gender:    c.CV.Gender.String(),
				Summary:   c.CV.Summary,
			},
			CreatedAt: c.CreatedAt,
		})
	}

	// staff
	var sVos []handler.StaffVo
	gss, err := gsrv.store.GameStaff().List(ctx, &model.GameStaff{GameID: g.ID}, nil)
	if err != nil {
		return nil, err
	}
	prMap := map[uint][]model.PersonRelation{}
	pIDs := []uint{}
	for _, gs := range gss {
		prs, ok := prMap[gs.PersonID]
		if ok {
			prs = append(prs, gs.Relation)
			prMap[gs.PersonID] = prs
		} else {
			pIDs = append(pIDs, gs.PersonID)
		}
	}
	node = &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "id",
				Operator: meta.IN,
				Value:    pIDs,
			},
		},
	}
	ss, err := gsrv.store.Person().ListComplex(ctx, &model.Person{}, node, nil)
	if err != nil {
		return nil, err
	}
	for _, s := range ss {
		var prs []string
		for _, pr := range prMap[s.ID] {
			prs = append(prs, pr.String())
		}
		sVos = append(sVos, handler.StaffVo{
			ID:        s.ID,
			Name:      s.Name,
			Alias:     s.Alias,
			Cover:     s.Cover,
			Images:    s.Images,
			Tags:      s.Tags,
			Summary:   s.Summary,
			Gender:    s.Gender.String(),
			Relation:  prs,
			CreatedAt: s.CreatedAt,
		})
	}
	return &handler.GameVo{
		ID:     g.ID,
		Name:   g.Name,
		Alias:  g.Alias,
		Cover:  g.Cover,
		Images: g.Images,
		// Versions:   version,
		Category:   g.Category,
		Series:     g.Series,
		Developer:  g.Developer,
		Publisher:  g.Publisher,
		Price:      g.Price,
		IssueDate:  g.IssueDate,
		Story:      g.Story,
		Platform:   g.Platform,
		Tags:       g.Tags,
		Characters: cVos,
		Staff:      sVos,
		Links:      g.Links,
		OtherInfo:  g.OtherInfo,
		CreatedAt:  g.CreatedAt,
	}, nil
}

func (gsrv *game) Search(ctx context.Context, param handler.GameListReq, opt *meta.ListOption, gwfs ...GameWhereNodeFunc) (int64, []handler.GameVo, error) {
	head := &meta.WhereNode{}
	node := head
	if opt == nil {
		opt = meta.NewListOption(nil, meta.WithPage(param.Page), meta.WithPageSize(param.PageSize))
	}
	opt.GetOption.Preload = append(opt.GetOption.Preload, []string{"Tags", "Category", "Series"}...)
	for _, f := range gwfs {
		node, opt = f(ctx, param, node, opt)
	}
	gs, err := gsrv.store.Game().ListComplex(ctx, &model.Game{}, head.Next, opt)
	if err != nil {
		return 0, nil, err
	}
	total, err := gsrv.store.Game().CountComplex(ctx, &model.Game{}, head.Next, &opt.GetOption)
	if err != nil {
		return 0, nil, err
	}
	gvos := make([]handler.GameVo, 0, len(gs))
	for _, g := range gs {
		gvos = append(gvos, handler.GameVo{
			ID:        g.ID,
			Name:      g.Name,
			Cover:     g.Cover,
			Alias:     g.Alias,
			Images:    g.Images,
			Category:  g.Category,
			Series:    g.Series,
			Price:     g.Price,
			IssueDate: g.IssueDate,
			Story:     g.Story,
			Platform:  g.Platform,
			Tags:      g.Tags,
			Links:     g.Links,
			OtherInfo: g.OtherInfo,
		})
	}

	return total, gvos, nil
}

func (gsrv *game) SaveFiles(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff) error {
	images := []string{}
	urls := map[string][]*string{}
	f := func(image *string) bool {
		if len(*image) > 4 && (*image)[:4] == "http" {
			_, ok := urls[*image]
			if !ok {
				urls[*image] = []*string{}
			}
			urls[*image] = append(urls[*image], image)

			return false
		} else if len(*image) > 0 {
			images = append(images, *image)
			return true
		}
		return true
	}

	if f(&g.Cover) {
		g.Cover = tools.GetFileName(g.Cover)
	}
	for i := range g.Images {
		if f(&g.Images[i]) {
			g.Images[i] = tools.GetFileName(g.Images[i])
		}
	}
	for _, c := range cs {
		if f(&c.Character.Cover) {
			c.Character.Cover = tools.GetFileName(c.Character.Cover)
		}
		for i := range c.Character.Images {
			if f(&c.Character.Images[i]) {
				c.Character.Images[i] = tools.GetFileName(c.Character.Images[i])
			}
		}
	}
	for _, s := range ss {
		if f(&s.Person.Cover) {
			s.Person.Cover = tools.GetFileName(s.Person.Cover)
		}
	}
	wait := sync.WaitGroup{}
	for k, vs := range urls {
		url := k
		wait.Add(1)
		gopool.CtxGo(ctx, func() {
			defer wait.Done()

			data, code, err := tools.MakeRequest(http.MethodGet, url, config.GetConfig().ProxyConfig, nil, nil, nil, config.DefaultRetryCnt)
			if err != nil {
				zaplog.L().Error("fetch file error", zap.String("file url", url), zap.Error(err))
				return
			}
			if code != http.StatusOK {
				zaplog.L().Error("fetch file status code not 200", zap.Int("status code", code))
				return
			}
			path, err := tools.SaveFile(filepath.Ext(url), bytes.NewBuffer(data), config.Dir)
			if code != http.StatusOK {
				zaplog.L().Error("save file error", zap.Error(err))
				return
			}
			p := tools.GetFileName(path)
			for _, v := range vs {
				*v = p
			}
		})
	}
	zaplog.L().Info("images", zap.Any("data", images))
	zaplog.L().Info("dir", zap.Any("name", config.Dir))
	for _, image := range images {
		err := os.Rename(image, filepath.Join(config.Dir, tools.GetFileName(image)))
		if err != nil {
			zaplog.L().Error("rename file error", zap.String("file path", image), zap.Error(err))
		}
	}

	wait.Wait()
	return nil
}

func NewGame(store db.IStore) *game {
	return &game{store: store}
}
