package game

import (
	"context"
	"encoding/json"
	"hermes/constant"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"
	"os"
	"path/filepath"

	"github.com/dokidokikoi/go-common/core"
	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/dokidokikoi/go-common/notice"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type LoadInfoResponse struct {
	Rid string `json:"rid"`
}

func (h Handler) LoadInfo(ctx context.Context, req *struct{}) (any, error) {
	logger := zaplog.From(ctx)
	p, err := data.GetDataFactory().Policy().Get(ctx, &model.Policy{Key: model.SystemPolicy{}.Key()}, nil)
	if err != nil {
		return nil, err
	}
	sp, err := model.Parse[model.SystemPolicy](p.Policy)
	if err != nil {
		return nil, err
	}
	infos, err := h.srv.Library().Ls(ctx, sp.GameLibrary)
	if err != nil {
		if os.IsNotExist(err) {
			core.WithMsg(ctx, "game library not exist")
		}
		return nil, err
	}
	rid := uuid.NewString()
	go func() {
		var errs []error
		defer func() {
			msg := "success"
			if len(errs) > 0 {
				msg = "failed"
			}
			notice.HubIns.SendBroadcast(constant.TOPIC_SCRAPER, notice.NoticeResponse{
				Rid:     rid,
				Event:   constant.EVENT_INFO_FILE_LOAD,
				Message: msg,
				Data:    nil,
			})
		}()

		for _, info := range infos {
			if !info.IsDir {
				continue
			}
			path := filepath.Join(info.Path, "info.json")
			f, err := os.Open(path)
			if err != nil {
				errs = append(errs, err)
				logger.With(zap.String("path", path)).Error("os.Open", zap.Error(err))
				continue
			}
			defer f.Close()

			gVo := &handler.GameVo{}
			err = json.NewDecoder(f).Decode(gVo)
			if err != nil {
				logger.Error("json.Decode", zap.Error(err))
				continue
			}

			err = h.srv.Game().Load(ctx, gVo, path)
			if err != nil {
				logger.With(zap.Error(err)).With(zap.String("path", info.Path)).Error("load failed")
				continue
			}
		}
	}()

	return nil, nil
}
