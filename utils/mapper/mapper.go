package mapper

import (
		"github.com/WebGameLinux/cms/utils/iterator"
		"strings"
)

type Mapper map[string]interface{}

type MapSetter interface {
		Set(string, interface{})
		SetStrict(key string, v interface{}, ty string) error
}

type MapGetter interface {
		GetAny(string) interface{}
		Get(string, ...string) string
		GetBool(string, ...bool) bool
		GetInt(string, ...int) int
		GetBigInt(string, ...int64) int64
		GetFloatN(string, ...float64) float64
		GetFloat(string, ...float32) float32
		Len() int
}

type ValueItemOptCode int

const (
		OptCodeDel ValueItemOptCode = -1
		OptCodeGet ValueItemOptCode = 0
		OptCodeSet ValueItemOptCode = 1
)

type IteratorValueItem interface {
		Key() string
		Value() interface{}
		Set(string, interface{})
		OptCode() ValueItemOptCode
		Get() (string, interface{})
}

type Value struct {
		key     string
		value   interface{}
		optCode ValueItemOptCode
}

func (v *Value) Key() string {
		return v.key
}

func (v *Value) Value() interface{} {
		return v.value
}

func (v *Value) Get() (string, interface{}) {
		return v.key, v.value
}

func (v *Value) Set(key string, val interface{}) {
		v.key = key
		v.value = val
}

func (v *Value) OptCode() ValueItemOptCode {
		return v.optCode
}

func NewIteratorValueItem(code ValueItemOptCode) IteratorValueItem {
		var val = new(Value)
		val.optCode = code
		return val
}

type ContainerMapWrapper struct {
		Container   *Mapper
		ScopeParser func(string) []string
		Iterator    func() iterator.IIterator
}

func (container *ContainerMapWrapper) Set(key string, value interface{}) {
		var items []IteratorValueItem
		scopes := container.ScopeParser(key)
		last := scopes[len(scopes)-1]
		it := container.Iterator()
		for it.NextAble() || 0 == len(scopes) {
				v := it.Current()
				if !v.CompareKey(scopes[0]) {
						continue
				}
				val := v.Value()
				if tmp, ok := val.(iterator.IAggregate); ok {
						it = tmp.Iterator()
						scopes = scopes[1:]
				}
				tmp := NewIteratorValueItem(OptCodeSet)

				if !v.CompareKey(last) {
						tmp.Set(v.KeyString(), val)
				} else {
						tmp.Set(v.KeyString(), value)
				}
				items = append(items, tmp)
		}
		index := cap(scopes) - 1
		if index != len(items)-1 {
				return
		}
		for ; index > 0; index-- {
				item := items[index]
				if item == nil {
						break
				}
				if index-1 < 0 {
						(*container.Container)[item.Key()] = item.Value()
						break
				}
				panic("impl set")
		}
}

func (container *ContainerMapWrapper) SetStrict(key string, v interface{}, ty string) error {
		// @todo 自动类型转换
		switch ty {
		case "String":
				container.Set(key, v.(string))
		case "Int":
				container.Set(key, v.(int))
		case "FloatN":
				container.Set(key, v.(float64))
		case "Float":
				container.Set(key, v.(float32))
		case "Bool":
				container.Set(key, v.(bool))
		default:
				container.Set(key, v)
		}
		return nil
}

func (container *ContainerMapWrapper) GetAny(key string) interface{} {
		var (
				last string
				item = NewIteratorValueItem(OptCodeGet)
		)
		scopes := container.ScopeParser(key)
		last = scopes[len(scopes)-1]
		item.Set(last, nil)
		it := container.Iterator()
		for it.NextAble() || 0 == len(scopes) {
				v := it.Current()
				if !v.CompareKey(scopes[0]) {
						continue
				}
				val := v.Value()
				if tmp, ok := val.(iterator.IAggregate); ok {
						it = tmp.Iterator()
						scopes = scopes[1:]
				}
				item.Set(v.KeyString(), val)
		}
		if item.Key() == last {
				return item.Value()
		}
		return nil
}

func (container *ContainerMapWrapper) Get(key string, def ...string) string {
		if len(def) == 0 {
				def = append(def, "")
		}
		v := container.GetAny(key)
		if v != nil {
				if str, ok := v.(string); ok {
						return str
				}
		}
		return def[0]
}

func (container *ContainerMapWrapper) GetBool(key string, def ...bool) bool {
		if len(def) == 0 {
				def = append(def, false)
		}
		v := container.GetAny(key)
		if v != nil {
				if b, ok := v.(bool); ok {
						return b
				}
		}
		return def[0]
}

func (container *ContainerMapWrapper) GetInt(key string, def ...int) int {
		if len(def) == 0 {
				def = append(def, 0)
		}
		v := container.GetAny(key)
		if v != nil {
				if n, ok := v.(int); ok {
						return n
				}
		}
		return def[0]
}

func (container *ContainerMapWrapper) GetBigInt(key string, def ...int64) int64 {
		if len(def) == 0 {
				def = append(def, 0)
		}
		v := container.GetAny(key)
		if v != nil {
				if n, ok := v.(int64); ok {
						return n
				}
				if n, ok := v.(int); ok {
						return int64(n)
				}
		}
		return def[0]
}

func (container *ContainerMapWrapper) GetFloatN(key string, def ...float64) float64 {
		if len(def) == 0 {
				def = append(def, 0)
		}
		v := container.GetAny(key)
		if v != nil {
				if n, ok := v.(float64); ok {
						return n
				}
				if n, ok := v.(float32); ok {
						return float64(n)
				}
		}
		return def[0]
}

func (container *ContainerMapWrapper) GetFloat(key string, def ...float32) float32 {
		if len(def) == 0 {
				def = append(def, 0)
		}
		v := container.GetAny(key)
		if v != nil {
				if n, ok := v.(float32); ok {
						return n
				}
		}
		return def[0]
}

func (container *ContainerMapWrapper) Len() int {
		if container.Container == nil {
				return 0
		}
		return len(*container.Container)
}

func Setter(v *Mapper) MapSetter {
		var getter = new(ContainerMapWrapper)
		initHandlers(getter)
		return getter
}

func Getter(v *Mapper) MapGetter {
		var setter = new(ContainerMapWrapper)
		setter.Container = v
		initHandlers(setter)
		return setter
}

func initHandlers(m *ContainerMapWrapper) *ContainerMapWrapper {
		m.ScopeParser = dotSplit
		m.Iterator = iteratorMapper
		return m
}

func dotSplit(s string) []string {
		if strings.Contains(s, ".") {
				return strings.SplitN(s, ".", -1)
		}
		return []string{s}
}

// 迭代器原理
func iteratorMapper() iterator.IIterator {

		panic("impl to iteratorGen")
}
