package game

import (
	"context"
	"encoding/json"
	"hermes/db/data"
	"hermes/model"
	"io"
	"os"
	"path/filepath"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	meta "github.com/dokidokikoi/go-common/meta/option"
	"go.uber.org/zap"
)

type DownloadInfoReq struct {
	GameId uint
}

func (h Handler) DownloadInfo(ctx context.Context, input *DownloadInfoReq) (any, error) {
	logger := zaplog.From(ctx)
	if input.GameId > 0 {
		gVo, err := h.srv.Game().GetVOByID(ctx, input.GameId)
		if err != nil {
			return nil, err
		}
		ins, err := data.GetDataFactory().
			GameInstance().List(ctx, &model.GameInstance{GameID: input.GameId}, nil)
		if err != nil {
			return nil, err
		}
		for _, i := range ins {
			info, err := os.Stat(i.Path)
			if err != nil {
				logger.Error("os.Stat", zap.Error(err))
				continue
			}
			if !info.IsDir() {
				continue
			}
			f, err := os.OpenFile(filepath.Join(i.Path, "info.json"), os.O_WRONLY, 0666)
			if err != nil {
				if os.IsNotExist(err) {
					err = nil
					f, err = os.Create(filepath.Join(i.Path, "info.json"))
				}
				if err != nil {
					logger.Error("Open", zap.Error(err))
					continue
				}
			}
			gVo.Version = i.Version
			gVo.Language = i.Language
			gVo.Comment = i.Comment
			err = downloadInfo(gVo, f)
			if err != nil {
				logger.Error("downloadInfo", zap.Error(err))
			}
		}
	}
	return nil, nil
}

func (h Handler) DownloadAllInfo(ctx context.Context, input struct{}) (any, error) {
	logger := zaplog.From(ctx)
	gs, err := data.GetDataFactory().Game().List(ctx, &model.Game{}, &meta.ListOption{GetOption: meta.GetOption{Select: []string{"ID"}}})
	if err != nil {
		return nil, err
	}
	go func() {
		for _, g := range gs {
			_, e := h.DownloadInfo(ctx, &DownloadInfoReq{GameId: g.ID})
			if e != nil {
				logger.Error("downloadInfo", zap.Any("error", e))
			}
		}
	}()

	return nil, nil
}

func downloadInfo(obj any, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	return enc.Encode(obj)
}
