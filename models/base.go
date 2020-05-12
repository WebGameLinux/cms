package models

import (
		"errors"
		utils "github.com/WebGameLinux/cms/utils/beego"
		"github.com/WebGameLinux/cms/utils/reflects"
		"github.com/astaxie/beego/orm"
		uuid "github.com/satori/go.uuid"
		validation "gopkg.in/go-playground/validator.v9"
		"reflect"
		"strconv"
		"strings"
		"time"
)

type BaseWrapper struct {
		Options map[string]interface{}
		Model   interface{}
		Fields  map[string]string
		Error   error
}

type Table interface {
		TableName() string
}

type ModelWrapper interface {
		BindModel(interface{})
		GetOrm() orm.Ormer
		GetOrmUsing() string
		Option(key string) (interface{}, bool)
		OptStr(key string, def ...string) string
		OptInt(key string, def ...int) int
		SetOpt(key string, val interface{})
		OptBool(key string, def ...bool) bool
		OptFloat(key string, def ...float64) float64
		ClearOpt(keys ...string)
		GetQuery() (utils.SqlQueryBuilder, error)
		GetFields() map[string]string
		HasField(key string) bool
		GetError(clean ...bool) error
		Table() string
		NewModel(data ...map[string]interface{}) interface{}
		OptMap(key string, def ...map[string]interface{}) map[string]interface{}
}

type AutoLoaderModel interface {
		LoadByMap(data map[string]interface{}) error
}

type VerifyAbleModel interface {
		Valid() (*validation.Validate, error)
}

func IsOrmUsingError(err error) bool {
		return strings.Contains(err.Error(), ormUsingError)
}

func (this *BaseWrapper) BindModel(m interface{}) {
		if m != nil && this.Model == nil {
				if reflect.TypeOf(m).Elem().Kind() == reflect.Struct {
						this.Model = m
				}
		}
}

func (this *BaseWrapper) GetOrm() orm.Ormer {
		var o = orm.NewOrm()
		connection := this.GetOrmUsing()
		err := o.Using(connection)
		// 重置
		if err != nil {
				// 强制 链接
				if this.OptBool("force") {
						ormOkSetter(err, o, connection)
						this.SetOpt("connection", "default")
				}
		}
		return o
}

func ormOkSetter(err error, o orm.Ormer, using string) {
		utils.Onerror(err)
		if IsOrmUsingError(err) && using != "default" {
				err = o.Using("default")
		}
		if err != nil {
				panic(err)
		}
}

func (this *BaseWrapper) GetOrmUsing() string {
		return this.OptStr("connection", "default")
}

func (this *BaseWrapper) GetQuery() (utils.SqlQueryBuilder, error) {
		var (
				err   error
				query orm.QueryBuilder
		)
		if query, err = orm.NewQueryBuilder(utils.GetDatabaseDriver()); err == nil {
				return utils.NewQueryBuilderWrapper(query, this.GetOrm(), this.NewModel()), nil
		}
		utils.Onerror(err)
		return nil, err
}

func (this *BaseWrapper) Option(key string) (interface{}, bool) {
		if this.Options == nil {
				return nil, false
		}
		if v, ok := this.Options[key]; ok {
				return v, true
		}
		return nil, false
}

func (this *BaseWrapper) OptStr(key string, def ...string) string {
		if len(def) == 0 {
				def = append(def, "")
		}
		if v, ok := this.Option(key); ok {
				if str, ok := v.(string); ok {
						return str
				}
		}
		return def[0]
}

func (this *BaseWrapper) OptInt(key string, def ...int) int {
		if len(def) == 0 {
				def = append(def, 0)
		}
		if v, ok := this.Option(key); ok {
				if num, ok := v.(int); ok {
						return num
				}
		}
		return def[0]
}

func (this *BaseWrapper) OptBool(key string, def ...bool) bool {
		if len(def) == 0 {
				def = append(def, false)
		}
		if v, ok := this.Option(key); ok {
				if b, ok := v.(bool); ok {
						return b
				}
		}
		return def[0]
}

func (this *BaseWrapper) OptFloat(key string, def ...float64) float64 {
		if len(def) == 0 {
				def = append(def, 0)
		}
		if v, ok := this.Option(key); ok {
				if b, ok := v.(float64); ok {
						return b
				}
		}
		return def[0]
}

func (this *BaseWrapper) OptMap(key string, def ...map[string]interface{}) map[string]interface{} {
		if len(def) == 0 {
				def = append(def, nil)
		}
		if v, ok := this.Option(key); ok {
				if m, ok := v.(map[string]interface{}); ok {
						return m
				}
		}
		return def[0]
}

func (this *BaseWrapper) SetOpt(key string, val interface{}) {
		if this.Options == nil {
				this.Options = make(map[string]interface{})
		}
		this.GetOrm()
		this.Options[key] = val
}

func (this *BaseWrapper) ClearOpt(keys ...string) {
		if this.Options == nil {
				return
		}
		if len(keys) == 0 {
				for k := range this.Options {
						delete(this.Options, k)
				}
				return
		}
		for _, k := range keys {
				delete(this.Options, k)
		}
}

func (this *BaseWrapper) GetFields() map[string]string {
		if this.Fields == nil || len(this.Fields) == 0 {
				if this.Model != nil {
						this.Fields = reflects.GetItemsAllTypes(this.Model)
				}
		}
		return this.Fields
}

func (this *BaseWrapper) HasField(key string) bool {
		if _, ok := this.GetFields()[key]; ok {
				return true
		}
		return false
}

func (this *BaseWrapper) GetError(clean ...bool) error {
		if len(clean) == 0 {
				clean = append(clean, true)
		}
		var err = this.Error
		if clean[0] {
				this.Error = nil
		}
		return err
}

func (this *BaseWrapper) Table() string {
		if this.Model == nil {
				return ""
		}
		if t, ok := this.Model.(Table); ok {
				return t.TableName()
		}
		// @todo snake camel
		return strings.ToLower(reflects.Name(this.Model))
}

func (this *BaseWrapper) NewModel(data ...map[string]interface{}) interface{} {
		if this.Model == nil {
				this.Error = errors.New("miss model")
				return nil
		}
		original := reflects.RealValue(reflect.ValueOf(this.Model))
		cpy := reflect.New(original.Type())
		model := cpy.Interface()
		if len(data) > 0 {
				this.Error = reflects.CopyMap2Struct(data[0], model)
		}
		return model
}

func WrapperInitOptions(wrapper ModelWrapper, options ...interface{}) {
		if len(options) == 0 {
				return
		}
		for _, v := range options {
				if str, ok := v.(string); ok && str != "" {
						wrapper.SetOpt("connection", str)
				}
				if m, ok := v.(map[string]interface{}); ok {
						for k, val := range m {
								wrapper.SetOpt(k, val)
						}
				}
		}
}

func (this *BaseWrapper) CreateUUid(table string, key string) string {
		times := 0
		uu := uuid.NewV4()
		seqId := uu.String()
		for this.GetOrm().QueryTable(table).Filter(key, seqId).Exist() {
				uu = uuid.NewV4()
				seqId = uu.String()
				if times > 5 {
						return seqId + strconv.Itoa(time.Now().Second())
				}
				times++
		}
		return seqId
}

func (this *BaseWrapper) FilterFields(arr []string) []string {
		if len(arr) == 0 {
				return arr
		}
		var (
				keys  []string
				table = this.Table()
		)
		fields := this.GetFields()
		for _, key := range arr {
				if strings.Contains(key, "*") {
						keys = append(keys, key)
						if key == "*" {
								break
						}
				}
				if strings.Contains(key, ".") {
						if !strings.Contains(key, table) {
								keys = append(keys, strings.Replace(key, ".", "__", -1))
								continue
						}
						key = strings.Replace(key, table+".", "", 1)
				}
				if _, ok := fields[key]; !ok {
						continue
				}
				keys = append(keys, key)
		}
		return keys
}

func (this *BaseWrapper) GetFieldKeys() []string {
		var keys []string
		for key, v := range this.GetFields() {
				if v == "N" {
						continue
				}
				keys = append(keys, key)
		}
		return keys
}

func (this *BaseWrapper) QueryResolver(builder orm.QuerySeter, allowFields []string, conditions map[string]interface{}) (query orm.QuerySeter, effects map[string]interface{}) {
		effects = make(map[string]interface{})
		for _, key := range allowFields {
				v, ok := conditions[key]
				if !ok || v == nil || v == "" || v == 0 {
						continue
				}
				simple, ok := v.([]string)
				if ok {
						builder, ok = this.ResolverQuerySimple(builder, key, simple)
						if ok {
								effects[key] = simple
						}
						continue
				}
				arr, ok := v.([][]string)
				if ok {
						builder, ok = this.ResolverQueryComplex(builder, key, arr)
						if ok {
								effects[key] = arr
						}
						continue
				}
		}
		return builder, effects
}

func (this *BaseWrapper) ResolverQuerySimple(builder orm.QuerySeter, key string, v []string) (query orm.QuerySeter, ok bool) {
		if len(v) == 0 {
				return builder, false
		}
		return builder.FilterRaw(key, strings.Join(v, " ")), true
}

func (this *BaseWrapper) ResolverQueryComplex(builder orm.QuerySeter, key string, v [][]string) (query orm.QuerySeter, ok bool) {
		var i = 0
		for _, it := range v {
				builder, ok = this.ResolverQuerySimple(builder, key, it)
				if ok {
						i++
				}
		}
		return builder, i > 0
}
