package bangumi

type SearchParam struct {
	Type          int    `query:"type"`
	ResponseGroup string `query:"responseGroup"`
	Start         int    `query:"start"`
	MaxResults    int    `query:"max_results"`
}
