package service

import (
	"context"
	"hermes/db"
	"hermes/internal/handler"
	"hermes/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

type CharacterWhereNodeFunc func(ctx context.Context, param handler.CharacterListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption)

type ICharacter interface {
	Search(ctx context.Context, param handler.CharacterListReq, opt *meta.ListOption, cwfs ...CharacterWhereNodeFunc) (int64, []*model.Character, error)
}

var _ ICharacter = (*character)(nil)

type character struct {
	store db.IStore
}

func (csrv *character) Search(ctx context.Context, param handler.CharacterListReq, opt *meta.ListOption, cwfs ...CharacterWhereNodeFunc) (int64, []*model.Character, error) {
	head := &meta.WhereNode{}
	node := head
	if opt == nil {
		opt = meta.NewListOption(nil)
	}
	opt.GetOption.Preload = append(opt.GetOption.Preload, []string{"Tags", "CV"}...)
	for _, f := range cwfs {
		node, opt = f(ctx, param, node, opt)
	}
	cs, err := csrv.store.Character().ListComplex(ctx, &model.Character{}, node, opt)
	if err != nil {
		return 0, nil, err
	}
	total, err := csrv.store.Character().CountComplex(ctx, &model.Character{}, node, &opt.GetOption)
	if err != nil {
		return 0, nil, err
	}
	return total, cs, err
}

func NewCharacter(store db.IStore) *character {
	return &character{store: store}
}
