package tag

import (
	"context"
	"izumi/db"
	"izumi/internal/service"
	"izumi/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/dokidokikoi/go-common/middleware"
)

type EtagSyncResponse struct {
	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
	Total   int   `json:"total"`
}

// Sync 拉取 EhTagTranslation 数据库并同步到 etags 表。
func (h Handler) EtagSync(ctx context.Context, req *struct{}, op *middleware.PreHandleOptions) (*EtagSyncResponse, error) {
	result, err := service.NewEtag(db.GetStore()).SyncEtags(ctx)
	if err != nil {
		return nil, err
	}
	return &EtagSyncResponse{
		Created: result.Created,
		Updated: result.Updated,
		Total:   result.Total,
	}, nil
}

// EtagSearchRequest 按 ns / 关键字（key 或中文译名模糊匹配）搜索 etags，
// 供手动匹配时挑选规范 tag。
type EtagSearchRequest struct {
	NS      string `json:"ns" form:"ns"`
	Keyword string `json:"keyword" form:"keyword"`
}

// EtagSearch 搜索 etags 表。同一节点内的条件按 OR 连接，
// ns 过滤通过 Next 节点以 AND 挂接。
func (h Handler) EtagSearch(ctx context.Context, req *EtagSearchRequest, op *middleware.PreHandleOptions) ([]*model.Etag, error) {
	node := &meta.WhereNode{}
	if req.NS != "" {
		node.Conditions = append(node.Conditions, &meta.Condition{
			Field:    "ns",
			Operator: meta.EQUAL,
			Value:    req.NS,
		})
	}
	if req.Keyword != "" {
		node.Next = &meta.WhereNode{
			Conditions: []*meta.Condition{
				{
					Field:    "key",
					Operator: meta.LIKE,
					Value:    "%" + req.Keyword + "%",
				},
				{
					Field:    "name",
					Operator: meta.LIKE,
					Value:    "%" + req.Keyword + "%",
				},
			},
		}
	}
	return db.GetStore().Etag().ListComplex(ctx, &model.Etag{}, node,
		&meta.ListOption{GetOption: meta.GetOption{Select: []string{"id", "ns", "key", "name"}}, PageSize: 50, Page: 1})
}

// DecidedTagResponse 是一条待匹配的原始标签及其当前 tags 行。
type DecidedTagResponse struct {
	Tag    string     `json:"tag"`
	TagID  uint       `json:"tag_id"`
	Detail *model.Tag `json:"detail"`
}

// ListDecided 列出所有 decided_tags（ns=reclass 的 detail 即待人工匹配）。
func (h Handler) ListDecided(ctx context.Context, req *struct{}, op *middleware.PreHandleOptions) ([]*DecidedTagResponse, error) {
	dtags, err := db.GetStore().DecidedTag().List(ctx, &model.DecidedTag{}, &meta.ListOption{Order: "tag"})
	if err != nil {
		return nil, err
	}
	resp := make([]*DecidedTagResponse, 0, len(dtags))
	ids := make([]uint, 0, len(dtags))
	for _, d := range dtags {
		resp = append(resp, &DecidedTagResponse{Tag: d.Tag, TagID: d.TagID})
		ids = append(ids, d.TagID)
	}
	if len(ids) > 0 {
		tags, err := db.GetStore().Tag().ListComplex(ctx, &model.Tag{}, &meta.WhereNode{
			Conditions: []*meta.Condition{{Field: "id", Operator: meta.IN, Value: ids}},
		}, nil)
		if err != nil {
			return nil, err
		}
		tmap := make(map[uint]*model.Tag, len(tags))
		for _, t := range tags {
			tmap[t.ID] = t
		}
		for _, r := range resp {
			r.Detail = tmap[r.TagID]
		}
	}
	return resp, nil
}

// DecideRequest 提交一次手动匹配：原始标签 raw -> 规范 etag (ns, key)。
type DecideRequest struct {
	Tag string `json:"tag" binding:"required"`
	NS  string `json:"ns" binding:"required"`
	Key string `json:"key" binding:"required"`
}

// Decide 将 decided_tag 指向的 tags 行更新为规范 etag 信息。
func (h Handler) Decide(ctx context.Context, req *DecideRequest, op *middleware.PreHandleOptions) (*struct{}, error) {
	return nil, service.NewEtag(db.GetStore()).Decide(ctx, req.Tag, req.NS, req.Key)
}
