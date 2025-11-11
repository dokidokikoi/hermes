package service

import (
	"context"
	"izumi/db"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

var libraryCache = *gocache.New(time.Minute, time.Hour)

type PathInfo struct {
	Path  string     `json:"path"`
	IsDir bool       `json:"is_dir"`
	Child []PathInfo `json:"child"`
}

type ILibrary interface {
	Ls(ctx context.Context, path string, onlyNoScrap bool) ([]PathInfo, error)
}

var _ ILibrary = (*library)(nil)

type library struct {
	store db.IStore
}

func (lsrv *library) Ls(ctx context.Context, path string, onlyNoScrap bool) ([]PathInfo, error) {
	ls, ok := libraryCache.Get(path)
	if ok {
		return ls.([]PathInfo), nil
	}

	paths, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var result []PathInfo
	for _, p := range paths {
		if strings.HasPrefix(p.Name(), ".") {
			continue
		}
		if onlyNoScrap {
			if !p.IsDir() {
				continue
			}
			if _, err := os.Stat(filepath.Join(path, p.Name(), "info.json")); err == nil {
				continue
			}
		}
		child := []PathInfo{}
		if p.IsDir() {
			cPaths, err := os.ReadDir(filepath.Join(path, p.Name()))
			if err != nil {
				return nil, err
			}
			for _, cp := range cPaths {
				if strings.HasPrefix(cp.Name(), ".") {
					continue
				}
				child = append(child, PathInfo{
					Path:  filepath.Join(path, p.Name(), cp.Name()),
					IsDir: p.IsDir(),
				})
				sort.Slice(child, func(i, j int) bool {
					if child[i].IsDir == child[j].IsDir {
						return strings.Compare(strings.ToLower(child[i].Path), strings.ToLower(child[j].Path)) <= 0
					} else {
						return child[i].IsDir
					}
				})
			}
		}
		result = append(result, PathInfo{
			Path:  filepath.Join(path, p.Name()),
			IsDir: p.IsDir(),
			Child: child,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].IsDir == result[j].IsDir {
			return strings.Compare(strings.ToLower(result[i].Path), strings.ToLower(result[j].Path)) <= 0
		} else {
			return result[i].IsDir
		}
	})

	libraryCache.SetDefault(path, result)
	return result, nil
}

func NewLibrary(store db.IStore) *library {
	return &library{store: store}
}
