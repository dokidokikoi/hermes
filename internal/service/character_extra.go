package service

import (
	"context"
	"fmt"
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"

	comm_tools "github.com/dokidokikoi/go-common/tools"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

var CharacterBasicSearchNode = []CharacterWhereNodeFunc{
	CharacterWhereNodeCV,
	CharacterWhereNodeCreatedAtRange,
	CharacterWhereNodeGender,
	CharacterWhereNodeKeyword,
	CharacterWhereNodeTag,
}

func CharacterWhereNodeKeyword(ctx context.Context, param handler.CharacterListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	defer func() {
		n = node.Next
		o = opt
	}()

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
					Field:    "alias",
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
	}
	return
}
func CharacterWhereNodeTag(ctx context.Context, param handler.CharacterListReq, node *meta.WhereNode, opt *meta.ListOption) (n *meta.WhereNode, o *meta.ListOption) {
	if len(param.Tags) < 1 {
		return node, opt
	}

	defer func() {
		o = opt
	}()
	db := data.GetDataFactory().CharacterTag().ListComplexDB(ctx, &model.CharacterTag{}, &meta.WhereNode{
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
		Table:           model.Character{}.TableName(),
		InnerQuery:      db,
		InnerQueryAlias: "character_tag",
		TableField:      "id",
		JoinTableField:  "character_id",
	})
	return
}
func CharacterWhereNodeCreatedAtRange(ctx context.Context, param handler.CharacterListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption) {
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
func CharacterWhereNodeGender(ctx context.Context, param handler.CharacterListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption) {
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
func CharacterWhereNodeCV(ctx context.Context, param handler.CharacterListReq, node *meta.WhereNode, opt *meta.ListOption) (*meta.WhereNode, *meta.ListOption) {
	if param.CV == 0 {
		return node, opt
	}

	node.Next = &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "person_id",
				Operator: meta.EQUAL,
				Value:    param.CV,
			},
		},
	}

	return node.Next, opt
}
