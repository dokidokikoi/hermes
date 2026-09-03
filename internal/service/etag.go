package service

import (
	"context"
	"fmt"
	"izumi/db"
	"izumi/model"
	"izumi/scraper/ehtag"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IEtag interface {
	// SyncEtags 拉取 EhTagTranslation 数据库并同步到 etags 表
	// （新增缺失条目、更新变化的译名/介绍）。
	SyncEtags(ctx context.Context) (*EtagSyncResult, error)
	// Decide 手动匹配：将 decided_tag 指向的 tags 行更新为规范 etag 的
	// ns/key/name/intro。tag ID 不变，game_tags 引用无需迁移。
	Decide(ctx context.Context, raw string, ns, key string) error
}

var _ IEtag = (*etagSrv)(nil)

// 同步批量写入的分片大小，避免单条 SQL 参数过多。
const etagBatchSize = 500

type etagSrv struct {
	store db.IStore
}

type EtagSyncResult struct {
	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
	Total   int   `json:"total"`
}

func (s *etagSrv) SyncEtags(ctx context.Context) (*EtagSyncResult, error) {
	database, err := ehtag.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	existing, err := s.store.Etag().List(ctx, &model.Etag{}, nil)
	if err != nil {
		return nil, fmt.Errorf("list etags error: %w", err)
	}
	type etagKey struct{ ns, key string }
	existMap := make(map[etagKey]*model.Etag, len(existing))
	for _, e := range existing {
		existMap[etagKey{e.NS, e.Key}] = e
	}

	result := &EtagSyncResult{}
	var (
		creates []*model.Etag
		updates []*model.Etag
	)
	for _, entries := range database {
		for _, entry := range entries {
			result.Total++
			old, ok := existMap[etagKey{entry.NS, entry.Key}]
			if !ok {
				creates = append(creates, &model.Etag{
					NS:    entry.NS,
					Key:   entry.Key,
					Name:  entry.Name,
					Intro: entry.Intro,
				})
				continue
			}
			if old.Name != entry.Name || old.Intro != entry.Intro {
				old.Name = entry.Name
				old.Intro = entry.Intro
				updates = append(updates, old)
			}
		}
	}

	for i := 0; i < len(creates); i += etagBatchSize {
		end := min(i+etagBatchSize, len(creates))
		if err := s.store.Etag().Creates(ctx, creates[i:end], nil); err != nil {
			return nil, fmt.Errorf("create etags error: %w", err)
		}
		result.Created += int64(end - i)
	}
	if len(updates) > 0 {
		errs := s.store.Etag().UpdateCollection(ctx, updates, nil)
		for _, e := range errs {
			if e != nil {
				return nil, fmt.Errorf("update etags error: %w", e)
			}
		}
		result.Updated = int64(len(updates))
	}
	return result, nil
}

func (s *etagSrv) Decide(ctx context.Context, raw string, ns, key string) error {
	et, err := s.store.Etag().Get(ctx, &model.Etag{NS: ns, Key: key}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("etag %s:%s not found, sync etags first", ns, key)
		}
		return fmt.Errorf("get etag error: %w", err)
	}

	dtag, err := s.store.DecidedTag().Get(ctx, &model.DecidedTag{Tag: raw}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("decided tag %q not found", raw)
		}
		return fmt.Errorf("get decided tag error: %w", err)
	}

	// 用规范 etag 信息覆盖 tags 行；Select 确保 Intro 空串也能写入。
	err = s.store.Tag().Update(ctx, &model.Tag{
		ID:    dtag.TagID,
		NS:    et.NS,
		Key:   et.Key,
		Name:  et.Name,
		Intro: et.Intro,
	}, &meta.UpdateOption{Select: []string{"ns", "key", "name", "intro", "lang"}})
	if err != nil {
		return fmt.Errorf("update tag error: %w", err)
	}
	return nil
}

func NewEtag(store db.IStore) *etagSrv {
	return &etagSrv{store: store}
}
