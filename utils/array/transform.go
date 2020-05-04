package array

import (
		"github.com/WebGameLinux/cms/utils/reflects"
		"reflect"
)

func IsArray(v interface{}) bool {
		var t = reflect.TypeOf(v)
		t = reflects.RealType(t)
		k := t.Kind()
		return k == reflect.Array || k == reflect.Slice
}

func JoinArrays(arr []interface{}, arr2 interface{}, start ...int) []interface{} {
		if !IsArray(arr2) {
				return arr
		}
		if len(start) == 0 {
				start = append(start, -1)
		}
		var (
				capSize    = cap(arr)
				length     = len(arr)
				val        = reflect.ValueOf(arr2)
				startIndex = start[0]
		)

		if val.CanSet() || val.Kind() == reflect.Slice {
				size := val.Len()
				if length < capSize && capSize > 0 && startIndex <= 0 && size < capSize {
						arr = arr[:]
				}
				for i := 0; i < size; i++ {
						if startIndex != -1 && startIndex < capSize && startIndex > 0 {
								arr[startIndex] = val.Index(i).Interface()
								startIndex++
						} else {
								arr = append(arr, val.Index(i).Interface())
						}
				}
		}
		return arr
}
