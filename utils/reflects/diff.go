package reflects

import (
		"fmt"
		"reflect"
)

const (
		UnCompressLess  = 0  // 类型不同无法比较
		CompressEq      = 1  // 相等
		CompressNotEq   = -1 // 不相等
		CompressLg      = 2  // 大于
		CompressLt      = 3  // 小于
		CompressSimilar = 4  // 类似
)

// 对两值
func DiffValue(source, dist interface{}) int {
		if source == nil && dist != nil {
				return UnCompressLess
		}
		// 获取类型
		sType := RealType(reflect.TypeOf(source))
		dType := RealType(reflect.TypeOf(dist))
		if dType.Kind() != sType.Kind() {
				return UnCompressLess
		}
		// 获取值
		sVal := RealValue(reflect.ValueOf(source))
		dVal := RealValue(reflect.ValueOf(dist))
		sValue := sVal.Interface()
		dValue := dVal.Interface()
		switch sType.Kind() {
		case reflect.Bool:
				n := sValue.(bool)
				n2 := dValue.(bool)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Int:
				n := sValue.(int)
				n2 := dValue.(int)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Float32:
				n := sValue.(float32)
				n2 := dValue.(float32)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Func:
				if sValue == dValue {
						return CompressEq
				}
				if sVal.Kind().String() == dVal.Kind().String() {
						return CompressSimilar
				}
				return CompressNotEq
		case reflect.Float64:
				n := sValue.(float64)
				n2 := dValue.(float64)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		case reflect.String:
				n := sValue.(string)
				n2 := dValue.(string)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Ptr:
				if sValue == dValue {
						return CompressEq
				}
				return CompressNotEq

		case reflect.Struct:
				if sVal.Kind().String() != dVal.Kind().String() {
						return UnCompressLess
				}
				if ClassName(sValue) != ClassName(dValue) {
						return UnCompressLess
				}
				if sVal == dVal {
						return CompressEq
				}
				return diffStructValue(sType, sVal, dType, dVal)
		case reflect.Map:
				if sVal.Kind().String() != dVal.Kind().String() {
						return UnCompressLess
				}
				if sVal.Len() != dVal.Len() {
						return CompressNotEq
				}
				if sVal == dVal {
						return CompressEq
				}
				return diffMapValue(sType, sVal, dType, dVal)
		case reflect.Complex64:
				n := sValue.(complex64)
				n2 := dValue.(complex64)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Complex128:
				n := sValue.(complex128)
				n2 := dValue.(complex128)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Array:
				fallthrough
		case reflect.Slice:
				if sVal.Len() != dVal.Len() {
						return CompressNotEq
				}
				if sVal.Kind().String() != dVal.Kind().String() {
						return CompressNotEq
				}
				if fmt.Sprintf("%+v", sValue) == fmt.Sprintf("%+v", dValue) {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Chan:
				if sVal.Kind().String() != dVal.Kind().String() {
						return UnCompressLess
				}
				if sVal == dVal {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Int8:
				n := sValue.(int8)
				n2 := dValue.(int8)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Int16:
				n := sValue.(int16)
				n2 := dValue.(int16)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Int32:
				n := sValue.(int32)
				n2 := dValue.(int32)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Uint:
				n := sValue.(uint)
				n2 := dValue.(uint)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Uint16:
				n := sValue.(uint16)
				n2 := dValue.(uint16)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Uintptr:
				n := sValue.(uintptr)
				n2 := dValue.(uintptr)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Uint32:
				n := sValue.(uint32)
				n2 := dValue.(uint32)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		case reflect.Uint64:
				n := sValue.(uint64)
				n2 := dValue.(uint64)
				if n == n2 {
						return CompressEq
				}
				return CompressNotEq
		}
		return UnCompressLess
}

// 对比结构体内容
func diffStructValue(sourceType reflect.Type, sourceValue reflect.Value, distType reflect.Type, distValue reflect.Value) int {
		var num = sourceType.NumField()
		if sourceType.Kind() != distType.Kind() {
				return CompressNotEq
		}
		if sourceValue.Kind() != reflect.Struct || distValue.Kind() != reflect.Struct {
				return CompressNotEq
		}
		for i := 0; i < num; i++ {
				v1 := sourceValue.Field(i)
				v2 := distValue.Field(i)
				if sourceType.Field(i).Name != distType.Field(i).Name {
						return CompressNotEq
				}
				if DiffValue(v1, v2) == CompressEq {
						continue
				}
		}
		return CompressEq
}

// 对比map内容
func diffMapValue(sourceType reflect.Type, sourceValue reflect.Value, distType reflect.Type, distValue reflect.Value) int {
		if sourceType.Kind() != distType.Kind() {
				return CompressNotEq
		}
		if len(sourceValue.MapKeys()) != len(distValue.MapKeys()) {
				return CompressNotEq
		}
		var Iter = sourceValue.MapRange()
		for Iter.Next() {
				k := Iter.Key()
				v1 := Iter.Value()
				v2 := distValue.MapIndex(k)
				if DiffValue(v1.Interface(), v2.Interface()) == CompressEq {
						continue
				}
				return CompressNotEq
		}
		return CompressEq
}
