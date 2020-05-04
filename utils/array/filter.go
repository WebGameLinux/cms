package array

import (
		"regexp"
		"strings"
)

// 含空字符正则
func GetNilCharRegexp() *regexp.Regexp {
		return regexp.MustCompile(`( |\n|\r|\t|\r\n|\t\n|\r\t| )+`)
}

// 空字符串正则
func GetNilStringRegexp() *regexp.Regexp {
		return regexp.MustCompile(`^( |\n|\r|\t|\r\n|\t\n|\r\t|\s| )+$`)
}

// 字符串数组过滤
func Filter(array []string, handlers ...interface{}) []string {
		if len(handlers) == 0 {
				handlers = append(handlers, StrArrayRegexpFilter)
		}
		var (
				strArr      []string
				regexps     []*regexp.Regexp
				boolFuncs   []func(string) bool
				compares    []func(int, string, []string) bool
				regexpFuncs []func([]string, ...*regexp.Regexp) []string
		)
		// 处理器分类
		for _, handler := range handlers {
				if fn, ok := handler.(func([]string, ...*regexp.Regexp) []string); ok {
						regexpFuncs = append(regexpFuncs, fn)
				}
				if fn, ok := handler.(func(string) bool); ok {
						boolFuncs = append(boolFuncs, fn)
				}
				if reg, ok := handler.(*regexp.Regexp); ok {
						regexps = append(regexps, reg)
				}
				if fn, ok := handler.(func(int, string, []string) bool); ok {
						compares = append(compares, fn)
				}
				if str, ok := handler.(string); ok {
						strArr = append(strArr, str)
				}
		}
		// 单个过滤
		for i := 0; i < len(array); {
				str := array[i]
				if str == "" {
						array = RemoveStringArrayItem(i, array)
						continue
				}
				// 字符串对比
				for _, s := range strArr {
						if s == str {
								array = RemoveStringArrayItem(i, array)
								continue
						}
				}
				// 正则对比
				for _, reg := range regexps {
						if reg.MatchString(str) {
								array = RemoveStringArrayItem(i, array)
								continue
						}
				}
				// 自定义对比函数
				for _, boolFunc := range boolFuncs {
						if boolFunc(str) {
								array = RemoveStringArrayItem(i, array)
								continue
						}
				}
				// 自定义全索引对比
				for _, compare := range compares {
						if compare(i, str, array) {
								array = RemoveStringArrayItem(i, array)
								continue
						}
				}
				i++
		}
		// 整体过滤
		if len(regexpFuncs) > 0 {
				for _, handler := range regexpFuncs {
						array = handler(array)
				}
		}
		return array
}

// 是否包含某个字符串
// @param array []string : 查找源
// @param compare string : 要查找的字符串
// @param ignore bool    : 是否忽略大小 (可选,默认false不忽略)
func Contains(array []string, compare string, ignore ...bool) bool {
		if len(Index(array, compare, 1, ignore...)) > 0 {
				return true
		}
		return false
}

// 查找字符串出现下标
// @param array []string : 查找源
// @param compare string : 要查找的字符串
// @param limit int      : 要查找的字符串索引数量限制(-1,表示查找出所有)
// @param ignore bool    : 是否忽略大小 (可选,默认false不忽略)
func Index(array []string, compare string, limit int, ignore ...bool) []int {
		var ret []int
		if len(ignore) == 0 {
				ignore = append(ignore, false)
		}
		if len(array) == 0 {
				return ret
		}
		for i, v := range array {
				if !ignore[0] && v != compare {
						continue
				}
				if ignore[0] && !strings.EqualFold(v, compare) {
						continue
				}
				ret = append(ret, i)
				if limit != -1 && len(ret) >= limit {
						break
				}
		}
		return ret
}

// 查找字符串出现次数
// @param array []string : 查找源
// @param compare string : 要查找的字符串
// @param ignore bool    : 是否忽略大小 (可选,默认false不忽略)
func Count(array []string, compare string, ignore ...bool) int {
		return len(Index(array, compare, -1, ignore...))
}

// 移除
func RemoveStringArrayItem(index int, array []string) []string {
		length := len(array)
		if length == 0 {
				return array
		}
		if index+1 >= length {
				array = array[:index]
		}
		if index+1 < length && index >= 0 {
				return append(array[:index], array[index+1:]...)
		}
		return array
}

// 过滤
func StrArrayRegexpFilter(arr []string, reg ...*regexp.Regexp) []string {
		if len(arr) == 0 {
				return arr
		}
		if len(reg) == 0 {
				reg = append(reg, GetNilStringRegexp())
		}
		for _, match := range reg {
				for i := 0; i < len(arr); {
						str := arr[i]
						if !match.MatchString(str) {
								i++
								continue
						}
						arr = RemoveStringArrayItem(i, arr)
				}
		}
		return arr
}

// 对比过滤
func StrArrayCompareFilter(arr []string, compare ...func(int, string, []string) bool) []string {

		if len(arr) == 0 {
				return arr
		}
		if len(compare) == 0 {
				compare = append(compare, CompareNilStr)
		}
		for _, match := range compare {
				for i := 0; i < len(arr); {
						str := arr[i]
						if !match(i, str, arr) {
								i++
								continue
						}
						arr = RemoveStringArrayItem(i, arr)
				}
		}
		return arr
}

// 对比空字符串
func CompareNilStr(index int, str string, arr []string) bool {
		if len(arr) == 0 || index < 0 {
				return false
		}
		return IsNilString(str)
}

// 是否空字符串
func IsNilString(str string) bool {
		return GetNilStringRegexp().MatchString(str)
}

// 是否有空字符串
func HasNilChar(str string) bool {
		return GetNilCharRegexp().MatchString(str)
}
