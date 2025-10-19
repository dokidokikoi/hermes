package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hermes/config"
	"hermes/db"
	"hermes/db/data"
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
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

type GameWhereNodeFunc func(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption)

type IGame interface {
	CreateL(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff, gIns *model.GameInstance) error
	UpdateL(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff) error
	GetVOByID(ctx context.Context, id uint) (*handler.GameVo, error)

	Search(ctx context.Context, param handler.GameListReq, opt *meta.ListOption, gwfs ...GameWhereNodeFunc) (int64, []handler.GameVo, error)
	SaveFiles(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff) error

	Load(ctx context.Context, gVo *handler.GameVo, path string) error
}

var _ IGame = (*game)(nil)

type game struct {
	store db.IStore
}

func (gsrv *game) CreateL(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff, gIns *model.GameInstance) error {
	tx := gsrv.store.Transaction().Begin()
	err := tx.Game().Create(ctx, g, &meta.CreateOption{Omit: []string{"Series", "Developer", "Category"}})
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	if gIns != nil {
		err = tx.GameInstance().Create(ctx, gIns, nil)
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
			c.Character.UUID = uuid.NewString()
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
	if len(cs) > 0 {
		err = tx.GameCharacter().Creates(ctx, cs, nil)
		if err != nil {
			tx.Transaction().Rollback()
			return err
		}
	}

	// staff
	staffCreate := []*model.Person{}
	staffUpdate := []*model.Person{}
	for _, s := range ss {
		if s.Person.ID == 0 {
			s.Person.UUID = uuid.NewString()
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

	for _, s := range ss {
		s.GameID = g.ID
		s.PersonID = s.Person.ID
	}
	if len(ss) > 0 {
		err = tx.GameStaff().Creates(ctx, ss, nil)
		if err != nil {
			tx.Transaction().Rollback()
			return err
		}
	}

	return tx.Transaction().Commit()
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
			c.Character.UUID = uuid.NewString()
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
			s.Person.UUID = uuid.NewString()
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

	err = tx.GameStaff().Creates(ctx, ss, nil)
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
	g, err := gsrv.store.Game().Get(ctx, &model.Game{ID: uint(id)}, &meta.GetOption{Preload: []string{"Tags", "Category", "Series", "Developer"}})
	if err != nil {
		return nil, err
	}

	// character
	var cVos []handler.CharacterVo
	gcs, err := gsrv.store.GameCharacter().List(
		ctx,
		&model.GameCharacter{GameID: g.ID},
		nil,
	)
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
	cs, err := gsrv.store.Character().ListComplex(
		ctx,
		&model.Character{},
		node,
		&meta.ListOption{
			Page:      1,
			PageSize:  100,
			GetOption: meta.GetOption{Preload: []string{"CV", "Tags"}},
			Order:     "weight desc",
		})
	if err != nil {
		return nil, err
	}
	for _, c := range cs {
		cVos = append(cVos, handler.CharacterVo{
			ID:      c.ID,
			Name:    c.Name,
			Alias:   c.Alias,
			Gender:  c.Gender,
			Rlation: crMap[c.ID],
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
				Gender:    c.CV.Gender,
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
	prMap := map[uint]model.PersonRelations{}
	pIDs := []uint{}
	for _, gs := range gss {
		prMap[gs.PersonID] = gs.Relations
		pIDs = append(pIDs, gs.PersonID)
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
		sVos = append(sVos, handler.StaffVo{
			ID:        s.ID,
			Name:      s.Name,
			Alias:     s.Alias,
			Cover:     s.Cover,
			Images:    s.Images,
			Tags:      s.Tags,
			Summary:   s.Summary,
			Gender:    s.Gender,
			Relation:  prMap[s.ID],
			CreatedAt: s.CreatedAt,
		})
	}
	return &handler.GameVo{
		ID:         g.ID,
		Name:       g.Name,
		Alias:      g.Alias,
		Cover:      g.Cover,
		Images:     g.Images,
		Category:   g.Category,
		Series:     g.Series,
		Developer:  g.Developer,
		Price:      g.Price,
		IssueDate:  g.IssueDate,
		Story:      g.Story,
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

			rsp, err := tools.Req(http.MethodGet, url, nil)
			if err != nil {
				zaplog.L().Error("fetch file error", zap.String("file url", url), zap.Error(err))
				return
			}
			if rsp.StatusCode() != http.StatusOK {
				zaplog.L().Error("fetch file status code not 200", zap.Int("status code", rsp.StatusCode()))
				return
			}
			path, err := tools.SaveFile(filepath.Ext(url), bytes.NewBuffer(rsp.Bytes()), config.Dir)
			if err != nil {
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

func (gsrv *game) Load(ctx context.Context, gVo *handler.GameVo, path string) (e error) {
	db := data.GetDataFactory()
	cate, err := db.Category().Get(ctx, &model.Category{Name: gVo.Category.Name}, &meta.GetOption{Select: []string{"ID"}})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("create category error: %w", err)
		}
	}
	if cate == nil {
		err = db.Category().Create(ctx, gVo.Category, nil)
		if err != nil {
			return fmt.Errorf("create category error: %w", err)
		}
	} else {
		gVo.Category.ID = cate.ID
	}

	dev, err := db.Developer().Get(ctx, &model.Developer{Name: gVo.Developer.Name}, &meta.GetOption{Select: []string{"ID"}})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("create developer error: %w", err)
		}
	}
	if dev == nil {
		err = db.Developer().Create(ctx, gVo.Developer, nil)
		if err != nil {
			return fmt.Errorf("create developer error: %w", err)
		}
	} else {
		gVo.Developer.ID = dev.ID
	}

	for _, series := range gVo.Series {
		s, err := db.Series().Get(ctx, &model.Series{Name: series.Name}, &meta.GetOption{Select: []string{"ID"}})
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("create series error: %w", err)
			}
		}
		if s == nil {
			err = db.Series().Create(ctx, series, nil)
			if err != nil {
				return fmt.Errorf("create category error: %w", err)
			}
		} else {
			series.ID = s.ID
		}
	}
	for _, tag := range gVo.Tags {
		t, err := db.Tag().Get(ctx, &model.Tag{Name: tag.Name}, &meta.GetOption{Select: []string{"ID"}})
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("create tag error: %w", err)
			}
		}
		if t == nil {
			err = db.Tag().Create(ctx, tag, nil)
			if err != nil {
				return fmt.Errorf("create tag error: %w", err)
			}
		} else {
			tag.ID = t.ID
		}
	}

	var cs []*model.GameCharacter
	charactersCreate := []*model.Character{}
	charactersUpdate := []*model.Character{}
	for _, character := range gVo.Characters {
		c := &model.Character{
			UUID:     character.UUID,
			Name:     character.Name,
			Alias:    character.Alias,
			Cover:    character.Cover,
			Images:   character.Images,
			Tags:     character.Tags,
			Summary:  character.Summary,
			Gender:   character.Gender,
			PersonID: character.CV.ID,
		}
		cs = append(cs, &model.GameCharacter{
			Character: c,
			Relation:  character.Rlation,
		})

		cRsp, err := db.Character().Get(ctx, &model.Character{ID: character.ID}, &meta.GetOption{Select: []string{"ID", "UUID"}})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				charactersCreate = append(charactersCreate, c)
				continue
			}
			return fmt.Errorf("get character error: %w", err)
		}
		if cRsp.UUID != character.UUID {
			charactersCreate = append(charactersCreate, c)
		} else {
			charactersUpdate = append(charactersUpdate, c)
		}
	}

	// staff
	var staffIDM = make(map[uint]*model.Person)
	var ss []*model.GameStaff
	staffCreate := []*model.Person{}
	staffUpdate := []*model.Person{}
	for _, staff := range gVo.Staff {
		s := &model.Person{
			UUID:    staff.UUID,
			Name:    staff.Name,
			Alias:   staff.Alias,
			Cover:   staff.Cover,
			Images:  staff.Images,
			Tags:    staff.Tags,
			Summary: staff.Summary,
			Gender:  staff.Gender,
		}
		ss = append(ss, &model.GameStaff{
			Person:    s,
			Relations: staff.Relation,
		})
		staffIDM[staff.ID] = s

		sRsp, err := db.Person().Get(ctx, &model.Person{ID: staff.ID}, &meta.GetOption{Select: []string{"ID", "UUID"}})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				staffCreate = append(staffCreate, s)
				continue
			}
			return fmt.Errorf("get character error: %w", err)
		}
		if sRsp.UUID != staff.UUID {
			staffCreate = append(staffCreate, s)
		} else {
			staffUpdate = append(staffUpdate, s)
		}
	}

	// 启动事务
	tx := gsrv.store.Transaction().Begin()
	defer func() {
		if e != nil {
			tx.Transaction().Rollback()
		}
	}()

	if len(staffCreate) > 0 {
		err = tx.Person().Creates(ctx, staffCreate, nil)
		if err != nil {
			tx.Transaction().Rollback()
			return err
		}
	}
	errs := tx.Person().UpdateCollection(ctx, staffUpdate, nil)
	if len(errs) > 0 {
		for _, err := range errs {
			if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
				tx.Transaction().Rollback()
				return err
			}
		}
	}

	// 设置 PersonID
	for _, character := range cs {
		staff := staffIDM[character.Character.PersonID]
		if staff != nil {
			character.Character.PersonID = staff.ID
		}
	}
	if len(charactersCreate) > 0 {
		err = tx.Character().Creates(ctx, charactersCreate, nil)
		if err != nil {
			return err
		}
	}
	errs = tx.Character().UpdateCollection(ctx, charactersUpdate, nil)
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	g := &model.Game{
		Name:        gVo.Name,
		UUID:        gVo.UUID,
		Alias:       gVo.Alias,
		JanCode:     gVo.JanCode,
		Code:        gVo.Code,
		Cover:       gVo.Cover,
		Images:      gVo.Images,
		Tags:        gVo.Tags,
		CategoryID:  gVo.Category.ID,
		DeveloperID: gVo.Developer.ID,
		OtherInfo:   gVo.OtherInfo,
		Story:       gVo.Story,
		Price:       gVo.Price,
		IssueDate:   gVo.IssueDate,
	}
	gRsp, err := tx.Game().Get(ctx, &model.Game{ID: gVo.ID}, &meta.GetOption{Select: []string{"ID", "UUID"}})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	} else {
		if gRsp.UUID != gVo.UUID {
			err = tx.Game().Create(ctx, g, &meta.CreateOption{Omit: []string{clause.Associations}})
			if err != nil {
				return err
			}
			if gVo.Version != "" {
				err = tx.GameInstance().Create(ctx, &model.GameInstance{
					GameID:   g.ID,
					Version:  gVo.Version,
					Path:     path,
					Size:     gVo.Size,
					Language: gVo.Language,
					Comment:  gVo.Comment,
				}, nil)
				if err != nil {
					return err
				}
			}
		} else {
			err = tx.Game().Update(ctx, g, &meta.UpdateOption{Omit: []string{clause.Associations}})
			if err != nil {
				return err
			}
			err = tx.GameInstance().Update(ctx, &model.GameInstance{
				GameID:   g.ID,
				Version:  gVo.Version,
				Path:     path,
				Size:     gVo.Size,
				Language: gVo.Language,
				Comment:  gVo.Comment,
			}, nil)
			if err != nil {
				return err
			}
			err = tx.GameCharacter().Delete(ctx, &model.GameCharacter{GameID: g.ID}, nil)
			if err != nil {
				return err
			}
			err = tx.GameSeries().Delete(ctx, &model.GameSeries{GameID: g.ID}, nil)
			if err != nil {
				return err
			}
			err = tx.GameTag().Delete(ctx, &model.GameTag{GameID: g.ID}, nil)
			if err != nil {
				return err
			}
			err = tx.GameStaff().Delete(ctx, &model.GameStaff{GameID: g.ID}, nil)
			if err != nil {
				return err
			}
		}
	}

	var gameSeries []*model.GameSeries
	for _, series := range gVo.Series {
		gameSeries = append(gameSeries, &model.GameSeries{
			GameID:   g.ID,
			SeriesID: series.ID,
		})
	}
	err = tx.GameSeries().Creates(ctx, gameSeries, nil)
	if err != nil {
		return err
	}

	var gameTags []*model.GameTag
	for _, tag := range gVo.Tags {
		gameTags = append(gameTags, &model.GameTag{
			GameID: g.ID,
			TagID:  tag.ID,
		})
	}
	err = tx.GameTag().Creates(ctx, gameTags, nil)
	if err != nil {
		return err
	}

	for _, staff := range ss {
		staff.GameID = g.ID
	}
	err = tx.GameStaff().Creates(ctx, ss, nil)
	if err != nil {
		return err
	}

	for _, character := range cs {
		character.GameID = g.ID
	}
	err = tx.GameCharacter().Creates(ctx, cs, nil)
	if err != nil {
		return err
	}

	return tx.Transaction().Commit()
}

func NewGame(store db.IStore) *game {
	return &game{store: store}
}
