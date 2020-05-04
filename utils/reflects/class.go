package reflects

import (
		"reflect"
		"strings"
)

func ClassName(v interface{}) string {
		if v == nil {
				return "<nil>"
		}
		var t = reflect.TypeOf(v)
		t = RealType(t)
		className := t.PkgPath() + "::" + t.Name()
		if className == "::" {
				return t.String()
		}
		return className
}

func Name(v interface{}) string {
		if v == nil {
				return "<nil>"
		}
		var t = reflect.TypeOf(v)
		t = RealType(t)
		return t.Name()
}

func IsBuiltInType(v interface{}) bool {
		var t = ClassName(v)
		if strings.Contains(t, "/") {
				return false
		}
		return true
}

func IsPointer(t reflect.Type) bool {
		// strings.Index(t.String(), "*") == 0
		return t.Kind() == reflect.Ptr
}

func RealType(t reflect.Type) reflect.Type {
		for t.Kind() != reflect.Func && IsPointer(t) {
				t = t.Elem()
		}
		return t
}

func RealValue(t reflect.Value) reflect.Value {
		for t.Kind() != reflect.Func && IsPointerValue(t) {
				t = t.Elem()
		}
		return t
}

func IsPointerValue(v reflect.Value) bool {
		return v.Kind() == reflect.Ptr
}

func PointerN(t reflect.Type) int {
		var i = 0
		for IsPointer(t) {
				t = t.Elem()
				i++
		}
		return i
}

func PointerNum(v interface{}) int {
		return PointerN(reflect.TypeOf(v))
}
