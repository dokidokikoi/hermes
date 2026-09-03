package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"izumi/config"
	"izumi/db"
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

	"github.com/pkg/errors"

	db_pkg "github.com/dokidokikoi/go-common/db"
	comm_errors "github.com/dokidokikoi/go-common/errors"
	"github.com/dokidokikoi/go-common/gopool"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/dokidokikoi/go-common/tools"
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
	UpsertFull(ctx context.Context, gVo *handler.GameVo, gIns *model.GameInstance, proccess func(step int)) error

	Search(ctx context.Context, param handler.GameListReq, opt *meta.ListOption, gwfs ...GameWhereNodeFunc) (int64, []handler.GameVo, error)
	SaveFiles(ctx context.Context, g *handler.GameVo, process ...func(total int)) error

	Load(ctx context.Context, gVo *handler.GameVo, path string) error
	DownloadInfo(ctx context.Context, gameID uint, t time.Time) error
}

var _ IGame = (*game)(nil)

type game struct {
	store db.IStore
}

func (gsrv *game) UpsertFull(ctx context.Context, gVo *handler.GameVo, gIns *model.GameInstance, proccess func(step int)) error {
	var err error
	// err := gsrv.SaveFiles(ctx, gVo, func(total int) {
	// 	step := 4000 / total
	// 	proccess(step)
	// })
	// if err != nil {
	// 	return err
	// }

	db := gsrv.store
	if gVo.ID == 0 {
		gVo.ID, err = gsrv.SearchExistGame(ctx, gVo.RelIDs)
		if err != nil {
			return errors.Wrap(err, "gsrv.SearchExistGame()")
		}
	}

	g := tools.Ptr(handler.Vo2Game(*gVo))
	// 处理category
	{
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
		tagSrv := NewTag(gsrv.store)
		for _, tag := range gVo.Tags {
			err := tagSrv.CreateL(ctx, tag)
			if err != nil {
				return err
			}
			gameTags = append(gameTags, &model.GameTag{
				GameID: g.ID,
				TagID:  tag.ID,
			})
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
			if len(staff.RelIDs) > 0 {
				persons, err := db.Person().ListComplex(
					ctx,
					&model.Person{},
					&meta.WhereNode{Conditions: []*meta.Condition{
						{
							Field:    "rel_ids",
							Operator: "",
							Value:    gorm.Expr(meta.PG_OVERLAP + " " + db_pkg.PgArray(staff.RelIDs)),
						},
					}},
					&meta.ListOption{GetOption: meta.GetOption{Select: []string{"ID"}}})
				if err != nil {
					return err
				}
				if len(persons) > 0 {
					staff.ID = persons[0].ID
				}
			}
			if staff.ID == 0 {
				list, err := db.Person().ListComplex(
					ctx,
					&model.Person{},
					&meta.WhereNode{
						Conditions: []*meta.Condition{
							{
								Field:    "name",
								Operator: meta.EQUAL,
								Value:    staff.Name,
							},
						},
					},
					&meta.ListOption{GetOption: meta.GetOption{Select: []string{"ID"}}})
				if err != nil {
					return err
				}
				if len(list) > 0 {
					staff.ID = list[0].ID
				}
			}
			s := tools.Ptr(handler.Vo2Person(staff))
			if staff.ID != 0 {
				err := db.Person().Update(ctx, s, nil)
				if err != nil {
					return err
				}
			} else {
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
			if len(char.RelIDs) > 0 {
				cs, err := db.Character().ListComplex(ctx, &model.Character{}, &meta.WhereNode{Conditions: []*meta.Condition{
					{
						Field:    "rel_ids",
						Operator: "",
						Value:    gorm.Expr(meta.PG_OVERLAP + " " + db_pkg.PgArray(char.RelIDs)),
					},
				}}, &meta.ListOption{GetOption: meta.GetOption{Select: []string{"ID"}}})
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				if len(cs) > 0 {
					char.ID = cs[0].ID
				}
			}
			c := tools.Ptr(handler.Vo2Character(char))
			if char.CV.Name != "" {
				c.PersonID = staffM[char.CV.Name]
			}
			omits := []string{clause.Associations}
			if char.CV.ID == 0 {
				omits = append(omits, "PersonID")
			}
			if char.ID != 0 {
				err := db.Character().Update(ctx, c, &meta.UpdateOption{Omit: omits})
				if err != nil {
					return err
				}
			} else {
				err := db.Character().Create(ctx, c, &meta.CreateOption{Omit: omits})
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

// runInTx 包装一段事务逻辑：任一步出错即回滚并返回，全部成功则提交。
// 这样调用方无需在每个 DB 操作后重复 Rollback()+return 样板。
func (gsrv *game) runInTx(ctx context.Context, fn func(tx db.IStore) error) (err error) {
	tx := gsrv.store.Transaction().Begin()
	defer func() {
		if p := recover(); p != nil {
			tx.Transaction().Rollback()
			panic(p)
		}
		if err != nil {
			tx.Transaction().Rollback()
			return
		}
		err = tx.Transaction().Commit()
	}()
	return fn(tx)
}

// firstRealErr 从 UpdateCollection 返回的 []error 中取出第一个非 ErrNoUpdateRows 的错误。
func firstRealErr(errs []error) error {
	for _, e := range errs {
		if !errors.Is(e, comm_errors.ErrNoUpdateRows) {
			return e
		}
	}
	return nil
}

// saveGameStory 把 Story 中的外链图片本地化，失败仅记录日志，不中断流程。
func saveGameStory(ctx context.Context, g *model.Game) {
	story, err := SaveHtmlImg(ctx, g.Story)
	if err != nil {
		zaplog.From(ctx).With(zap.Error(err)).Error("SaveHtmlImg")
		return
	}
	g.Story = story
}

// replaceGameBrands 重建游戏与品牌的关联（先删后建）。
func replaceGameBrands(ctx context.Context, tx db.IStore, g *model.Game) error {
	if err := tx.GameBrands().Delete(ctx, &model.GameBrands{GameID: g.ID}, nil); err != nil {
		return err
	}
	gBrands := make([]*model.GameBrands, 0, len(g.Brands))
	for _, b := range g.Brands {
		gBrands = append(gBrands, &model.GameBrands{GameID: g.ID, BrandID: b.ID})
	}
	return tx.GameBrands().Creates(ctx, gBrands, nil)
}

// replaceGameSeries 重建游戏与系列的关联（先删后建）。
func replaceGameSeries(ctx context.Context, tx db.IStore, g *model.Game) error {
	if err := tx.GameSeries().Delete(ctx, &model.GameSeries{GameID: g.ID}, nil); err != nil {
		return err
	}
	gSeries := make([]*model.GameSeries, 0, len(g.Series))
	for _, s := range g.Series {
		gSeries = append(gSeries, &model.GameSeries{GameID: g.ID, SeriesID: s.ID})
	}
	return tx.GameSeries().Creates(ctx, gSeries, nil)
}

// replaceGameCharacters 重建游戏角色：按 ID 分流 create/update，再重建关联表。
func replaceGameCharacters(ctx context.Context, tx db.IStore, gameID uint, cs []*model.GameCharacter) error {
	create, update := []*model.Character{}, []*model.Character{}
	for _, c := range cs {
		if c.Character.ID != 0 {
			update = append(update, c.Character)
		} else {
			create = append(create, c.Character)
		}
	}
	if len(create) > 0 {
		if err := tx.Character().Creates(ctx, create, nil); err != nil {
			return err
		}
	}
	if err := firstRealErr(tx.Character().UpdateCollection(ctx, update, nil)); err != nil {
		return err
	}
	if err := tx.GameCharacter().Delete(ctx, &model.GameCharacter{GameID: gameID}, nil); err != nil {
		return err
	}
	for _, c := range cs {
		c.GameID = gameID
		c.CharacterID = c.Character.ID
	}
	if len(cs) > 0 {
		return tx.GameCharacter().Creates(ctx, cs, nil)
	}
	return nil
}

// replaceGameStaff 重建游戏制作人员：按 ID 分流 create/update，再重建关联表。
func replaceGameStaff(ctx context.Context, tx db.IStore, gameID uint, ss []*model.GameStaff) error {
	create, update := []*model.Person{}, []*model.Person{}
	for _, s := range ss {
		if s.Person.ID == 0 {
			create = append(create, s.Person)
		} else {
			update = append(update, s.Person)
		}
	}
	if len(create) > 0 {
		if err := tx.Person().Creates(ctx, create, nil); err != nil {
			return err
		}
	}
	if err := firstRealErr(tx.Person().UpdateCollection(ctx, update, nil)); err != nil {
		return err
	}
	if err := tx.GameStaff().Delete(ctx, &model.GameStaff{GameID: gameID}, nil); err != nil {
		return err
	}
	for _, s := range ss {
		s.GameID = gameID
		s.PersonID = s.Person.ID
	}
	if len(ss) > 0 {
		return tx.GameStaff().Creates(ctx, ss, nil)
	}
	return nil
}

func (gsrv *game) CreateL(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff, gIns *model.GameInstance) error {
	saveGameStory(ctx, g)

	return gsrv.runInTx(ctx, func(tx db.IStore) error {
		if err := tx.Game().Create(ctx, g, &meta.CreateOption{Omit: []string{"Series", "Brands", "Category"}}); err != nil {
			return err
		}
		if gIns != nil {
			gIns.GameID = g.ID
			if err := tx.GameInstance().Create(ctx, gIns, nil); err != nil {
				return err
			}
		}

		// brands / series（新建时无需先删）
		gBrands := make([]*model.GameBrands, 0, len(g.Brands))
		for _, b := range g.Brands {
			gBrands = append(gBrands, &model.GameBrands{GameID: g.ID, BrandID: b.ID})
		}
		if err := tx.GameBrands().Creates(ctx, gBrands, nil); err != nil {
			return err
		}
		gSeries := make([]*model.GameSeries, 0, len(g.Series))
		for _, s := range g.Series {
			gSeries = append(gSeries, &model.GameSeries{GameID: g.ID, SeriesID: s.ID})
		}
		if err := tx.GameSeries().Creates(ctx, gSeries, nil); err != nil {
			return err
		}

		if err := replaceGameCharacters(ctx, tx, g.ID, cs); err != nil {
			return err
		}
		return replaceGameStaff(ctx, tx, g.ID, ss)
	})
}

func (gsrv *game) UpdateL(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff) error {
	saveGameStory(ctx, g)

	return gsrv.runInTx(ctx, func(tx db.IStore) error {
		// 先清空各关联表
		if err := tx.GameTag().Delete(ctx, &model.GameTag{GameID: g.ID}, nil); err != nil {
			return err
		}
		if err := replaceGameSeries(ctx, tx, g); err != nil {
			return err
		}
		if err := replaceGameBrands(ctx, tx, g); err != nil {
			return err
		}
		if err := replaceGameCharacters(ctx, tx, g.ID, cs); err != nil {
			return err
		}
		if err := replaceGameStaff(ctx, tx, g.ID, ss); err != nil {
			return err
		}
		err := tx.Game().Update(ctx, g, &meta.UpdateOption{Omit: []string{"Series"}})
		if err != nil && !errors.Is(err, comm_errors.ErrNoUpdateRows) {
			return err
		}
		return nil
	})
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
			GetOption: meta.GetOption{Preload: []string{"CV"}},
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
	prMap := map[uint][]model.PersonRelation{}
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
	return tools.Ptr(handler.Game2Vo(*g, cVos, sVos)), nil
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
		if len(g.Tags) > 10 {
			g.Tags = g.Tags[:10]
		}
		gvos = append(gvos, handler.Game2Vo(*g, nil, nil))
	}

	return total, gvos, nil
}

func (gsrv *game) SaveFiles(ctx context.Context, g *handler.GameVo, process ...func(total int)) error {
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
				if len(process) > 0 {
					process[0](len(url))
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
	db := db.GetStore()
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
		c := tools.Ptr(handler.Vo2Character(character))
		cs = append(cs, &model.GameCharacter{
			Character: c,
			Relation:  character.Rlation,
		})

		// 按 rel_ids（["platform:id"]）判断角色是否已导入过：命中则更新已有记录，否则新建
		if len(character.RelIDs) > 0 {
			list, err := db.Character().ListComplex(ctx, &model.Character{}, &meta.WhereNode{Conditions: []*meta.Condition{
				{Field: "rel_ids", Operator: "", Value: gorm.Expr(meta.PG_OVERLAP + " " + db_pkg.PgArray(character.RelIDs))},
			}}, &meta.ListOption{GetOption: meta.GetOption{Select: []string{"ID"}}})
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("get character error: %w", err)
			}
			if len(list) > 0 {
				c.ID = list[0].ID
				charactersUpdate = append(charactersUpdate, c)
				continue
			}
		}
		charactersCreate = append(charactersCreate, c)
	}

	// staff
	var staffIDM = make(map[uint]*model.Person)
	var ss []*model.GameStaff
	staffCreate := []*model.Person{}
	staffUpdate := []*model.Person{}
	for _, staff := range gVo.Staff {
		s := tools.Ptr(handler.Vo2Person(staff))
		ss = append(ss, &model.GameStaff{
			Person:    s,
			Relations: staff.Relation,
		})
		staffIDM[staff.ID] = s

		// 按 rel_ids（["platform:id"]）判断制作人员是否已导入过：命中则更新已有记录，否则新建
		if len(staff.RelIDs) > 0 {
			list, err := db.Person().ListComplex(ctx, &model.Person{}, &meta.WhereNode{Conditions: []*meta.Condition{
				{Field: "rel_ids", Operator: "", Value: gorm.Expr(meta.PG_OVERLAP + " " + db_pkg.PgArray(staff.RelIDs))},
			}}, &meta.ListOption{GetOption: meta.GetOption{Select: []string{"ID"}}})
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("get person error: %w", err)
			}
			if len(list) > 0 {
				s.ID = list[0].ID
				staffUpdate = append(staffUpdate, s)
				continue
			}
		}
		staffCreate = append(staffCreate, s)
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
		for _, err := range errs {
			if !errors.Is(err, comm_errors.ErrNoUpdateRows) {
				return err
			}
		}
	}

	g := tools.Ptr(handler.Vo2Game(*gVo))
	// 按 rel_ids（["platform:id"]）判断游戏是否已导入过：已存在则走更新，否则新建
	if gVo.ID == 0 && len(gVo.RelIDs) > 0 {
		gVo.ID, err = gsrv.SearchExistGame(ctx, gVo.RelIDs)
		if err != nil {
			return err
		}
	}
	if gVo.ID == 0 {
		err = tx.Game().Create(ctx, g, &meta.CreateOption{Omit: []string{clause.Associations}})
		if err != nil {
			return err
		}
		if gVo.Version != "" {
			err = tx.GameInstance().Create(ctx, &model.GameInstance{
				GameID:  g.ID,
				Version: gVo.Version,
				Path:    path,
				Size:    gVo.Size,
				Comment: gVo.Comment,
			}, nil)
			if err != nil {
				return err
			}
		}
	} else {
		g.ID = gVo.ID
		err = tx.Game().Update(ctx, g, &meta.UpdateOption{Omit: []string{clause.Associations}})
		if err != nil {
			return err
		}
		err = tx.GameInstance().Update(ctx, &model.GameInstance{
			GameID:  g.ID,
			Version: gVo.Version,
			Path:    path,
			Size:    gVo.Size,
			Comment: gVo.Comment,
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
	ins, err := db.GetStore().
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
		gVo.Comment = i.Comment
		gVo.Size = i.Size
		err = downloadInfo(gVo, f)
		if err != nil {
			zaplog.From(ctx).Error("downloadInfo", zap.Error(err))
		}
	}

	return nil
}

func (gsrv *game) SearchExistGame(ctx context.Context, relIds []string) (uint, error) {
	if len(relIds) == 0 {
		return 0, nil
	}

	whereNode := &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "rel_ids",
				Operator: "",
				Value:    gorm.Expr(meta.PG_OVERLAP + " " + db_pkg.PgArray(relIds)),
			},
		},
	}

	games, err := gsrv.store.Game().ListComplex(
		ctx,
		&model.Game{},
		whereNode,
		&meta.ListOption{
			PageSize: 1,
		},
	)
	if err != nil {
		return 0, err
	}
	if len(games) == 0 {
		return 0, nil
	}
	return games[0].ID, nil
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
