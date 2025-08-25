package bangumi

type SearchParamFilter struct {
	Type []int `json:"type"`
	Nsfw bool  `json:"nsfw"`
}

type SearchParam struct {
	Keyword string            `json:"keyword"`
	Filter  SearchParamFilter `json:"filter"`

	Limit  int `query:"limit" json:"-"`
	Offset int `query:"offset" json:"-"`
}
