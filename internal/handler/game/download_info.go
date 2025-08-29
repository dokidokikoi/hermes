package game

import (
	"context"
	"encoding/json"
	"io"

	"github.com/dokidokikoi/go-common/errors"
)

type DownloadInfoReq struct {
	GameId uint
}

func (h Handler) DownloadInfo(ctx context.Context, input *DownloadInfoReq) (any, *errors.APIError) {
	if input.GameId > 0 {
		// gVo, err := h.srv.Game().GetVOByID(ctx, input.GameId)

	}
	return nil, nil
}

func downloadInfo(obj any, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	return enc.Encode(obj)
}
