package mapper

import "github.com/WebGameLinux/cms/utils/reflects"

type CompareHandler func(source, dist interface{}) int

// 任意两变量值比较
func Compare(source, dist interface{}) int {
		return reflects.DiffValue(source, dist)
}
