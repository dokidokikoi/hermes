package service

import (
	"context"
	"hermes/db"
	"hermes/internal/handler"
	"hermes/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

type PersonWhereNodeFunc func(ctx context.Context, param handler.PersonListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption)

type IPerson interface {
	Search(ctx context.Context, param handler.PersonListReq, opt *meta.ListOption, pwfs ...PersonWhereNodeFunc) (int64, []*model.Person, error)
}

var _ IPerson = (*person)(nil)

type person struct {
	store db.IStore
}

func (psrv *person) Search(ctx context.Context, param handler.PersonListReq, opt *meta.ListOption, pwfs ...PersonWhereNodeFunc) (int64, []*model.Person, error) {
	head := &meta.WhereNode{}
	node := head
	if opt == nil {
		opt = meta.NewListOption(nil)
	}
	opt.GetOption.Preload = append(opt.GetOption.Preload, []string{"Tags"}...)
	for _, f := range pwfs {
		node, opt = f(ctx, param, node, opt)
	}
	cs, err := psrv.store.Person().ListComplex(ctx, &model.Person{}, node, opt)
	if err != nil {
		return 0, nil, err
	}
	total, err := psrv.store.Person().CountComplex(ctx, &model.Person{}, node, &opt.GetOption)
	if err != nil {
		return 0, nil, err
	}
	return total, cs, err
}

func NewPerson(store db.IStore) *person {
	return &person{store: store}
}
