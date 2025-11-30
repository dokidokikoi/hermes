package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"izumi/config"
	"izumi/db"
	"izumi/db/data"
	"izumi/internal/handler"
	"izumi/model"
	"izumi/utils"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/abadojack/whatlanggo"
	comm_errors "github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/dokidokikoi/go-common/tools"
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
	UpsertFull(ctx context.Context, gVo *handler.GameVo, gIns *model.GameInstance) error

	Search(ctx context.Context, param handler.GameListReq, opt *meta.ListOption, gwfs ...GameWhereNodeFunc) (int64, []handler.GameVo, error)
	SaveFiles(ctx context.Context, g *handler.GameVo) error

	Load(ctx context.Context, gVo *handler.GameVo, path string) error
	DownloadInfo(ctx context.Context, gameID uint, t time.Time) error
}

var _ IGame = (*game)(nil)

type game struct {
	store db.IStore
}

func (gsrv *game) UpsertFull(ctx context.Context, gVo *handler.GameVo, gIns *model.GameInstance) error {
	err := gsrv.SaveFiles(ctx, gVo)
	if err != nil {
		return err
	}
	db := gsrv.store
	if gVo.VNDBID != "" {
		g, err := db.Game().Get(ctx, &model.Game{VNDBID: gVo.VNDBID}, &meta.GetOption{Include: []string{"ID"}})
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if g != nil {
			gVo.ID = g.ID
		}
	}
	g := tools.GetPtr(handler.Vo2Game(*gVo))
	if gVo.Category != nil && gVo.Category.Name != "" {
		c, err := db.Category().Get(ctx, &model.Category{Name: gVo.Category.Name}, &meta.GetOption{Include: []string{"ID"}})
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if c == nil {
			err = db.Category().Create(ctx, gVo.Category, nil)
			if err != nil {
				return err
			}
			g.CategoryID = gVo.Category.ID
		} else {
			g.CategoryID = c.ID
		}
	}
	{
		g.Story, err = SaveHtmlImg(ctx, g.Story)
		if err != nil {
			zaplog.From(ctx).With(zap.Error(err)).Error("SaveHtmlImg")
		}
		if g.ID != 0 {
			err := db.Game().Update(ctx, g, &meta.UpdateOption{Omit: []string{clause.Associations}})
			if err != nil {
				return err
			}
		} else {
			g.UUID = uuid.New().String()
			err := db.Game().Create(ctx, g, &meta.CreateOption{Omit: []string{clause.Associations}})
			if err != nil {
				return err
			}
		}
	}
	if gIns != nil {
		gIns.GameID = g.ID
		i, err := db.GameInstance().Get(ctx, &model.GameInstance{GameID: g.ID, Version: gIns.Version}, &meta.GetOption{Include: []string{"ID"}})
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if i == nil {
			err := db.GameInstance().Create(ctx, gIns, nil)
			if err != nil {
				return err
			}
		} else {
			gIns.ID = i.ID
			err := db.GameInstance().Update(ctx, gIns, nil)
			if err != nil {
				return err
			}
		}
	}
	if len(gVo.Brands) > 0 {
		gameBrands := []*model.GameBrands{}
		for _, brand := range gVo.Brands {
			b, err := db.Brand().Get(ctx, &model.Brand{Name: brand.Name}, &meta.GetOption{Include: []string{"ID"}})
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if b == nil {
				err = db.Brand().Create(ctx, brand, nil)
				if err != nil {
					return err
				}
				gameBrands = append(gameBrands, &model.GameBrands{
					GameID:  g.ID,
					BrandID: brand.ID,
				})
			} else {
				gameBrands = append(gameBrands, &model.GameBrands{
					GameID:  g.ID,
					BrandID: b.ID,
				})
			}
			tx := db.Transaction().Begin()
			err = tx.GameBrands().Delete(ctx, &model.GameBrands{GameID: g.ID}, nil)
			if err != nil {
				tx.Transaction().Rollback()
				return err
			}
			err = tx.GameBrands().Creates(ctx, gameBrands, nil)
			if err != nil {
				tx.Transaction().Rollback()
				return err
			}
			err = tx.Transaction().Commit()
			if err != nil {
				return err
			}
		}
	}
	if len(gVo.Series) > 0 {
		gameSeries := []*model.GameSeries{}
		for _, series := range gVo.Series {
			s, err := db.Series().Get(ctx, &model.Series{Name: series.Name}, &meta.GetOption{Include: []string{"ID"}})
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if s == nil {
				err = db.Series().Create(ctx, series, nil)
				if err != nil {
					return err
				}
				gameSeries = append(gameSeries, &model.GameSeries{
					GameID:   g.ID,
					SeriesID: series.ID,
				})
			} else {
				gameSeries = append(gameSeries, &model.GameSeries{
					GameID:   g.ID,
					SeriesID: s.ID,
				})
			}
			tx := db.Transaction().Begin()
			err = tx.GameSeries().Delete(ctx, &model.GameSeries{GameID: g.ID}, nil)
			if err != nil {
				tx.Transaction().Rollback()
				return err
			}
			err = tx.GameSeries().Creates(ctx, gameSeries, nil)
			if err != nil {
				tx.Transaction().Rollback()
				return err
			}
			err = tx.Transaction().Commit()
			if err != nil {
				tx.Transaction().Rollback()
				return err
			}
		}
	}
	if len(gVo.Tags) > 0 {
		gameTags := []*model.GameTag{}
		for _, tag := range gVo.Tags {
			t, err := db.Tag().Get(ctx, &model.Tag{Name: tag.Name}, &meta.GetOption{Include: []string{"ID"}})
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if t == nil {
				lang := whatlanggo.DetectLang(tag.Name)
				if lang.Iso6391() == "zh" || lang.Iso6391() == "ja" {
					tag.Lang = lang.Iso6391()
				} else {
					tag.Lang = "en"
				}
				err = db.Tag().Create(ctx, tag, nil)
				if err != nil {
					return err
				}
				gameTags = append(gameTags, &model.GameTag{
					GameID: g.ID,
					TagID:  tag.ID,
				})
			} else {
				gameTags = append(gameTags, &model.GameTag{
					GameID: g.ID,
					TagID:  t.ID,
				})
			}
		}
		tx := db.Transaction().Begin()
		err := tx.GameTag().Delete(ctx, &model.GameTag{GameID: g.ID}, nil)
		if err != nil {
			tx.Transaction().Rollback()
			return err
		}
		err = tx.GameTag().Creates(ctx, gameTags, nil)
		if err != nil {
			tx.Transaction().Rollback()
			return err
		}
		err = tx.Transaction().Commit()
		if err != nil {
			tx.Transaction().Rollback()
			return err
		}
	}
	staffM := map[string]uint{}
	if len(gVo.Staff) > 0 {
		gameStaff := []*model.GameStaff{}
		for _, staff := range gVo.Staff {
			if staff.VNDBID != "" {
				s, err := db.Person().Get(ctx, &model.Person{VNDBID: staff.VNDBID}, &meta.GetOption{Include: []string{"ID"}})
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				if s != nil {
					staff.ID = s.ID
				}
			}
			if staff.ID == 0 {
				list, err := data.GetDataFactory().Person().ListComplex(
					ctx,
					&model.Person{},
					&meta.WhereNode{
						Conditions: []*meta.Condition{
							{
								Field:    "vndb_id",
								Operator: meta.EQUAL,
								Value:    "",
							},
							{
								Field:    "name",
								Operator: meta.EQUAL,
								Value:    staff.Name,
							},
						},
					},
					&meta.ListOption{GetOption: meta.GetOption{Include: []string{"ID"}}})
				if err != nil {
					return err
				}
				if len(list) > 0 {
					staff.ID = list[0].ID
					staff.UUID = list[0].UUID
				}
			}
			s := tools.GetPtr(handler.Vo2Person(staff))
			if staff.ID != 0 {
				err := db.Person().Update(ctx, s, nil)
				if err != nil {
					return err
				}
			} else {
				s.UUID = uuid.New().String()
				err := db.Person().Create(ctx, s, nil)
				if err != nil {
					return err
				}
			}
			pRelationM := tools.NewSet(staff.Relation...)
			if pRelationM.Contains(model.PRelationUnknown) && pRelationM.Size() > 1 {
				pRelationM.Remove(model.PRelationUnknown)
			}
			staff.Relation = pRelationM.Slice()

			gameStaff = append(gameStaff, &model.GameStaff{
				PersonID:  s.ID,
				GameID:    g.ID,
				Relations: staff.Relation,
			})
			tx := db.Transaction().Begin()
			err := tx.GameStaff().Delete(ctx, &model.GameStaff{GameID: g.ID}, nil)
			if err != nil {
				tx.Transaction().Rollback()
				return err
			}
			err = tx.GameStaff().Creates(ctx, gameStaff, nil)
			if err != nil {
				tx.Transaction().Rollback()
				return err
			}
			err = tx.Transaction().Commit()
			if err != nil {
				tx.Transaction().Rollback()
				return err
			}
			staffM[s.Name] = s.ID
		}
	}
	if len(gVo.Characters) > 0 {
		gameCharacters := []*model.GameCharacter{}
		for _, char := range gVo.Characters {
			if char.VNDBID != "" {
				c, err := db.Character().Get(ctx, &model.Character{VNDBID: char.VNDBID}, &meta.GetOption{Include: []string{"ID"}})
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				if c != nil {
					char.ID = c.ID
				}
			}
			c := tools.GetPtr(handler.Vo2Character(char))
			if char.CV.Name != "" {
				c.PersonID = staffM[char.CV.Name]
			}
			if char.ID != 0 {
				err := db.Character().Update(ctx, c, &meta.UpdateOption{Omit: []string{clause.Associations}})
				if err != nil {
					return err
				}
			} else {
				c.UUID = uuid.New().String()
				err := db.Character().Create(ctx, c, &meta.CreateOption{Omit: []string{clause.Associations}})
				if err != nil {
					return err
				}
			}
			gameCharacters = append(gameCharacters, &model.GameCharacter{
				CharacterID: c.ID,
				GameID:      g.ID,
				Relation:    char.Rlation,
			})
			tx := db.Transaction().Begin()
			err := tx.GameCharacter().Delete(ctx, &model.GameCharacter{GameID: g.ID}, nil)
			if err != nil {
				tx.Transaction().Rollback()
				return err
			}
			err = tx.GameCharacter().Creates(ctx, gameCharacters, nil)
			if err != nil {
				tx.Transaction().Rollback()
				return err
			}
			err = tx.Transaction().Commit()
			if err != nil {
				tx.Transaction().Rollback()
				return err
			}
		}
	}

	return nil
}

func (gsrv *game) CreateL(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff, gIns *model.GameInstance) error {
	var err error
	g.Story, err = SaveHtmlImg(ctx, g.Story)
	if err != nil {
		zaplog.From(ctx).With(zap.Error(err)).Error("SaveHtmlImg")
	}
	tx := gsrv.store.Transaction().Begin()
	err = tx.Game().Create(ctx, g, &meta.CreateOption{Omit: []string{"Series", "Brands", "Category"}})
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	if gIns != nil {
		gIns.GameID = g.ID
		err = tx.GameInstance().Create(ctx, gIns, nil)
		if err != nil {
			tx.Transaction().Rollback()
			return err
		}
	}

	var gBrands = []*model.GameBrands{}
	for _, b := range g.Brands {
		gBrands = append(gBrands, &model.GameBrands{
			GameID:  g.ID,
			BrandID: b.ID,
		})
	}
	err = tx.GameBrands().Creates(ctx, gBrands, nil)
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
	var err error
	g.Story, err = SaveHtmlImg(ctx, g.Story)
	if err != nil {
		zaplog.From(ctx).With(zap.Error(err)).Error("SaveHtmlImg")
	}
	tx := gsrv.store.Transaction().Begin()
	err = tx.GameTag().Delete(ctx, &model.GameTag{GameID: g.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	err = tx.GameSeries().Delete(ctx, &model.GameSeries{GameID: g.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	err = tx.GameBrands().Delete(ctx, &model.GameBrands{GameID: g.ID}, nil)
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
	var gBrands = []*model.GameBrands{}
	for _, b := range g.Brands {
		gBrands = append(gBrands, &model.GameBrands{
			GameID:  g.ID,
			BrandID: b.ID,
		})
	}
	err = tx.GameBrands().Creates(ctx, gBrands, nil)
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
	g, err := gsrv.store.Game().Get(ctx, &model.Game{ID: uint(id)}, &meta.GetOption{Preload: []string{"Tags", "Category", "Series", "Brands"}})
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
		cVos = append(cVos, handler.Character2Vo(*c, crMap[c.ID], nil))
	}

	sort.Slice(cVos, func(i, j int) bool {
		if cVos[i].Rlation == model.CRelationMain {
			return true
		}
		if cVos[j].Rlation == model.CRelationMain {
			return false
		}
		if cVos[i].Rlation == model.CRelationMinor {
			return true
		}
		if cVos[j].Rlation == model.CRelationMinor {
			return false
		}
		return true
	})

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
		sVos = append(sVos, handler.Person2Vo(*s, prMap[s.ID]))
	}
	sort.Slice(g.Tags, func(i, j int) bool {
		if g.Tags[j].Lang == "zh" {
			return false
		} else if g.Tags[j].Lang == "ja" && g.Tags[i].Lang != "zh" {
			return false
		} else {
			return true
		}
	})
	return tools.GetPtr(handler.Game2Vo(*g, cVos, sVos)), nil
}

func (gsrv *game) Search(ctx context.Context, param handler.GameListReq, opt *meta.ListOption, gwfs ...GameWhereNodeFunc) (int64, []handler.GameVo, error) {
	head := &meta.WhereNode{}
	node := head
	if opt == nil {
		opt = meta.NewListOption(nil, meta.WithPage(param.Page), meta.WithPageSize(param.PageSize))
	}
	opt.GetOption.Preload = append(opt.GetOption.Preload, []string{"Tags", "Category", "Series", "Brands"}...)
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
		sort.Slice(g.Tags, func(i, j int) bool {
			if g.Tags[j].Lang == "zh" {
				return false
			} else if g.Tags[j].Lang == "ja" && g.Tags[i].Lang != "zh" {
				return false
			} else {
				return true
			}
		})
		if len(g.Tags) > 10 {
			g.Tags = g.Tags[:10]
		}
		gvos = append(gvos, handler.Game2Vo(*g, nil, nil))
	}

	return total, gvos, nil
}

func (gsrv *game) SaveFiles(ctx context.Context, g *handler.GameVo) error {
	cs := g.Characters
	ss := g.Staff
	images := []string{}
	urls := map[string]map[string][]*string{}
	f := func(image *string) bool {
		if len(*image) > 4 && (*image)[:4] == "http" {
			u, err := url.Parse(*image)
			if err != nil {
				return false
			}
			_, ok := urls[u.Host]
			if !ok {
				urls[u.Host] = make(map[string][]*string, 0)
			}
			urls[u.Host][*image] = append(urls[u.Host][*image], image)

			return false
		} else if len(*image) > 0 {
			images = append(images, *image)
			return true
		}
		return true
	}

	if f(&g.Cover) {
		g.Cover = utils.GetFileName(g.Cover)
	}
	for i := range g.Images {
		if f(&g.Images[i]) {
			g.Images[i] = utils.GetFileName(g.Images[i])
		}
	}
	for i := range cs {
		if f(&cs[i].Cover) {
			cs[i].Cover = utils.GetFileName(cs[i].Cover)
		}
		for j := range cs[i].Images {
			if f(&cs[i].Images[j]) {
				cs[i].Images[j] = utils.GetFileName(cs[i].Images[j])
			}
		}
	}
	for i := range ss {
		if f(&ss[i].Cover) {
			ss[i].Cover = utils.GetFileName(ss[i].Cover)
		}
		for j := range ss[i].Images {
			if f(&ss[i].Images[j]) {
				ss[i].Images[j] = utils.GetFileName(ss[i].Images[j])
			}
		}
	}
	wait := sync.WaitGroup{}
	for _, us := range urls {
		wait.Add(1)
		gopool.CtxGo(ctx, func() {
			defer wait.Done()

			for url, vs := range us {
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
				for _, v := range vs {
					*v = p
				}
				time.Sleep(time.Millisecond * 50)
			}
		})
	}
	zaplog.L().Info("images", zap.Any("data", images))
	zaplog.L().Info("dir", zap.Any("name", config.DataDir))
	for _, image := range images {
		err := tools.Cp(image, filepath.Join(config.DataDir, utils.GetFileName(image)))
		if err != nil && !errors.Is(err, os.ErrExist) {
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
	for _, brand := range gVo.Brands {
		dev, err := db.Brand().Get(ctx, &model.Brand{Name: brand.Name}, &meta.GetOption{Select: []string{"ID"}})
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("create brand error: %w", err)
			}
		}
		if dev == nil {
			err = db.Brand().Create(ctx, brand, nil)
			if err != nil {
				return fmt.Errorf("create brand error: %w", err)
			}
		} else {
			brand.ID = dev.ID
		}
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
		c := tools.GetPtr(handler.Vo2Character(character))
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
		s := tools.GetPtr(handler.Vo2Person(staff))
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
		} else {
			loadGameAllImages(zaplog.From(ctx), path)
		}
	}()

	if len(staffCreate) > 0 {
		err = tx.Person().Creates(ctx, staffCreate, nil)
		if err != nil {
			return err
		}
	}
	errs := tx.Person().UpdateCollection(ctx, staffUpdate, nil)
	if len(errs) > 0 {
		for _, err := range errs {
			if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
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

	g := tools.GetPtr(handler.Vo2Game(*gVo))
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
			err = tx.GameBrands().Delete(ctx, &model.GameBrands{GameID: g.ID}, nil)
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

	var gameBrands []*model.GameBrands
	for _, brand := range gVo.Brands {
		gameBrands = append(gameBrands, &model.GameBrands{
			GameID:  g.ID,
			BrandID: brand.ID,
		})
	}
	err = tx.GameBrands().Creates(ctx, gameBrands, nil)
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

func (gsrv *game) DownloadInfo(ctx context.Context, gameID uint, t time.Time) error {
	gVo, err := gsrv.GetVOByID(ctx, gameID)
	if err != nil {
		return err
	}
	ins, err := data.GetDataFactory().
		GameInstance().List(ctx, &model.GameInstance{GameID: gameID}, nil)
	if err != nil {
		return err
	}
	if len(ins) == 0 {
		return errors.New("no game instance")
	}
	for _, i := range ins {
		if !filepath.IsAbs(i.Path) {
			return errors.New("game path need absoulte path")
		}
		info, err := os.Stat(i.Path)
		if err != nil {
			zaplog.From(ctx).Error("os.Stat", zap.Error(err))
			continue
		}
		if !info.IsDir() {
			continue
		}
		infoPath := filepath.Join(i.Path, "info.json")
		info, err = os.Stat(infoPath)
		if err != nil && !os.IsNotExist(err) {
			zaplog.From(ctx).Error("os.Stat", zap.Error(err))
			continue
		}
		if info != nil &&
			(t.Before(info.ModTime()) || (!gVo.UpdatedAt.IsZero() && gVo.UpdatedAt.Before(info.ModTime()))) {
			continue
		}
		err = cpGameAllImages(zaplog.From(ctx), i.Path, gVo)
		if err != nil {
			zaplog.From(ctx).Error("cpGameAllImages", zap.Error(err))
			continue
		}
		f, err := os.OpenFile(infoPath, os.O_WRONLY, 0666)
		if err != nil {
			if os.IsNotExist(err) {
				err = nil
				f, err = os.Create(infoPath)
			}
			if err != nil {
				zaplog.From(ctx).Error("Open", zap.Error(err))
				continue
			}
		}
		defer f.Close()

		gVo.Version = i.Version
		gVo.Language = i.Language
		gVo.Comment = i.Comment
		gVo.Size = i.Size
		err = downloadInfo(gVo, f)
		if err != nil {
			zaplog.From(ctx).Error("downloadInfo", zap.Error(err))
		}
	}

	return nil
}

func NewGame(store db.IStore) *game {
	return &game{store: store}
}

func loadGameAllImages(logger *zap.Logger, path string) error {
	path = filepath.Join(path, "images")
	files, err := os.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		_, fileName := filepath.Split(file.Name())
		dstPath := filepath.Join(config.DataDir, fileName)
		_, err := os.Stat(dstPath)
		if err != nil {
			if os.IsExist(err) {
				continue
			}
			logger.With(zap.Error(err)).Error("", zap.String("path", dstPath))
			continue
		}

		dstF, err := os.Create(dstPath)
		if err != nil {
			logger.With(zap.Error(err)).Error("", zap.String("path", dstPath))
			continue
		}
		srcF, err := os.Open(file.Name())
		if err != nil {
			logger.With(zap.Error(err)).Error("", zap.String("path", file.Name()))
			continue
		}
		_, err = io.Copy(dstF, srcF)
		if err != nil {
			logger.With(zap.Error(err)).Error("", zap.String("path", file.Name()))
			continue
		}
	}
	return nil
}
