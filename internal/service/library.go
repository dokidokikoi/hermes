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
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
}

type ILibrary interface {
	Ls(ctx context.Context, path string, onlyNoScrap bool) ([]PathInfo, error)
}

var _ ILibrary = (*library)(nil)

type library struct {
	store db.IStore
}

func (lsrv *library) Ls(ctx context.Context, path string, onlyNoScrap bool) ([]PathInfo, error) {
	paths, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var result []PathInfo
	if onlyNoScrap {
		for _, p := range paths {
			if strings.HasPrefix(p.Name(), ".") {
				continue
			}
			if !p.IsDir() {
				continue
			}
			if _, err := os.Stat(filepath.Join(path, p.Name(), "info.json")); err == nil {
				continue
			}
			result = append(result, PathInfo{
				Path:  filepath.Join(path, p.Name()),
				IsDir: p.IsDir(),
			})
		}
	} else {
		for _, p := range paths {
			if strings.HasPrefix(p.Name(), ".") {
				continue
			}
			result = append(result, PathInfo{
				Path:  filepath.Join(path, p.Name()),
				IsDir: p.IsDir(),
			})
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].IsDir == result[j].IsDir {
			return strings.Compare(strings.ToLower(result[i].Path), strings.ToLower(result[j].Path)) <= 0
		} else {
			return result[i].IsDir
		}
	})
	return result, nil
}

func NewLibrary(store db.IStore) *library {
	return &library{store: store}
}
