package service

import (
	"context"
	"fmt"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"

	meta "github.com/dokidokikoi/go-common/meta/option"
	comm_tools "github.com/dokidokikoi/go-common/tools"
)

var PersonBasicSearchNode = []PersonWhereNodeFunc{
	PersonWhereNodeKeyword,
	PersonWhereNodeTag,
	PersonWhereNodeCreatedAtRange,
	PersonWhereNodeGender,
}

func PersonWhereNodeKeyword(ctx context.Context, param handler.PersonListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
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
					Field:    "summary",
					Operator: meta.LIKE,
					Value:    fmt.Sprintf("%%%s%%", keyword),
				},
			}...)
		}
		node = node.Next
	}
	return node, opt
}
func PersonWhereNodeTag(ctx context.Context, param handler.PersonListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	if len(param.Tags) < 1 {
		return node, opt
	}

	defer func() {
		o = opt
	}()
	db := data.GetDataFactory().PersonTag().ListComplexDB(ctx, &model.PersonTag{}, &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "tag_id",
				Operator: meta.IN,
				Value:    param.Tags,
			},
		},
	}, nil)
	opt.GetOption.Join = append(opt.GetOption.Join, &meta.Join{
		Method:          meta.INNER_JOIN,
		Table:           model.Person{}.TableName(),
		InnerQuery:      db,
		InnerQueryAlias: "person_tag",
		TableField:      "id",
		JoinTableField:  "person_id",
	})
	return
}
func PersonWhereNodeCreatedAtRange(ctx context.Context, param handler.PersonListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption) {
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
func PersonWhereNodeGender(ctx context.Context, param handler.PersonListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption) {
	if param.Gender == model.UnKnown {
		return node, opt
	}
	node.Next = &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "gender",
				Operator: meta.EQUAL,
				Value:    param.Gender,
			},
		},
	}

	return node.Next, opt
}
