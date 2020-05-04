package conditions

import (
		"github.com/WebGameLinux/cms/utils/mapper"
		"github.com/WebGameLinux/cms/utils/types"
)

type BasePageOptions types.PageOptionDto

// 创建分页基础参数
func MapCreatePageOptions(data map[string]interface{}) *BasePageOptions {
		var p = new(BasePageOptions)
		m := mapper.Mapper(data)
		p.Page = m.GetInt("page", 1)
		p.Count = m.GetInt("count", 10)
		return p
}
