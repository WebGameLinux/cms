package types

import "github.com/astaxie/beego/orm"

type SearchParams struct {
		Query   orm.QuerySeter
		Fields  []string
		Effects map[string]interface{}
}

func NewSearchParams() *SearchParams {
		params := new(SearchParams)
		params.Effects = make(map[string]interface{})
		return params
}
