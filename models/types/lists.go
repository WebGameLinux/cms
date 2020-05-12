package types

import "math"

// 分页数据
type Meta struct {
		Page  int  `json:"page"`
		Count int  `json:"count"`
		Size  int  `json:"size"`
		More  bool `json:"more"`
		Total int  `json:"total"`
		query interface{}
}

// 分页信息数据
func NewMeta() *Meta {
		meta := new(Meta)
		meta.Count = 10
		meta.Page = 1
		meta.More = true
		meta.Size = 10
		meta.query = ""
		return meta
}

func (this *Meta) GetQuerySql() interface{} {
		return this.query
}

func (this *Meta) SetQuery(query interface{}) *Meta {
		this.query = query
		return this
}

func (this *Meta) HasMore() bool {
		return this.Total > 0 && this.More
}

func (this *Meta) Next() *Meta {
		if !this.HasMore() {
				return nil
		}
		if this.Total/this.Count <= this.Page || this.Page >= math.MaxInt64 {
				return nil
		}
		meta := NewMeta()
		meta.More = this.More
		meta.query = this.query
		meta.Size = this.Size
		meta.Total = this.Total

		meta.Page = this.Page + 1
		meta.Count = this.Count
		return meta
}

func (this *Meta) Prev() *Meta {
		if this.Page <= 1 || this.Page <= math.MinInt64 {
				return nil
		}
		meta := NewMeta()
		meta.More = this.More
		meta.query = this.query
		meta.Total = this.Total
		meta.Size = this.Size
		meta.Page = this.Page + 1
		meta.Count = this.Count
		return meta
}

func (this *Meta) Limit() int {
		if this.Count == 0 {
				return 10
		}
		return this.Count
}

func (this *Meta) Offset() int {
		if this.Page <= 1 {
				return 0
		}
		return (this.Page - 1) * this.Limit()
}
