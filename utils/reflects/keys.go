package reflects

import (
		"reflect"
		"strings"
)

type TagFilter func(tag reflect.StructTag) (string, bool)

func GetItemsTypes(o interface{}, filter ...TagFilter) map[string]string {
		var items = make(map[string]string)
		if o == nil {
				return items
		}
		if len(filter) == 0 {
				filter = append(filter, JsonTagFilter)
		}
		ty := reflect.TypeOf(o)
		ty = RealType(ty)
		if ty.Kind() != reflect.Struct {
				return items
		}
		handler := filter[0]
		val := RealValue(reflect.ValueOf(o))

		for i := 0; i < ty.NumField(); i++ {
				it := val.Field(i)
				itTy := ty.Field(i)
				if itTy.Anonymous {
						continue
				}
				name, ok := handler(itTy.Tag)
				kd := it.Kind()
				kdStr := kd.String()
				if kd == reflect.Struct {
						kdStr = kdStr + "." + ClassName(it.Interface())
				}
				if ok {
						items[name] = kdStr
				} else {
						items[itTy.Name] = kdStr
				}
		}
		return items
}

func GetAllItemsTypes(ty reflect.Type, value reflect.Value, container *map[string]string, filter ...TagFilter) map[string]string {
		if container == nil {
				m := make(map[string]string)
				container = &m
		}
		if len(filter) == 0 {
				filter = append(filter, JsonTagFilter)
		}
		ty = RealType(ty)
		if ty.Kind() != reflect.Struct {
				return *container
		}
		handler := filter[0]
		value = RealValue(value)
		for i := 0; i < ty.NumField(); i++ {
				it := value.Field(i)
				itTy := ty.Field(i)
				if itTy.Anonymous {
						GetAllItemsTypes(itTy.Type, it, container, filter...)
						continue
				}
				name, ok := handler(itTy.Tag)
				kd := it.Kind()
				kdStr := kd.String()
				if kd == reflect.Struct {
						kdStr = kdStr + "." + ClassName(it.Interface())
				}
				if ok {
						(*container)[name] = kdStr
				} else {
						(*container)[itTy.Name] = kdStr
				}
		}
		return *container
}

func GetItemsAllTypes(o interface{}, filter ...TagFilter) map[string]string {
		var items = make(map[string]string)
		if o == nil {
				return items
		}
		if len(filter) == 0 {
				filter = append(filter, JsonTagFilter)
		}
		ty := reflect.TypeOf(o)
		ty = RealType(ty)
		if ty.Kind() != reflect.Struct {
				return items
		}
		val := RealValue(reflect.ValueOf(o))
		return GetAllItemsTypes(ty, val, &items, filter...)
}

func GetItemsAllValues(o interface{}, filter ...TagFilter) map[string]interface{} {
		var items = make(map[string]interface{})
		if o == nil {
				return items
		}
		if len(filter) == 0 {
				filter = append(filter, JsonTagFilter)
		}
		ty := reflect.TypeOf(o)
		ty = RealType(ty)
		if ty.Kind() != reflect.Struct {
				return items
		}
		val := RealValue(reflect.ValueOf(o))
		return GetAllItemsValues(ty, val, &items, filter...)

}

func GetItemsValues(o interface{}, filter ...TagFilter) map[string]interface{} {
		var items = make(map[string]interface{})
		if o == nil {
				return items
		}
		if len(filter) == 0 {
				filter = append(filter, JsonTagFilter)
		}
		ty := reflect.TypeOf(o)
		ty = RealType(ty)
		if ty.Kind() != reflect.Struct {
				return items
		}
		handler := filter[0]
		val := RealValue(reflect.ValueOf(o))
		for i := 0; i < ty.NumField(); i++ {
				it := val.Field(i)
				itTy := ty.Field(i)
				name, ok := handler(itTy.Tag)
				if itTy.Anonymous {
						continue
				}
				if ok {
						items[name] = it.Interface()
				} else {
						items[itTy.Name] = it.Interface()
				}
		}
		return items
}

func GetAllItemsValues(ty reflect.Type, value reflect.Value, container *map[string]interface{}, filter ...TagFilter) map[string]interface{} {
		if container == nil {
				m := make(map[string]interface{})
				container = &m
		}
		if len(filter) == 0 {
				filter = append(filter, JsonTagFilter)
		}
		ty = RealType(ty)
		if ty.Kind() != reflect.Struct {
				return *container
		}
		handler := filter[0]
		value = RealValue(value)
		for i := 0; i < ty.NumField(); i++ {
				it := value.Field(i)
				itTy := ty.Field(i)
				if itTy.Anonymous {
						GetAllItemsValues(itTy.Type, it, container, filter...)
						continue
				}
				name, ok := handler(itTy.Tag)
				if ok {
						(*container)[name] = it.Interface()
				} else {
						(*container)[itTy.Name] = it.Interface()
				}
		}
		return *container
}

func JsonTagFilter(tag reflect.StructTag) (string, bool) {
		tagStr := tag.Get("json")
		if tagStr == "" {
				return "", false
		}
		if strings.Contains(tagStr, ";") {
				tags := strings.SplitN(tagStr, ";", -1)
				return tags[0], true
		}
		return tagStr, true
}

