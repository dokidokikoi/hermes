package game

import (
	"hermes/db/data"
	"hermes/internal/handler"
	"hermes/model"
	"strconv"

	"github.com/dokidokikoi/go-common/core"
	"github.com/dokidokikoi/go-common/errors"
	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/gin-gonic/gin"
)

func (h Handler) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrValidation, nil)
		return
	}

	g, err := data.GetDataFactory().Game().Get(ctx, &model.Game{ID: uint(id)}, nil)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}

	// character
	var cVos []handler.CharacterVo
	gcs, err := data.GetDataFactory().GameCharacter().List(ctx, &model.GameCharacter{GameID: g.ID}, nil)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	crMap := map[uint]model.CharacterRelation{}
	cIDs := []uint{}
	for _, gc := range gcs {
		crMap[gc.CharacterID] = gc.Relation
		cIDs = append(cIDs, gc.CharacterID)
	}
	node := &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "id",
				Operator: meta.IN,
				Value:    cIDs,
			},
		},
	}
	cs, err := data.GetDataFactory().Character().ListComplex(ctx, &model.Character{}, node, nil)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	for _, c := range cs {
		cVos = append(cVos, handler.CharacterVo{
			ID:      c.ID,
			Name:    c.Name,
			Alias:   c.Alias,
			Gender:  c.Gender.String(),
			Rlation: crMap[c.ID].String(),
			Summary: c.Summary,
			Cover:   c.Cover,
			Images:  c.Images,
			Tags:    c.Tags,
			CV: handler.StaffVo{
				ID:        c.CV.ID,
				Name:      c.CV.Name,
				Cover:     c.CV.Cover,
				Images:    c.CV.Images,
				Alias:     c.CV.Alias,
				CreatedAt: c.CV.CreatedAt,
				Tags:      c.CV.Tags,
				Gender:    c.CV.Gender.String(),
				Summary:   c.CV.Summary,
			},
			CreatedAt: c.CreatedAt,
		})
	}

	// staff
	var sVos []handler.StaffVo
	gss, err := data.GetDataFactory().GameStaff().List(ctx, &model.GameStaff{GameID: g.ID}, nil)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	prMap := map[uint][]model.PersonRelation{}
	pIDs := []uint{}
	for _, gs := range gss {
		prs, ok := prMap[gs.PersonID]
		if ok {
			prs = append(prs, gs.Relation)
			prMap[gs.PersonID] = prs
		} else {
			cIDs = append(cIDs, gs.PersonID)
		}
	}
	node = &meta.WhereNode{
		Conditions: []*meta.Condition{
			{
				Field:    "id",
				Operator: meta.IN,
				Value:    pIDs,
			},
		},
	}
	ss, err := data.GetDataFactory().Person().ListComplex(ctx, &model.Person{}, node, nil)
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	for _, s := range ss {
		var prs []string
		for _, pr := range prMap[s.ID] {
			prs = append(prs, pr.String())
		}
		sVos = append(sVos, handler.StaffVo{
			ID:        s.ID,
			Name:      s.Name,
			Alias:     s.Alias,
			Cover:     s.Cover,
			Images:    s.Images,
			Tags:      s.Tags,
			Summary:   s.Summary,
			Gender:    s.Gender.String(),
			Relation:  prs,
			CreatedAt: s.CreatedAt,
		})
	}

	// veriosn
	gins, err := data.GetDataFactory().GameInstance().List(ctx, &model.GameInstance{GameID: g.ID}, &meta.ListOption{GetOption: meta.GetOption{Select: []string{"ID", "GameID", "Version"}}})
	if err != nil {
		core.WriteResponse(ctx, errors.ApiErrSystemErr, nil)
		return
	}
	var version []string
	for _, v := range gins {
		version = append(version, v.Version)
	}
	gvo := handler.GameVo{
		ID:         g.ID,
		Name:       g.Name,
		Alias:      g.Alias,
		Cover:      g.Cover,
		Images:     g.Images,
		Versions:   version,
		Category:   g.Category,
		Series:     g.Series,
		Developer:  g.Developer,
		Publisher:  g.Publisher,
		Price:      g.Price,
		IssueDate:  g.IssueDate,
		Story:      g.Story,
		Platform:   g.Platform,
		Tags:       g.Tags,
		Characters: cVos,
		Staff:      sVos,
		Links:      g.Links,
		OtherInfo:  g.OtherInfo,
		CreatedAt:  g.CreatedAt,
	}

	core.WriteResponse(ctx, nil, gvo)
}
