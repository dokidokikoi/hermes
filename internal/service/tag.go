package service

import (
	"context"
	"fmt"
	"izumi/db"
	"izumi/model"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ITag interface {
	CreateL(ctx context.Context, tag *model.Tag) error
}

var _ ITag = (*tag)(nil)

type tag struct {
	store db.IStore
}

func (t *tag) CreateL(ctx context.Context, tag *model.Tag) error {
	if tag.NS == "" {
		tag.NS = model.NS_RECLASS
	}
	if tag.Key == "" {
		tag.Key = tag.Name
	}
	oriTag, err := t.store.Tag().Get(ctx, &model.Tag{NS: tag.NS, Name: tag.Name}, nil)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("get original tag error: %w", err)
		}
	} else {
		tag.ID = oriTag.ID
	}
	if tag.NS == model.NS_RECLASS {
		// 之前已经处理过的 tag 直接复用之前的 tag 信息
		dtag, err := t.store.DecidedTag().Get(ctx, &model.DecidedTag{Tag: tag.Name}, nil)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if oriTag == nil {
					err = t.store.Tag().Create(ctx, tag, nil)
					if err != nil {
						return fmt.Errorf("create tag error: %w", err)
					}
				}
				err = t.store.DecidedTag().Create(ctx, &model.DecidedTag{
					Tag:   tag.Name,
					TagID: tag.ID,
				}, nil)
				if err != nil {
					zaplog.L().Error("store.DecidedTag().Create()", zap.Error(err))
				}
				return nil
			}
			return fmt.Errorf("get decided tag error: %w", err)
		}
		refTag, err := t.store.Tag().Get(ctx, &model.Tag{ID: dtag.TagID}, nil)
		if err != nil {
			return nil
		}
		*tag = *refTag
	} else {
		if oriTag == nil {
			err = t.store.Tag().Create(ctx, tag, nil)
			if err != nil {
				return fmt.Errorf("create tag error: %w", err)
			}
		}
	}

	return nil
}

func NewTag(store db.IStore) *tag {
	return &tag{store: store}
}
