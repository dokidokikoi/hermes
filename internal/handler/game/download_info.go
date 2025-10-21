package game

import (
	"context"
	"encoding/json"
	"hermes/config"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"
	"hermes/tools"
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
			err = cpGameAllImages(logger, i.Path, gVo)
			if err != nil {
				logger.Error("cpGameAllImages", zap.Error(err))
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
			defer f.Close()

			gVo.Version = i.Version
			gVo.Language = i.Language
			gVo.Comment = i.Comment
			gVo.Size = i.Size
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

func cpGameAllImages(logger *zap.Logger, path string, gVo *handler.GameVo) error {
	if gVo.Cover != "" {
		err := tools.Cp(filepath.Join(config.DataDir, gVo.Cover), filepath.Join(path, gVo.Cover))
		if err != nil {
			logger.With(zap.Error(err)).Sugar().Errorf("tools.Cp(%s)", filepath.Join(path, gVo.Cover))
		}
	}
	for _, image := range gVo.Images {
		if image != "" {
			err := tools.Cp(filepath.Join(config.DataDir, image), filepath.Join(path, image))
			if err != nil {
				logger.With(zap.Error(err)).Sugar().Errorf("tools.Cp(%s)", filepath.Join(path, image))
			}
		}
	}
	for _, c := range gVo.Characters {
		if c.Cover != "" {
			err := tools.Cp(filepath.Join(config.DataDir, c.Cover), filepath.Join(path, c.Cover))
			if err != nil {
				logger.With(zap.Error(err)).Sugar().Errorf("tools.Cp(%s)", filepath.Join(path, c.Cover))
			}
		}
		for _, image := range c.Images {
			if image != "" {
				err := tools.Cp(filepath.Join(config.DataDir, image), filepath.Join(path, image))
				if err != nil {
					logger.With(zap.Error(err)).Sugar().Errorf("tools.Cp(%s)", filepath.Join(path, image))
				}
			}
		}
	}
	for _, s := range gVo.Staff {
		if s.Cover != "" {
			err := tools.Cp(filepath.Join(config.DataDir, s.Cover), filepath.Join(path, s.Cover))
			if err != nil {
				logger.With(zap.Error(err)).Sugar().Errorf("tools.Cp(%s)", filepath.Join(path, s.Cover))
			}
		}
		for _, image := range s.Images {
			if image != "" {
				err := tools.Cp(filepath.Join(config.DataDir, image), filepath.Join(path, image))
				if err != nil {
					logger.With(zap.Error(err)).Sugar().Errorf("tools.Cp(%s)", filepath.Join(path, image))
				}
			}
		}
	}
	return nil
}
