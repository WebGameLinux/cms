package models

import (
		"errors"
		utils "github.com/WebGameLinux/cms/utils/beego"
		"github.com/WebGameLinux/cms/utils/reflects"
		"github.com/astaxie/beego/orm"
		validation "gopkg.in/go-playground/validator.v9"
		"reflect"
		"strings"
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

func (base *BaseWrapper) BindModel(m interface{}) {
		if m != nil && base.Model == nil {
				if reflect.TypeOf(m).Elem().Kind() == reflect.Struct {
						base.Model = m
				}
		}
}

func (base *BaseWrapper) GetOrm() orm.Ormer {
		var o = orm.NewOrm()
		connection := base.GetOrmUsing()
		err := o.Using(connection)
		// 重置
		if err != nil {
				// 强制 链接
				if base.OptBool("force") {
						ormOkSetter(err, o, connection)
						base.SetOpt("connection", "default")
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

func (base *BaseWrapper) GetOrmUsing() string {
		return base.OptStr("connection", "default")
}

func (base *BaseWrapper) GetQuery() (utils.SqlQueryBuilder, error) {
		var (
				err   error
				query orm.QueryBuilder
		)
		if query, err = orm.NewQueryBuilder(utils.GetDatabaseDriver()); err == nil {
				return utils.NewQueryBuilderWrapper(query, base.GetOrm(), base.NewModel()), nil
		}
		utils.Onerror(err)
		return nil, err
}

func (base *BaseWrapper) Option(key string) (interface{}, bool) {
		if base.Options == nil {
				return nil, false
		}
		if v, ok := base.Options[key]; ok {
				return v, true
		}
		return nil, false
}

func (base *BaseWrapper) OptStr(key string, def ...string) string {
		if len(def) == 0 {
				def = append(def, "")
		}
		if v, ok := base.Option(key); ok {
				if str, ok := v.(string); ok {
						return str
				}
		}
		return def[0]
}

func (base *BaseWrapper) OptInt(key string, def ...int) int {
		if len(def) == 0 {
				def = append(def, 0)
		}
		if v, ok := base.Option(key); ok {
				if num, ok := v.(int); ok {
						return num
				}
		}
		return def[0]
}

func (base *BaseWrapper) OptBool(key string, def ...bool) bool {
		if len(def) == 0 {
				def = append(def, false)
		}
		if v, ok := base.Option(key); ok {
				if b, ok := v.(bool); ok {
						return b
				}
		}
		return def[0]
}

func (base *BaseWrapper) OptFloat(key string, def ...float64) float64 {
		if len(def) == 0 {
				def = append(def, 0)
		}
		if v, ok := base.Option(key); ok {
				if b, ok := v.(float64); ok {
						return b
				}
		}
		return def[0]
}

func (base *BaseWrapper) OptMap(key string, def ...map[string]interface{}) map[string]interface{} {
		if len(def) == 0 {
				def = append(def, nil)
		}
		if v, ok := base.Option(key); ok {
				if m, ok := v.(map[string]interface{}); ok {
						return m
				}
		}
		return def[0]
}

func (base *BaseWrapper) SetOpt(key string, val interface{}) {
		if base.Options == nil {
				base.Options = make(map[string]interface{})
		}
		base.GetOrm()
		base.Options[key] = val
}

func (base *BaseWrapper) ClearOpt(keys ...string) {
		if base.Options == nil {
				return
		}
		if len(keys) == 0 {
				for k, _ := range base.Options {
						delete(base.Options, k)
				}
				return
		}
		for _, k := range keys {
				delete(base.Options, k)
		}
}

func (base *BaseWrapper) GetFields() map[string]string {
		if base.Fields == nil || len(base.Fields) == 0 {
				if base.Model != nil {
						base.Fields = reflects.GetItemsAllTypes(base.Model)
				}
		}
		return base.Fields
}

func (base *BaseWrapper) HasField(key string) bool {
		if _, ok := base.GetFields()[key]; ok {
				return true
		}
		return false
}

func (base *BaseWrapper) GetError(clean ...bool) error {
		if len(clean) == 0 {
				clean = append(clean, true)
		}
		var err = base.Error
		if clean[0] {
				base.Error = nil
		}
		return err
}

func (base *BaseWrapper) Table() string {
		if base.Model == nil {
				return ""
		}
		if t, ok := base.Model.(Table); ok {
				return t.TableName()
		}
		// @todo snake camel
		return strings.ToLower(reflects.Name(base.Model))
}

func (base *BaseWrapper) NewModel(data ...map[string]interface{}) interface{} {
		if base.Model == nil {
				base.Error = errors.New("miss model")
				return nil
		}
		original := reflects.RealValue(reflect.ValueOf(base.Model))
		cpy := reflect.New(original.Type())
		model := cpy.Interface()
		if len(data) > 0 {
				base.Error = reflects.CopyMap2Struct(data[0], model)
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
