package game

import (
	"context"
	"izumi/db"
	"izumi/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/middleware"
)

type PanelResp struct {
	Games      int64 `json:"games"`
	LocalGames int64 `json:"local_games"`
	Tags       int64 `json:"tags"`
	UsedTags   int64 `json:"used_tags"`
}

func (h *Handler) Panel(ctx context.Context, req *struct{}, op *middleware.PreHandleOptions) (*PanelResp, error) {
	games, err := db.GetStore().Game().Count(ctx, &model.Game{}, nil)
	if err != nil {
		return nil, err
	}
	p, err := db.GetStore().Policy().Get(context.Background(), &model.Policy{Key: model.SystemPolicy{}.Key()}, nil)
	if err != nil {
		return nil, err
	}
	sp, err := model.Parse[model.SystemPolicy](p.Policy)
	if err != nil {
		return nil, err
	}
	var localGames int64
	for _, path := range sp.GameLibrary {
		infos, err := h.srv.Library().Ls(ctx, path, false, true)
		if err != nil {
			return nil, err
		}
		localGames += int64(len(infos))
	}
	tags, err := db.GetStore().Tag().Count(ctx, &model.Tag{}, nil)
	if err != nil {
		return nil, err
	}
	usedTags, err := db.GetStore().GameTag().Count(ctx, &model.GameTag{}, &meta.GetOption{
		Group: meta.Group{
			Fields: []string{"tag_id"},
		},
	})
	if err != nil {
		return nil, err
	}

	return &PanelResp{
		Games:      games,
		LocalGames: localGames,
		Tags:       tags,
		UsedTags:   usedTags,
	}, nil
}
