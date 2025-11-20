package service

import (
	"context"
	"izumi/db/data"
	"izumi/model"
	"time"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	gocache "github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

var (
	libraryCache = *gocache.New(time.Hour, time.Hour)

	libraryCacheChan = make(chan string, 1)
)

func StartLibraryCache() {
	ticker := time.NewTicker(time.Hour)
	db := data.GetDataFactory()
	srv := NewLibrary(db)

	err := DoSetCache(srv)
	if err != nil {
		zaplog.L().With(zap.String("fn", "StartLibraryCache")).Error("set cache error", zap.Error(err))
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				err := DoSetCache(srv)
				if err != nil {
					zaplog.L().With(zap.String("fn", "StartLibraryCache")).Error("set cache error", zap.Error(err))
				}
			case path := <-libraryCacheChan:
				err := SetCache(srv, path)
				if err != nil {
					zaplog.L().With(zap.String("fn", "StartLibraryCache")).Error("set cache error", zap.Error(err))
				}
			}
		}
	}()
}

func DoSetCache(srv ILibrary) error {
	p, err := data.GetDataFactory().Policy().Get(context.Background(), &model.Policy{Key: model.SystemPolicy{}.Key()}, nil)
	if err != nil {
		return err
	}
	sp, err := model.Parse[model.SystemPolicy](p.Policy)
	if err != nil {
		return err
	}
	err = SetCache(srv, sp.GameLibrary...)
	if err != nil {
		return err
	}
	return nil
}

func SetCache(srv ILibrary, paths ...string) error {
	for _, path := range paths {
		ls, err := srv.Ls(context.Background(), path, false, false)
		if err != nil {
			return err
		}
		libraryCache.SetDefault(path, ls)
	}
	return nil
}
