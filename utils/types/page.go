package types

import (
		"encoding/json"
		"fmt"
)

const (
		PageDef      = 1
		PageCountDef = 20
)

type PageArgs interface {
		Page() int
		Count() int
		Get(string) interface{}
		Set(string, interface{})
		Map() map[string]int
		Dto() *PageOptionDto
		fmt.Stringer
}

type PageOptions struct {
		page  int
		count int
}

type PageOptionDto struct {
		Page  int
		Count int
}

// 创建分页参数
// page,count int
// @param page int : 页码
// @param count int : 页长
func NewPageArg(args ...int) PageArgs {
		var page = &PageOptions{}
		if len(args) == 0 {
				args = append(args, PageDef, PageCountDef)
		}
		if len(args) < 2 {
				args = append(args, PageCountDef)
		}
		page.page = args[0]
		page.count = args[1]
		return page
}

func (p *PageOptions) Page() int {
		if p.page == 0 {
				return PageDef
		}
		return p.page
}

func (p *PageOptions) Count() int {
		if p.count == 0 {
				return PageCountDef
		}
		return p.page
}

func (p *PageOptions) Get(key string) interface{} {
		switch key {
		case "page":
				fallthrough
		case "offset":
				return p.Page()
		case "count":
				fallthrough
		case "size":
				return p.Count()
		}
		return nil
}

func (p *PageOptions) Set(key string, v interface{}) {
		switch key {
		case "page":
				fallthrough
		case "offset":
				num := p.Page()
				if p, ok := v.(int); ok && p != 0 {
						num = p
				}
				if n, ok := v.(*int); ok && (*n) != 0 {
						num = *n
				}
				p.page = num
		case "count":
				fallthrough
		case "size":
				num := p.Count()
				if p, ok := v.(int); ok && p != 0 {
						num = p
				}
				if n, ok := v.(*int); ok && (*n) != 0 {
						num = *n
				}
				p.count = num
		}
}

// 输出 Map
func (p *PageOptions) Map() map[string]int {
		var m = make(map[string]int)
		m["page"] = p.Page()
		m["count"] = p.Count()
		return m
}

func (p *PageOptions) String() string {
		if v, err := json.Marshal(p.Map()); err == nil {
				return string(v)
		}
		return fmt.Sprintf(`{"page":%d,"count":%d}`, p.Page(), p.Count())
}

func (p *PageOptions) Dto() *PageOptionDto {
		var dto = new(PageOptionDto)
		dto.Page = p.Page()
		dto.Count = p.Count()
		return dto
}
