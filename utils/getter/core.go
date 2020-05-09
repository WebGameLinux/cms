package getter

import (
		"errors"
		"github.com/WebGameLinux/cms/utils/array"
		"github.com/WebGameLinux/cms/utils/mapper"
		"github.com/WebGameLinux/cms/utils/reflects"
		"github.com/WebGameLinux/cms/utils/times"
		"log"
		"reflect"
		"strconv"
		"strings"
		"time"
)

const RawTag = "_raw"
const JsonTag = "json"

var (
		boolTrueStr  = []string{"true", "True", "1"}
		boolFalseStr = []string{"false", "False", "0"}
)

type Item struct {
		Excluded bool
		Value    interface{}
		Keys     []map[string]string
}

type StructGetter struct {
		Tags      []string
		Ctx       interface{}
		Container []*Item
}

type FilterHandler func(string, interface{}) bool

func NewGetter(v interface{}, tags ...string) *StructGetter {
		var (
				g         = new(StructGetter)
				container []*Item
		)
		if v == nil {
				return nil
		}
		g.Container = container
		g.Ctx = v
		if len(tags) == 0 {
				tags = append(tags, JsonTag, RawTag)
		}
		g.AddTag(tags...)
		err := g.Resolver()
		if err != nil {
				log.Fatal("解析struct异常: ", err)
				return nil
		}
		return g
}

func resolver(ty reflect.Type, value reflect.Value, tag string) map[string]interface{} {
		var item = make(map[string]interface{})
		ty = reflects.RealType(ty)
		value = reflects.RealValue(value)
		if ty.Kind() != reflect.Struct {
				return item
		}
		size := ty.NumField()
		for i := 0; i < size; {
				t := ty.Field(i)
				v := value.Field(i)
				if tag == RawTag {
						item[t.Name] = v.Interface()
						return item
				}
				if t.Anonymous {
						return resolver(t.Type, v, tag)
				}
				it := resolverTag(t.Tag, tag)
				if it != nil && len(it) != 0 {
						item[it[tag]] = v.Interface()
						return item
				}
				break
		}
		return item
}

func resolverTag(tag reflect.StructTag, key string) map[string]string {
		var item = make(map[string]string)
		k := tag.Get(key)
		if k == "" {
				return item
		}
		if strings.Contains(k, ";") {
				keys := strings.SplitN(k, ";", -1)
				if len(keys) == 0 {
						return item
				}
				item[key] = keys[0]
		} else {
				item[key] = k
		}
		return item
}

func (this *StructGetter) Resolver(tags ...string) error {
		if len(tags) == 0 {
				tags = this.Tags
		}
		if len(tags) == 0 {
				return nil
		}
		if this.Ctx == nil {
				return errors.New("struct context is nil")
		}
		var (
				ty    = reflects.RealType(reflect.TypeOf(this.Ctx))
				value = reflects.RealValue(reflect.ValueOf(this.Ctx))
		)
		if value.Kind() != reflect.Struct {
				return errors.New("type error")
		}
		var (
				tagAppend = make(map[string]bool)
				size      = ty.NumField()
		)
		for i := 0; i < size; i++ {
				t := ty.Field(i)
				v := value.Field(i)
				item := new(Item)
				item.Keys = make([]map[string]string, 0)
				for _, tag := range tags {
						if tag == RawTag {
								item.Keys = append(item.Keys, map[string]string{tag: t.Name})
								tagAppend[tag] = true
								continue
						}
						if t.Anonymous {
								key := resolver(t.Type, v, tag)
								if key == nil || len(key) == 0 {
										continue
								}
								for key, v := range key {
										item.Keys = append(item.Keys, map[string]string{tag: key})
										if item.Value == nil {
												item.Value = v
										}
										break
								}
								tagAppend[tag] = true
								continue
						} else {
								key := resolverTag(t.Tag, tag)
								if key == nil || len(key) == 0 {
										continue
								}
								item.Keys = append(item.Keys, key)
								tagAppend[tag] = true
						}
				}
				if item.Value == nil {
						item.Value = v.Interface()
				}
				if len(item.Keys) == 0 {
						item.Keys = append(item.Keys, map[string]string{RawTag: t.Name})
				}
				this.Container = append(this.Container, item)
		}
		for tag, _ := range tagAppend {
				this.AddTag(tag)
		}
		return nil
}

func (this *StructGetter) GetTags() []string {
		return this.Tags
}

func (this *StructGetter) RemoveTag(name string) *StructGetter {
		var size = len(this.Tags)
		for i, key := range this.Tags {
				if key != name {
						continue
				}
				if i == 0 {
						this.Tags = this.Tags[i+1:]
				}
				if i > 0 {
						if size <= i+1 {
								this.Tags = this.Tags[:i]
						} else {
								this.Tags = append(this.Tags[:i], this.Tags[i+1:]...)
						}
				}
				break
		}
		return this
}

func (this *StructGetter) AddTag(name ...string) *StructGetter {
		if len(name) == 0 {
				return this
		}
		for _, n := range name {
				if this.HasTag(n) {
						continue
				}
				this.Tags = append(this.Tags, n)
		}
		return this
}

func (this *StructGetter) ExcludeKey(key string, tag ...string) *StructGetter {
		if len(tag) == 0 {
				tag = append(tag, RawTag)
		}
		for _, t := range tag {
				for _, item := range this.Container {
						if item == nil {
								continue
						}
						if item.Excluded {
								continue
						}
						for _, it := range item.Keys {
								if k, ok := it[t]; ok && k == key {
										item.Excluded = true
										return this
								}
						}
				}
		}
		return this
}

func (this *StructGetter) FilterKey(key string, tag ...string) *StructGetter {
		if len(tag) == 0 {
				tag = append(tag, RawTag)
		}
		var exclude = false
		if strings.Index(key, "!") == 0 {
				exclude = true
				key = key[1:]
		}
		for _, t := range tag {
				for _, item := range this.Container {
						if item == nil {
								continue
						}
						for _, it := range item.Keys {
								if k, ok := it[t]; ok && k == key {
										if !item.Excluded {
												item.Excluded = exclude
										}
										return this
								}
						}
				}
		}
		return this
}

func (this *StructGetter) FilterKeys(keys ...string) *StructGetter {
		if len(keys) == 0 {
				return this
		}
		for _, key := range keys {
				exclude := false
				if strings.Index(key, "!") == 0 {
						exclude = true
						key = key[1:]
				}
				for _, item := range this.Container {
						if item == nil {
								continue
						}
						for _, it := range item.Keys {
								ok := false
								for _, k := range it {
										if k == key {
												ok = true
												item.Excluded = exclude
												break
										}
								}
								if ok {
										break
								}
						}
				}
		}
		return this
}

func (this *StructGetter) Keys(tags ...string) []string {
		var keys []string
		if len(tags) == 0 {
				tags = append(tags, RawTag)
		}
		for _, t := range tags {
				for _, item := range this.Container {
						if item == nil {
								continue
						}
						if item.Excluded {
								continue
						}
						for _, it := range item.Keys {
								if k, ok := it[t]; ok {
										keys = append(keys, k)
								}
						}
				}
		}
		return keys
}

func (this *StructGetter) Values(tags ...string) []interface{} {
		var values []interface{}
		if len(tags) == 0 {
				tags = append(tags, RawTag)
		}
		for _, t := range tags {
				for _, item := range this.Container {
						if item == nil {
								continue
						}
						if item.Excluded {
								continue
						}
						for _, it := range item.Keys {
								if _, ok := it[t]; ok {
										values = append(values, item.Value)
								}
						}
				}
		}
		return values
}

func (this *StructGetter) Load(key string) (interface{}, bool) {
		for _, item := range this.Container {
				if item == nil {
						continue
				}
				for _, m := range item.Keys {
						if len(m) == 0 {
								continue
						}
						for _, k := range m {
								if k == key {
										return item.Value, true
								}
						}
				}
		}
		return nil, false
}

func (this *StructGetter) Get(key string, def ...string) string {
		if len(def) == 0 {
				def = append(def, "")
		}
		if v, ok := this.Load(key); ok {
				switch v.(type) {
				case string:
						str := v.(string)
						return str
				case *string:
						str := v.(*string)
						return *str
				}
		}
		return def[0]
}

func (this *StructGetter) GetInt(key string, def ...int) int {
		if len(def) == 0 {
				def = append(def, 0)
		}
		if v, ok := this.Load(key); ok {
				if n, ok := v.(int); ok {
						return n
				}
				if n, ok := v.(*int); ok {
						return *n
				}
				var str string
				if s, ok := v.(string); ok {
						str = s
				}
				if s, ok := v.(*string); ok {
						str = *s
				}
				if str != "" {
						if n, err := strconv.Atoi(str); err == nil {
								return n
						}
				}
		}
		return def[0]
}

func (this *StructGetter) GetIntN(key string, def ...int64) int64 {
		if len(def) == 0 {
				def = append(def, 0)
		}
		if v, ok := this.Load(key); ok {
				if n, ok := v.(int64); ok {
						return n
				}
				if n, ok := v.(*int64); ok {
						return *n
				}
				var str string
				if s, ok := v.(string); ok {
						str = s
				}
				if s, ok := v.(*string); ok {
						str = *s
				}
				if str != "" {
						if n, err := strconv.Atoi(str); err == nil {
								return int64(n)
						}
				}
		}
		return def[0]
}

func (this *StructGetter) GetFloat(key string, def ...float32) float32 {
		if len(def) == 0 {
				def = append(def, 0)
		}
		if v, ok := this.Load(key); ok {
				if n, ok := v.(float32); ok {
						return n
				}
				if n, ok := v.(*float32); ok {
						return *n
				}
				if n, ok := v.(float64); ok {
						return float32(n)
				}
				if n, ok := v.(*float64); ok {
						return float32(*n)
				}
		}
		return def[0]
}

func (this *StructGetter) GetTime(key string, def ...*time.Time) *time.Time {
		if len(def) == 0 {
				def = append(def, nil)
		}
		if value, ok := this.Load(key); ok {
				switch value.(type) {
				case time.Time:
						t := value.(time.Time)
						return &t
				case *time.Time:
						return value.(*time.Time)
				case int64:
						t := time.Unix(value.(int64), 0)
						return &t
				case *int64:
						t := time.Unix(*value.(*int64), 0)
						return &t
				case string:
						if t, ok := times.ParseTime(value.(string)); ok {
								return &t
						}
				case *string:
						if t, ok := times.ParseTime(*value.(*string)); ok {
								return &t
						}
				}
		}
		return def[0]
}

func (this *StructGetter) GetDuration(key string, def ...time.Duration) time.Duration {
		if len(def) == 0 {
				def = append(def, time.Duration(0))
		}
		if value, ok := this.Load(key); ok {
				switch value.(type) {
				case time.Duration:
						return value.(time.Duration)
				case *time.Duration:
						return *value.(*time.Duration)
				case int64:
						return time.Duration(value.(int64))
				case *int64:
						return time.Duration(*value.(*int64))
				case string:
						if d, err := time.ParseDuration(value.(string)); err == nil {
								return d
						}
				case *string:
						if d, err := time.ParseDuration(*value.(*string)); err == nil {
								return d
						}
				}
		}
		return def[0]
}

func (this *StructGetter) GetFloatN(key string, def ...float64) float64 {
		if len(def) == 0 {
				def = append(def, 0)
		}
		if value, ok := this.Load(key); ok {
				switch value.(type) {
				case float64:
						return value.(float64)
				case *float64:
						return *value.(*float64)
				}
		}
		return def[0]
}

func (this *StructGetter) GetBool(key string, def ...bool) bool {
		if len(def) == 0 {
				def = append(def, false)
		}
		if value, ok := this.Load(key); ok {
				switch value.(type) {
				case bool:
						return value.(bool)
				case *bool:
						return *value.(*bool)
				case string:
						str := value.(string)
						if array.Contains(boolTrueStr, str, true) {
								return true
						}
						if array.Contains(boolFalseStr, str, true) {
								return false
						}
				case *string:
						str := *value.(*string)
						if array.Contains(boolTrueStr, str, true) {
								return true
						}
						if array.Contains(boolFalseStr, str, true) {
								return false
						}
				}
		}
		return def[0]
}

func (this *StructGetter) GetArray(key string, def ...[]interface{}) []interface{} {
		if len(def) == 0 {
				def = append(def, []interface{}{})
		}
		if value, ok := this.Load(key); ok {
				switch value.(type) {
				case []interface{}:
						return value.([]interface{})
				case *[]string:
						return *value.(*[]interface{})
				}
		}
		strArr := this.GetStringArray(key)
		if len(strArr) != 0 {
				data := make([]interface{}, 0)
				for _, v := range strArr {
						data = append(data, v)
				}
				return data
		}
		m := this.GetIntArray(key)
		if len(m) != 0 {
				data := make([]interface{}, 0)
				for _, v := range strArr {
						data = append(data, v)
				}
				return data
		}
		return def[0]
}

func (this *StructGetter) GetIntArray(key string, def ...[]int) []int {
		if len(def) == 0 {
				def = append(def, []int{})
		}
		if value, ok := this.Load(key); ok {
				switch value.(type) {
				case []int:
						return value.([]int)
				case *[]string:
						return *value.(*[]int)
				}
		}
		return def[0]
}

func (this *StructGetter) GetStringArray(key string, def ...[]string) []string {
		if len(def) == 0 {
				def = append(def, []string{})
		}
		if value, ok := this.Load(key); ok {
				switch value.(type) {
				case []string:
						return value.([]string)
				case *[]string:
						return *value.(*[]string)
				}
		}
		return def[0]
}

func (this *StructGetter) GetMap(key string, def ...map[interface{}]interface{}) map[interface{}]interface{} {
		if len(def) == 0 {
				def = append(def, make(map[interface{}]interface{}))
		}
		if value, ok := this.Load(key); ok {
				switch value.(type) {
				case map[interface{}]interface{}:
						return value.(map[interface{}]interface{})
				case *map[interface{}]interface{}:
						return *value.(*map[interface{}]interface{})
				}
		}
		kv := this.GetKvMap(key)
		if len(kv) != 0 {
				data := make(map[interface{}]interface{})
				for k, v := range kv {
						data[k] = v
				}
				return data
		}
		m := this.GetMapper(key)
		if len(m) != 0 {
				data := make(map[interface{}]interface{})
				for k, v := range m {
						data[k] = v
				}
				return data
		}
		return def[0]
}

func (this *StructGetter) GetMapper(key string, def ...map[string]interface{}) map[string]interface{} {
		if len(def) == 0 {
				def = append(def, make(map[string]interface{}))
		}
		if value, ok := this.Load(key); ok {
				switch value.(type) {
				case map[string]interface{}:
						return value.(map[string]interface{})
				case *map[string]interface{}:
						return *value.(*map[string]interface{})
				}
		}
		return def[0]
}

func (this *StructGetter) GetKvMap(key string, def ...map[string]string) map[string]string {
		if len(def) == 0 {
				def = append(def, make(map[string]string))
		}
		if value, ok := this.Load(key); ok {
				switch value.(type) {
				case map[string]string:
						return value.(map[string]string)
				case *map[string]string:
						return *value.(*map[string]string)
				}
		}
		return def[0]
}

func (this *StructGetter) Filter(fn FilterHandler, tags ...string) *StructGetter {
		var size = len(tags)
		for _, item := range this.Container {
				if item == nil {
						continue
				}
				if item.Excluded {
						continue
				}
				for _, mKey := range item.Keys {
						if len(mKey) == 0 {
								continue
						}
						if size == 0 {
								for _, k := range mKey {
										if !fn(k, item.Value) {
												item.Excluded = true
												break
										}
								}
								continue
						}
						for _, name := range tags {
								if k, ok := mKey[name]; ok {
										if !fn(k, item.Value) {
												item.Excluded = true
												break
										}
								}
						}
				}
		}
		return this
}

func (this *StructGetter) EmptyExclude() *StructGetter {
		return this.Filter(filter, this.Tags...)
}

func (this *StructGetter) Map(tag ...string) mapper.Mapper {
		var m = make(mapper.Mapper)
		if len(tag) == 0 {
				tag = append(tag, this.Tags[0])
		}
		var name = tag[0]

		for _, item := range this.Container {
				if item == nil {
						continue
				}
				if item.Excluded {
						continue
				}
				for _, mKey := range item.Keys {
						if v, ok := mKey[name]; ok {
								m[v] = item.Value
						}
				}
		}
		return m
}

func (this *StructGetter) HasTag(name string) bool {
		for _, tag := range this.Tags {
				if name == tag {
						return true
				}
		}
		return false
}

func filter(key string, value interface{}) bool {
		if key == "" {
				return false
		}
		if value == nil {
				return false
		}
		switch value.(type) {
		case string:
				if value.(string) == "" {
						return false
				}
		case *string:
				if *value.(*string) == "" {
						return false
				}
		case time.Time:
				t := time.Time{}
				v := value.(time.Time)
				if v.Equal(t) {
						return false
				}
		case *time.Time:
				t := time.Time{}
				v := *value.(*time.Time)
				if v.Equal(t) {
						return false
				}
		}
		return true
}
