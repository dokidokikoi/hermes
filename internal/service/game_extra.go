package service

import (
	"context"
	"fmt"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"
	"strings"

	meta "github.com/dokidokikoi/go-common/meta/option"
	comm_tools "github.com/dokidokikoi/go-common/tools"
)

var GameBasicSearchNode = []GameWhereNodeFunc{
	GameWhereNodeKeyword,
	GameWhereNodeTag,
	GameWhereNodeCharacter,
	GameWhereNodeStaff,
	GameWhereNodeSeries,
	GameWhereNodeCategory,
	GameWhereNodeDeveloper,
	GameWhereNodePublisher,
	GameWhereNodeSizeRange,
	GameWhereNodeIssueDateRange,
	GameWhereNodeCreatedAtRange,
}

func GameWhereNodeKeyword(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	keyword := comm_tools.TrimBlankChar(param.Keyword)
	if keyword != "" {
		node.Next = &meta.WhereNode{
			Conditions: []*meta.Condition{
				{
					Field:    "name",
					Operator: meta.LIKE,
					Value:    fmt.Sprintf("%%%s%%", keyword),
				},
				{
					Field:    "alias::text",
					Operator: meta.LIKE,
					Value:    fmt.Sprintf("%%%s%%", keyword),
				},
			},
		}
		if param.FullText {
			node.Next.Conditions = append(node.Next.Conditions, []*meta.Condition{
				{
					Field:    "story",
					Operator: meta.LIKE,
					Value:    fmt.Sprintf("%%%s%%", keyword),
				},
				{
					Field:    "other_info",
					Operator: meta.LIKE,
					Value:    fmt.Sprintf("%%%s%%", keyword),
				},
			}...)
		}
		node = node.Next
	}
	return node, opt
}
func GameWhereNodeTag(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	if len(param.Tags) < 1 {
		return node, opt
	}

	defer func() {
		o = opt
	}()
	for i := range param.Tags {
		param.Tags[i] = strings.ToLower(param.Tags[i])
	}
	tmpdb := data.GetDataFactory().Tag().ListComplexDB(ctx, &model.Tag{}, &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "LOWER(name)",
				Operator: meta.IN,
				Value:    param.Tags,
			},
		},
	}, nil)
	db := data.GetDataFactory().GameTag().ListDB(ctx, &model.GameTag{}, &meta.ListOption{
		GetOption: meta.GetOption{
			Join: []*meta.Join{
				{
					Method:          meta.INNER_JOIN,
					Table:           model.GameTag{}.TableName(),
					InnerQuery:      tmpdb,
					InnerQueryAlias: "t",
					TableField:      "tag_id",
					JoinTableField:  "id",
				},
			},
		},
	})
	opt.GetOption.Join = append(opt.GetOption.Join, &meta.Join{
		Method:          meta.INNER_JOIN,
		Table:           model.Game{}.TableName(),
		InnerQuery:      db,
		InnerQueryAlias: "game_tag",
		TableField:      "id",
		JoinTableField:  "game_id",
	})
	return
}
func GameWhereNodeCharacter(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	if param.Character == 0 {
		return node, opt
	}
	db := data.GetDataFactory().GameCharacter().ListDB(ctx, &model.GameCharacter{CharacterID: param.Character}, nil)
	opt.GetOption.Join = append(opt.GetOption.Join, &meta.Join{
		Method:          meta.INNER_JOIN,
		Table:           model.Game{}.TableName(),
		InnerQuery:      db,
		InnerQueryAlias: "game_character",
		TableField:      "id",
		JoinTableField:  "game_id",
	})
	return node, opt
}
func GameWhereNodeStaff(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption) {
	if param.Staff == 0 {
		return node, opt
	}
	db := data.GetDataFactory().GameStaff().ListDB(ctx, &model.GameStaff{PersonID: param.Staff}, nil)
	opt.GetOption.Join = append(opt.GetOption.Join, &meta.Join{
		Method:          meta.INNER_JOIN,
		Table:           model.Game{}.TableName(),
		InnerQuery:      db,
		InnerQueryAlias: "game_staff",
		TableField:      "id",
		JoinTableField:  "game_id",
	})
	return node, opt
}
func GameWhereNodeSeries(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption) {
	if param.Series == 0 {
		return node, opt
	}
	db := data.GetDataFactory().GameSeries().ListDB(ctx, &model.GameSeries{SeriesID: param.Series}, nil)
	opt.GetOption.Join = append(opt.GetOption.Join, &meta.Join{
		Method:          meta.INNER_JOIN,
		Table:           model.Game{}.TableName(),
		InnerQuery:      db,
		InnerQueryAlias: "game_series",
		TableField:      "id",
		JoinTableField:  "game_id",
	})
	return node, opt
}
func GameWhereNodeCategory(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	if param.Category == 0 {
		return node, opt
	}
	defer func() {
		n = node.Next
		o = opt
	}()

	node.Next = &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "category_id",
				Operator: meta.EQUAL,
				Value:    param.Category,
			},
		},
	}
	return
}
func GameWhereNodeDeveloper(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	if param.Developer == 0 {
		return node, opt
	}
	defer func() {
		n = node.Next
		o = opt
	}()

	node.Next = &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "developer_id",
				Operator: meta.EQUAL,
				Value:    param.Developer,
			},
		},
	}
	return
}
func GameWhereNodePublisher(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	if param.Publisher == 0 {
		return node, opt
	}
	defer func() {
		n = node.Next
		o = opt
	}()

	node.Next = &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "publisher_id",
				Operator: meta.EQUAL,
				Value:    param.Publisher,
			},
		},
	}
	return
}
func GameWhereNodeSizeRange(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	if len(param.SizeRange) > 0 {
		node.Next = &meta.WhereNode{
			Conditions: []*meta.Condition{
				{
					Field:    "size",
					Operator: meta.GTE,
					Value:    param.SizeRange[0],
				},
			},
		}

		node = node.Next
		if len(param.SizeRange) > 1 {
			node.Next = &meta.WhereNode{
				Conditions: []*meta.Condition{
					{
						Field:    "size",
						Operator: meta.LTE,
						Value:    param.SizeRange[1],
					},
				},
			}
			node = node.Next
		}
	}

	return node, opt
}
func GameWhereNodeIssueDateRange(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption) {
	if len(param.IssueDateRange) > 0 {
		node.Next = &meta.WhereNode{
			Conditions: []*meta.Condition{
				{
					Field:    "issue_date",
					Operator: meta.GTE,
					Value:    param.IssueDateRange[0],
				},
			},
		}

		node = node.Next
		if len(param.IssueDateRange) > 1 {
			node.Next = &meta.WhereNode{
				Conditions: []*meta.Condition{
					{
						Field:    "issue_date",
						Operator: meta.LTE,
						Value:    param.IssueDateRange[1],
					},
				},
			}
			node = node.Next
		}
	}

	return node, opt
}
func GameWhereNodeCreatedAtRange(ctx context.Context, param handler.GameListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption) {
	if len(param.CreatedAtRange) > 0 {
		node.Next = &meta.WhereNode{
			Conditions: []*meta.Condition{
				{
					Field:    "created_at",
					Operator: meta.GTE,
					Value:    param.CreatedAtRange[0],
				},
			},
		}

		node = node.Next
		if len(param.CreatedAtRange) > 1 {
			node.Next = &meta.WhereNode{
				Conditions: []*meta.Condition{
					{
						Field:    "created_at",
						Operator: meta.LTE,
						Value:    param.CreatedAtRange[1],
					},
				},
			}
			node = node.Next
		}
	}

	return node, opt
}
