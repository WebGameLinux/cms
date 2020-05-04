package string

import (
		"github.com/WebGameLinux/cms/utils/array"
		"regexp"
		"strings"
)

// 分隔获取一个
func StrSplitFirst(str string, splits ...string) string {
		var arr = StrSplitN(str, splits...)
		if len(arr) > 0 {
				return arr[0]
		}
		return ""
}

// 多重分隔
func StrSplitN(str string, splits ...string) []string {
		var sliceStr []string
		if str == "" {
				return sliceStr
		}
		if len(splits) == 0 {
				splits = append(splits, " ", ",")
		}
		sliceStr = append(sliceStr, str)
		for _, split := range splits {
				for i := 0; i < len(sliceStr); {
						s := sliceStr[i]
						arr := array.Filter(strings.SplitN(s, split, -1), regexp.MustCompile(`^( )+$`))
						n := len(arr)
						if n <= 0 {
								continue
						}
						if n == 1 {
								i++
								continue
						}
						if i == 0 {
								sliceStr = append(arr, sliceStr[i+1:]...)
						}
						if i > 0 {
								end := sliceStr[i+1:]
								sliceStr = append(sliceStr[:i], arr...)
								if len(end) > 0 {
										sliceStr = append(splits, end...)
								}
						}
				}
		}
		return sliceStr
}
