package mapper

import (
		"encoding/json"
		"fmt"
		"github.com/WebGameLinux/cms/utils/reflects"
		"github.com/WebGameLinux/cms/utils/times"
		"reflect"
		"strconv"
		"strings"
		"time"
)

type Filter func(key string, value interface{}) bool

type MapMustGetter interface {
		GetInt(key string, def ...int) int
		GetBool(key string, def ...bool) bool
		Get(key string, def ...string) string
		GetBigInt(key string, def ...int64) int64
		GetTime(key string, def ...time.Time) time.Time
		GetDate(key string, def ...time.Time) time.Time
		GetWeek(key string, def ...time.Weekday) time.Weekday
		GetDuration(key string, def ...time.Duration) time.Duration
		GetFloat(key string, def ...float32) float32
		GetFloatN(key string, def ...float64) float64
}

type MapAggregationGetter interface {
		MapMustGetter
		GetMap(key string, def ...Mapper) Mapper
		GetAny(key string, def ...interface{}) interface{}
		GetType(key string) reflect.Type
		GetValue(key string) reflect.Value
		GetStruct(key string, structPtr interface{}) bool
		Struct(structPtr interface{}) bool
		Filter(filter ...Filter) Mapper
		Keys(filter ...Filter) []string
		Values(filter ...Filter) []interface{}
		Exists(key string) bool
}

func (this Mapper) GetBool(key string, def ...bool) bool {
		if len(def) == 0 {
				def = append(def, false)
		}
		if n, ok := this[key]; ok {
				if num, ok := n.(bool); ok {
						return num
				}
				if num, ok := n.(*bool); ok {
						return *num
				}
				if num, ok := n.(int); ok {
						if num <= 0 {
								return false
						}
						if num == 1 {
								return true
						}
				}
		}
		return def[0]
}

func (this Mapper) Get(key string, def ...string) string {
		if len(def) == 0 {
				def = append(def, "")
		}
		if n, ok := this[key]; ok {
				if num, ok := n.(string); ok {
						return num
				}
				if num, ok := n.(*string); ok {
						return *num
				}
				if strAble, ok := n.(fmt.Stringer); ok {
						return strAble.String()
				}
				if num, ok := n.(int); ok {
						return strconv.Itoa(num)
				}
				if num, ok := n.(int64); ok {
						return fmt.Sprintf("%d", num)
				}
				if num, ok := n.(*int64); ok {
						return fmt.Sprintf("%d", *num)
				}
				if num, ok := n.(float64); ok {
						return fmt.Sprintf("%f", num)
				}
				if num, ok := n.(*float64); ok {
						return fmt.Sprintf("%f", *num)
				}
				if num, ok := n.(float32); ok {
						return fmt.Sprintf("%f", num)
				}
				if num, ok := n.(*float32); ok {
						return fmt.Sprintf("%f", *num)
				}
		}
		return def[0]
}

func (this Mapper) GetBigInt(key string, def ...int64) int64 {
		if len(def) == 0 {
				def = append(def, 0)
		}
		if n, ok := this[key]; ok {
				if num, ok := n.(int64); ok {
						return num
				}
				if num, ok := n.(*int64); ok {
						return *num
				}
				if num, ok := n.(string); ok {
						if w, err := strconv.Atoi(num); err == nil {
								return int64(w)
						}
				}
		}
		return def[0]
}

func (this Mapper) GetTime(key string, def ...time.Time) time.Time {
		if len(def) == 0 {
				def = append(def, time.Now())
		}
		if n, ok := this[key]; ok {
				if num, ok := n.(time.Time); ok {
						return num
				}
				if num, ok := n.(*time.Time); ok {
						return *num
				}
				if num, ok := n.(string); ok {
						if w, ok := times.ParseTime(num); ok {
								return w
						}
				}
				if num, ok := n.(int64); ok {
						if num > 0 {
								return time.Unix(num, 0)
						}
				}
		}
		return def[0]
}

func (this Mapper) GetDate(key string, def ...time.Time) time.Time {
		if len(def) == 0 {
				def = append(def, time.Now())
		}
		if n, ok := this[key]; ok {
				if num, ok := n.(time.Time); ok {
						return num
				}
				if num, ok := n.(*time.Time); ok {
						return *num
				}
				if num, ok := n.(string); ok {
						if w, ok := times.ParseTime(num); ok {
								return w
						}
				}
				if num, ok := n.(int64); ok {
						if num > 0 {
								return time.Unix(num, 0)
						}
				}
		}
		return def[0]
}

func (this Mapper) GetWeek(key string, def ...time.Weekday) time.Weekday {
		if len(def) == 0 {
				def = append(def, 0)
		}
		if n, ok := this[key]; ok {
				if num, ok := n.(time.Weekday); ok {
						return num
				}
				if num, ok := n.(*time.Weekday); ok {
						return *num
				}
				if num, ok := n.(string); ok {
						if w, ok := times.ParseWeekday(num); ok {
								return w
						}
				}
				if num, ok := n.(int); ok {
						if num >= 0 && num < 6 {
								return time.Weekday(num)
						}
				}
		}
		return def[0]
}

func (this Mapper) GetDuration(key string, def ...time.Duration) time.Duration {
		if len(def) == 0 {
				def = append(def, 0)
		}
		if n, ok := this[key]; ok {
				if num, ok := n.(time.Duration); ok {
						return num
				}
				if num, ok := n.(*time.Duration); ok {
						return *num
				}
				if num, ok := n.(string); ok {
						if d, err := time.ParseDuration(num); err == nil {
								return d
						}
				}
				if num, ok := n.(int64); ok {
						return time.Duration(num)
				}
				if num, ok := n.(float64); ok {
						return time.Duration(int64(num))
				}
		}
		return def[0]
}

func (this Mapper) GetFloat(key string, def ...float32) float32 {
		if len(def) == 0 {
				def = append(def, 0)
		}
		if n, ok := this[key]; ok {
				if num, ok := n.(float32); ok {
						return num
				}
				if num, ok := n.(*float32); ok {
						return *num
				}
				if num, ok := n.(float64); ok {
						return float32(num)
				}
				if num, ok := n.(*float64); ok {
						return float32(*num)
				}
				if num, ok := n.(int64); ok {
						return float32(num)
				}
				if num, ok := n.(int); ok {
						return float32(num)
				}
		}
		return def[0]
}

func (this Mapper) GetFloatN(key string, def ...float64) float64 {
		if len(def) == 0 {
				def = append(def, 0)
		}
		if n, ok := this[key]; ok {
				if num, ok := n.(float64); ok {
						return num
				}
				if num, ok := n.(*float64); ok {
						return *num
				}
				if num, ok := n.(int64); ok {
						return float64(num)
				}
				if num, ok := n.(int); ok {
						return float64(num)
				}
		}
		return def[0]
}

func (this Mapper) GetInt(key string, def ...int) int {
		if len(def) == 0 {
				def = append(def, 0)
		}
		if n, ok := this[key]; ok {
				if num, ok := n.(int); ok {
						return num
				}
				if num, ok := n.(*int); ok {
						return *num
				}
				if num, ok := n.(int64); ok {
						return int(num)
				}
				if num, ok := n.(string); ok {
						if w, err := strconv.Atoi(num); err == nil {
								return w
						}
				}
		}
		return def[0]
}

func (this Mapper) GetMap(key string, def ...Mapper) Mapper {
		if len(def) == 0 {
				def = append(def, Mapper{})
		}
		if n, ok := this[key]; ok {
				if m, ok := n.(Mapper); ok {
						return m
				}
				if m, ok := n.(*Mapper); ok {
						return *m
				}
		}
		return def[0]
}

func (this Mapper) GetAny(key string, def ...interface{}) interface{} {
		if len(def) == 0 {
				def = append(def, nil)
		}
		if n, ok := this[key]; ok {
				return n
		}
		return def[0]
}

func (this Mapper) GetType(key string) reflect.Type {
		if n, ok := this[key]; ok {
				return reflect.TypeOf(n)
		}
		return nil
}

func (this Mapper) GetValue(key string) reflect.Value {
		if n, ok := this[key]; ok {
				return reflect.ValueOf(n)
		}
		return reflect.ValueOf(nil)
}

func (this Mapper) GetStruct(key string, structPtr interface{}) bool {
		if n, ok := this[key]; ok {
				if b, err := json.Marshal(n); err == nil {
						if json.Unmarshal(b, structPtr) == nil {
								return true
						}
				}
		}
		return false
}

func (this Mapper) Struct(structPtr interface{}) bool {
		if b, err := json.Marshal(this); err == nil {
				if json.Unmarshal(b, structPtr) == nil {
						return true
				}
		}
		return false
}

func (this Mapper) Filter(filter ...Filter) Mapper {
		if len(filter) == 0 {
				filter = append(filter, filterNil)
		}
		var m = make(Mapper)
		for k, v := range this {
				flag := true
				for _, fn := range filter {
						if !fn(k, v) {
								flag = false
						}
				}
				if flag {
						m[k] = v
				}
		}
		return m
}

func (this Mapper) Keys(filter ...Filter) []string {
		if len(filter) == 0 {
				filter = append(filter, filterNil)
		}
		var m []string
		for k, v := range this {
				flag := true
				for _, fn := range filter {
						if !fn(k, v) {
								flag = false
						}
				}
				if flag {
						m = append(m, k)
				}
		}
		return m
}

func (this Mapper) Values(filter ...Filter) []interface{} {
		if len(filter) == 0 {
				filter = append(filter, filterNil)
		}
		var m []interface{}
		for k, v := range this {
				flag := true
				for _, fn := range filter {
						if !fn(k, v) {
								flag = false
						}
				}
				if flag {
						m = append(m, v)
				}
		}
		return m
}

func (this Mapper) Exists(key string) bool {
		if _, ok := this[key]; ok {
				return true
		}
		return false
}

type MapMustSetter interface {
		Set(string, interface{})
		SetEx(string, interface{})
}

type MapAggregationSetter interface {
		MapMustSetter
		Sets(key string, value interface{}, dot ...string)
		Delete(key string)
}

func (this Mapper) Set(key string, value interface{}) {
		this[key] = value
}

func (this Mapper) SetEx(key string, value interface{}) {
		if _, ok := this[key]; ok {
				return
		}
		this[key] = value
}

func (this Mapper) Delete(key string) {
		delete(this, key)
}

func (this Mapper) Sets(key string, value interface{}, dot ...string) {
		if len(dot) == 0 {
				dot = append(dot, ".")
		}
		if !strings.Contains(key, dot[0]) {
				this.Set(key, value)
				return
		}
		// @todo
}

func (this Mapper) String() string {
		if buf, err := json.Marshal(this); err == nil {
				return string(buf)
		}
		return fmt.Sprintf("%+v", map[string]interface{}(this))
}

// 与 前者 map 对比不同 结果 (不比较各自的特有的键,比较共同有键)
func (this Mapper) Diff(m Mapper, compares ...CompareHandler) Mapper {
		var (
				k  string
				ok bool
				v  interface{}
				v2 interface{}
		)
		if len(compares) == 0 {
				compares = append(compares, Compare)
		}
		var resMap = make(Mapper)
		for k, v = range this {
				if v2, ok = m[k]; !ok {
						continue
				}
				if v2 == v {
						continue
				}
				for _, compare := range compares {
						switch compare(v, v2) {
						case reflects.UnCompressLess:
								// 无法比较
								continue
						case reflects.CompressEq:
								// 相等
								continue
						case reflects.CompressNotEq:
								// 不相等
								resMap[k] = v2
						case reflects.CompressLg:
								// 大于
								resMap[k] = v2
						case reflects.CompressLt:
								// 小于
								resMap[k] = v2
						case reflects.CompressSimilar:
								// 相似
								resMap[k] = []interface{}{v, v2}
						default:
								continue
						}
				}
		}
		return resMap
}

// 数据
func (this Mapper) Len() int {
		return len(map[string]interface{}(this))
}

func (this Mapper) ReName(key string, newKey string, delOld ...bool) Mapper {
		if len(delOld) == 0 {
				delOld = append(delOld, true)
		}
		if v, ok := this[key]; ok && newKey != key && newKey != "" {
				this[newKey] = v
				if delOld[0] {
						this.Delete(key)
				}
		}
		return this
}

// 重命名键
// 如果不存在 则使用默认值填充
// @param key string : 旧键名
// @param newKey string :  新键名
// @param def interface : 默认值
// @param delOld bool   : 重名后是否删除旧键,默认删除 [可选]
func (this Mapper) ReNameN(key string, newKey string, def interface{}, delOld ...bool) Mapper {
		if len(delOld) == 0 {
				delOld = append(delOld, true)
		}
		if v, ok := this[key]; ok && newKey != key && newKey != "" {
				this[newKey] = v
				if delOld[0] {
						this.Delete(key)
				}
		} else {
				this[newKey] = def
		}
		return this
}

func (this Mapper) Transform(key string, handler func(string, interface{}) (string, interface{}), delOld ...bool) Mapper {
		if len(delOld) == 0 {
				delOld = append(delOld, false)
		}
		if v, ok := this[key]; ok && handler != nil {
				k, newValue := handler(key, v)
				if k == "" {
						this.Delete(key)
				} else {
						this[k] = newValue
						if k != key && delOld[0] {
								this.Delete(key)
						}
				}
		}
		return this
}

func (this Mapper) Transforms(handler func(string, interface{}) (string, interface{}), delOld ...bool) Mapper {
		if handler == nil {
				return this
		}
		if len(delOld) == 0 {
				delOld = append(delOld, false)
		}
		for key, value := range this {
				k, v := handler(key, value)
				if k == "" {
						this.Delete(key)
						continue
				}
				if k != key {
						this[k] = v
						if delOld[0] {
								this.Delete(key)
						}
				}
		}
		return this
}

// 通过map 设置结构体的值
func SetByMap(v interface{}, m Mapper) bool {
		return reflects.SetStructByMap(v, m)
}

func filterNil(key string, value interface{}) bool {
		if key == "" {
				return false
		}
		if value == nil {
				return false
		}
		return true
}
