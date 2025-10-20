package service

import (
	"context"
	"encoding/json"
	"errors"
	"hermes/db"
	"hermes/model"
	"hermes/scraper/event"

	"gorm.io/gorm"
)

type IPolicy interface {
	SystemPolicyEffect(ctx context.Context, sp *model.SystemPolicy)
	ScraperPolicyEffect(ctx context.Context, sp *model.ScraperPolicy)
	PolicyEffect(ctx context.Context) error
}

var _ IPolicy = (*policy)(nil)

type policy struct {
	store db.IStore
}

func (p *policy) SystemPolicyEffect(ctx context.Context, sp *model.SystemPolicy) {
}

func (p *policy) ScraperPolicyEffect(ctx context.Context, sp *model.ScraperPolicy) {
	if sp == nil {
		return
	}
	if len(event.GameScraperMap) == 0 {
		for name, c := range event.GameScraperConstructors {
			if s, ok := (*sp)[name]; ok {
				event.RegisterScraper(c(s.Header, s.Proxy))
			}
		}
	} else {
		for name, scraper := range event.GameScraperMap {
			if s, ok := (*sp)[name]; ok {
				scraper.SetHeader(s.Header)
				scraper.SetProxy(s.Proxy)
			}
		}
	}
}

func (p *policy) PolicyEffect(ctx context.Context) error {
	po, err := p.store.Policy().Get(context.Background(), &model.Policy{Key: model.SystemPolicy{}.Key()}, nil)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		} else {
			d, err := json.Marshal(model.SystemPolicy{})
			if err != nil {
				return err
			}
			err = p.store.Policy().Create(ctx, &model.Policy{Key: model.SystemPolicy{}.Key(), Policy: string(d)}, nil)
			if err != nil {
				return err
			}
		}
	} else {
		sp, err := model.Parse[model.SystemPolicy](po.Policy)
		if err != nil {
			return err
		}
		p.SystemPolicyEffect(ctx, sp)
	}

	po, err = p.store.Policy().Get(context.Background(), &model.Policy{Key: model.ScraperPolicy{}.Key()}, nil)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		} else {
			d, err := json.Marshal(event.GameScraperPolicyMap)
			if err != nil {
				return err
			}
			po = &model.Policy{Key: model.ScraperPolicy{}.Key(), Policy: string(d)}
			err = p.store.Policy().Create(ctx, po, nil)
			if err != nil {
				return err
			}
		}
	}
	sp, err := model.Parse[model.ScraperPolicy](po.Policy)
	if err != nil {
		return err
	}
	p.ScraperPolicyEffect(ctx, sp)

	return nil
}

func NewPolicy(store db.IStore) *policy {
	return &policy{store: store}
}
