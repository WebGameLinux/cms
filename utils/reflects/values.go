package reflects

import (
		"fmt"
		"reflect"
)

// 通过 map 设置 struct 属性值
func SetStructByMap(v interface{}, data map[string]interface{}, filter ...TagFilter) bool {
		if v == nil {
				return false
		}
		if len(filter) == 0 {
				filter = append(filter, JsonTagFilter)
		}
		ty := reflect.TypeOf(v)
		ty = RealType(ty)
		if ty.Kind() != reflect.Struct {
				return false
		}
		var (
				i          = 0
				ok         bool
				key        string
				it         reflect.Value
				value      interface{}
				field      reflect.StructField
				m          map[string]interface{}
				setN       = 0
				val        = RealValue(reflect.ValueOf(v))
				getTagName = filter[0]
				num        = ty.NumField()
		)
		for i = 0; i < num; i++ {
				it = val.Field(i)
				field = ty.Field(i)
				key, ok = getTagName(field.Tag)
				if !ok {
						if field.Anonymous {
								k, num := getAnonymousTagFirst(field.Type, getTagName, field.Name)
								if k == "" || num == 0 {
										key = field.Name
								}
								if k != "" && num > 1 {
										if value, ok = data[key]; !ok {
												continue
										}
										if m, ok = Map2Mapper(value); ok {
												if SetStructByMap(it.Interface(), m, getTagName) {
														setN++
												}
										}
								}
								if k != "" && num == 0 {
										key = k
								}
						}
				}
				if key == "" {
						key = field.Name
				}
				if value, ok = data[key]; ok {
						it.Set(reflect.ValueOf(value))
						setN++
				}
		}
		return setN > 0
}

// 任意类型map转 mapper
func Map2Mapper(v interface{}) (map[string]interface{}, bool) {
		if v == nil {
				return nil, false
		}
		if m, ok := v.(map[string]interface{}); ok {
				return m, true
		}
		if m, ok := v.(*map[string]interface{}); ok {
				return *m, true
		}
		var ty = RealType(reflect.TypeOf(v))
		if ty.Kind() != reflect.Map {
				return nil, false
		}
		var (
				data = make(map[string]interface{})
				Iter = RealValue(reflect.ValueOf(v)).MapRange()
		)
		for Iter.Next() {
				k := Iter.Key()
				value := Iter.Value()
				key := k.Interface()
				name := Any2Str(key)
				if name == "" {
						continue
				}
				data[name] = value.Interface()
		}
		return data, len(data) > 0
}

// 任意类型输出字符串
func Any2Str(v interface{}) string {
		if v == nil {
				return ""
		}
		if str, ok := v.(fmt.Stringer); ok {
				return str.String()
		}
		var ty = RealType(reflect.TypeOf(v))
		switch ty.Kind() {
		case reflect.String:
				return fmt.Sprintf("%s", v)
		case reflect.Int:
				fallthrough
		case reflect.Int32:
				fallthrough
		case reflect.Int64:
				fallthrough
		case reflect.Int16:
				fallthrough
		case reflect.Int8:
				return fmt.Sprintf("%d", v)
		case reflect.Interface:
				return fmt.Sprintf("%+v", v)
		case reflect.Ptr:
				return fmt.Sprintf("<Ptr>%x", v)
		case reflect.Func:
				value := RealValue(reflect.ValueOf(v))
				return fmt.Sprintf("&%x.<Func:<%s>>", value.Interface(), value.Kind().String())
		case reflect.Array:
				value := RealValue(reflect.ValueOf(v))
				return fmt.Sprintf("&%x.<Array:<%s:%d>>%+v", value.Interface(), value.Kind().String(), value.Len(), value.Interface())
		case reflect.Slice:
				value := RealValue(reflect.ValueOf(v))
				return fmt.Sprintf("&%x.<Slice:<%s:%d>>%+v", value.Interface(), value.Kind().String(), value.Len(), value.Interface())
		case reflect.Invalid:
				return "<Invalid>"
		case reflect.Uint64:
				fallthrough
		case reflect.Uint8:
				fallthrough
		case reflect.Uint32:
				fallthrough
		case reflect.Uint:
				return fmt.Sprintf("%d", v)
		case reflect.Float64:
				return fmt.Sprintf("%e", v)
		case reflect.Float32:
				return fmt.Sprintf("%f", v)
		case reflect.Map:
				value := RealValue(reflect.ValueOf(v))
				return fmt.Sprintf("<Map:%s>%+v", value.Kind().String(), v)
		case reflect.Struct:
				return fmt.Sprintf("<Struct:%s>%+v", ClassName(v), v)
		case reflect.Chan:
				return fmt.Sprintf("<Chan>%v", v)
		case reflect.Bool:
				return fmt.Sprintf("%t", v)
		}
		return ""
}

// 获取匿名字段第一个属性的tag
func getAnonymousTagFirst(ty reflect.Type, filter TagFilter, def ...string) (string, int) {
		if len(def) == 0 {
				def = append(def, ty.Name())
		}
		if ty.Kind() != reflect.Struct {
				return def[0], 0
		}
		var num = ty.NumField()
		for i := 0; i < num; {
				itTy := ty.Field(i)
				if itTy.Anonymous {
						ky, num := getAnonymousTagFirst(ty, filter, def...)
						return itTy.Name + "." + ky, num
				}
				if name, ok := filter(itTy.Tag); ok {
						return name, num
				}
				return def[0], num
		}
		return def[0], num
}
