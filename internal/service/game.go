package service

import (
	"context"
	"hermes/db"
	"hermes/internal/handler"
	"hermes/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

type GameWhereNodeFunc func(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption)

type IGame interface {
	CreateL(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff) error
	UpdateL(ctx context.Context, g *model.Game) error

	Search(ctx context.Context, param handler.GameListReq, gwfs ...GameWhereNodeFunc) (int64, []handler.GameVo, error)
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

func (gsrv *game) CreateL(ctx context.Context, g *model.Game, cs []*model.GameCharacter, ss []*model.GameStaff) error {
	tx := gsrv.store.Transaction().Begin()
	err := tx.Game().Create(ctx, g, nil)
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
	err = tx.Character().Creates(ctx, charactersCreate, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	errs := tx.Character().UpdateCollection(ctx, charactersUpdate, nil)
	if len(errs) > 0 {
		tx.Transaction().Rollback()
		return errs[0]
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
		if s.Person.ID != 0 {
			staffCreate = append(staffCreate, s.Person)
		} else {
			staffUpdate = append(staffUpdate, s.Person)
		}
	}
	err = tx.Person().Creates(ctx, staffCreate, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	errs = tx.Person().UpdateCollection(ctx, staffUpdate, nil)
	if len(errs) > 0 {
		tx.Transaction().Rollback()
		return errs[0]
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
	err = tx.GameStaff().Creates(ctx, staff, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}

	return nil
}

func (gsrv *game) UpdateL(ctx context.Context, g *model.Game) error {
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
	err = tx.GameCharacter().Delete(ctx, &model.GameCharacter{GameID: g.ID}, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}

	err = tx.Game().Update(ctx, g, nil)
	if err != nil {
		tx.Transaction().Rollback()
		return err
	}
	return nil
}

func (gsrv *game) Search(ctx context.Context, param handler.GameListReq, gwfs ...GameWhereNodeFunc) (int64, []handler.GameVo, error) {
	head := &meta.WhereNode{}
	node := head
	opt := meta.NewListOption(nil, meta.WithPage(param.Page), meta.WithPageSize(param.PageSize))
	opt.GetOption.Preload = append(opt.GetOption.Preload, []string{"Tags", "Category", "Series"}...)
	for _, f := range gwfs {
		node, opt = f(ctx, param, node, opt)
	}
	gs, err := gsrv.store.Game().ListComplex(ctx, &model.Game{}, node, opt)
	if err != nil {
		return 0, nil, err
	}
	total, err := gsrv.store.Game().CountComplex(ctx, &model.Game{}, node, &opt.GetOption)
	if err != nil {
		return 0, nil, err
	}
	gvos := make([]handler.GameVo, 0, len(gs))
	for _, g := range gs {
		gvos = append(gvos, handler.GameVo{
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

func NewGame(store db.IStore) *game {
	return &game{store: store}
}
