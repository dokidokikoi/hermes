package game

import (
	"context"
)

type LoadInfoResponse struct {
	Rid string `json:"rid"`
}

func (h Handler) LoadInfo(ctx context.Context, req *struct{}) (any, error) {
	// p, err := data.GetDataFactory().Policy().Get(ctx, &model.Policy{Key: model.SystemPolicy{}.Key()}, nil)
	// if err != nil {
	// 	return nil, err
	// }
	// sp, err := model.Parse[model.SystemPolicy](p.Policy)
	// if err != nil {
	// 	return nil, err
	// }
	// infos, err := h.srv.Library().Ls(ctx, sp.GameLibrary)
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		core.WithMsg(ctx, "game library not exist")
	// 	}
	// 	return nil, err
	// }
	// rid := uuid.NewString()
	// go func() {
	// 	for _, info := range infos {
	// 		if !info.IsDir {
	// 			continue
	// 		}
	// 		f, err := os.Open(filepath.Join(info.Path, "info.json"))
	// 		if err != nil {
	// 			continue
	// 		}
	// 		defer f.Close()

	// 		json.NewDecoder(f).Decode()
	// 	}
	// }()

	return nil, nil
}
