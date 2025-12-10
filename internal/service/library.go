package service

import (
	"context"
	"izumi/db"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type PathInfo struct {
	Path  string     `json:"path"`
	IsDir bool       `json:"is_dir"`
	Child []PathInfo `json:"child"`
}

type ILibrary interface {
	Ls(ctx context.Context, path string, onlyNoScrap bool, cache bool) ([]PathInfo, error)
}

var _ ILibrary = (*library)(nil)

type library struct {
	store db.IStore
}

func (lsrv *library) Ls(ctx context.Context, path string, onlyNoScrap bool, cache bool) ([]PathInfo, error) {
	if cache {
		defer func() {
			select {
			case libraryCacheChan <- path:
			default:
			}
		}()
		ls, ok := libraryCache.Get(path)
		if ok {
			return filterLibrary(ls.([]PathInfo), onlyNoScrap), nil
		}
		return nil, nil
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
					IsDir: cp.IsDir(),
				})
			}
			sort.Slice(child, func(i, j int) bool {
				if child[i].IsDir == child[j].IsDir {
					return strings.Compare(strings.ToLower(child[i].Path), strings.ToLower(child[j].Path)) <= 0
				} else {
					return child[i].IsDir
				}
			})
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

	return filterLibrary(result, onlyNoScrap), nil
}

func filterLibrary(libs []PathInfo, onlyNoScrap bool) []PathInfo {
	if !onlyNoScrap {
		return libs
	}
	res := []PathInfo{}
loop:
	for _, l := range libs {
		if !l.IsDir {
			continue
		}
		for _, c := range l.Child {
			if filepath.Base(c.Path) == "info.json" {
				continue loop
			}
		}
		res = append(res, l)
	}
	return res
}

func NewLibrary(store db.IStore) *library {
	return &library{store: store}
}
