package common

import (
		"encoding/json"
		"fmt"
		"github.com/WebGameLinux/cms/dto/enums"
		"github.com/WebGameLinux/cms/dto/response"
		utils "github.com/WebGameLinux/cms/utils/beego"
		"github.com/WebGameLinux/cms/utils/getter"
		"github.com/WebGameLinux/cms/utils/reflects"
		"github.com/WebGameLinux/cms/utils/setter"
		"reflect"
		"strings"
)

type Decoder interface {
		Unmarshal(v interface{}) error
}

type Encoder interface {
		Marshal() string
}

type LengthAble interface {
		Len() int
}

type CountAble interface {
		Count() int
}

type JsonResponseInterface interface {
		Code() int
		Message() string
		Error() error
		Value(key string, def ...interface{}) interface{}
		Item() interface{}
		fmt.Stringer
		Mapper() map[string]interface{}
		Json() string
		Decoder
		Encoder
		IsSuccess() bool
		Set(string, interface{})
		GetBoot() func(interface{})
		GetDtoStruct() interface{}
		Coder() enums.Code
		IsEmpty() bool
}

type Number interface {
		Num() int
}

// 基础数据对象
type BaseResultDto struct {
		response.RespJson
}

// 服务结果数据对象
type ServiceResultDto struct {
		BaseResultDto
		err         error
		Encoder     func() string
		Decoder     func(v interface{}) error
		bootHandler func(v interface{})
		autoInit    bool
}

func NewResult(init ...bool) JsonResponseInterface {
		if len(init) == 0 || !init[0] {
				return new(ServiceResultDto)
		}
		return new(ServiceResultDto).Init()
}

func defaultBootResult(this *ServiceResultDto) {
		this.RespJson.Code = 0
		this.RespJson.Msg = ""
		this.err = nil
		this.RespJson.Data = nil
}

func (this *ServiceResultDto) boot() JsonResponseInterface {
		this.GetBoot()(this)
		return this
}

// 初始化相关处理器
func (this *ServiceResultDto) GetBoot() func(v interface{}) {
		if this.bootHandler == nil {
				this.bootHandler = func(v interface{}) {
						if dto, ok := v.(*ServiceResultDto); ok {
								defaultBootResult(dto)
						}
				}
		}
		return this.bootHandler
}

// 初始化
func (this *ServiceResultDto) Init() JsonResponseInterface {
		return this.boot()
}

// 序列化 json
func (this *ServiceResultDto) Marshal() string {
		if this.Encoder != nil {
				return this.Encoder()
		}
		return this.Json()
}

// 获取 业务码
func (this *ServiceResultDto) Code() int {
		return this.RespJson.Code
}

func (this *ServiceResultDto) Coder() enums.Code {
		return enums.Code(this.Code())
}

// 获取业务提示信息
func (this *ServiceResultDto) Message() string {
		return this.RespJson.Msg
}

// 获取业务异常
func (this *ServiceResultDto) Error() error {
		return this.err
}

// 获取业务值
func (this *ServiceResultDto) Value(key string, def ...interface{}) interface{} {
		var (
				err error
				v   interface{}
		)
		if len(def) == 0 {
				def = append(def, nil)
		}
		if key == "" {
				return this.RespJson.Data
		}
		g := getter.GetAnyGetter(this.RespJson.Data)
		if g != nil {
				if v, err = g.GetValue(key); err == nil {
						return v
				}
				utils.Onerror(err)
		}
		return def[0]
}

// 获取可以导出的数据
func (this *ServiceResultDto) Export() *BaseResultDto {
		var data = &BaseResultDto{}
		data.RespJson.Data = this.RespJson.Data
		data.RespJson.Code = this.RespJson.Code
		data.RespJson.Msg = this.RespJson.Msg
		return data
}

// 获取纯数据结果
func (this *ServiceResultDto) GetDtoStruct() interface{} {
		return this.Export()
}

// 获取返回内容体
func (this *ServiceResultDto) Item() interface{} {
		return this.RespJson.Data
}

// 字符串输出
func (this *ServiceResultDto) String() string {
		return this.Json()
}

// 转为 mapper
func (this *ServiceResultDto) Mapper() map[string]interface{} {
		var v = make(map[string]interface{})
		if this.Encoder != nil {
				if err := json.Unmarshal([]byte(this.Encoder()), &v); err != nil {
						utils.Onerror(err)
				}
		}
		if buf, err := this.MarshalJSON(); err == nil {
				_ = json.Unmarshal(buf, &v)
		}
		return v
}

// 输出json
func (this *ServiceResultDto) Json() string {
		var (
				buf []byte
				err error
		)
		if this.Encoder != nil {
				return this.Encoder()
		}
		if buf, err = this.MarshalJSON(); err != nil {
				utils.Onerror(err)
				return ""
		}
		return string(buf)
}

// 反序列化为为对应结构对象
func (this *ServiceResultDto) Unmarshal(v interface{}) error {
		if this.Decoder != nil {
				return this.Decoder(v)
		}
		return json.Unmarshal([]byte(this.Json()), v)
}

// 业务结果是否成功
func (this *ServiceResultDto) IsSuccess() bool {
		if this.Error() != nil {
				return false
		}
		if this.Coder().Equal(enums.SUCCESS) {
				return true
		}
		return true
}

func (this *ServiceResultDto) IsEmpty() bool {
		if this.RespJson.Data == nil {
				return true
		}
		switch this.RespJson.Data.(type) {
		case CountAble:
				c := this.RespJson.Data.(CountAble)
				return c.Count() == 0
		case LengthAble:
				l := this.RespJson.Data.(LengthAble)
				return l.Len() == 0
		}
		ty := reflects.RealType(reflect.TypeOf(this.RespJson.Data))
		if ty.Kind() == reflect.Array || ty.Kind() == reflect.Slice {
				return ty.Len() == 0
		}
		return true
}

// 设置返回结果体
func (this *ServiceResultDto) Set(key string, v interface{}) {
		switch key {
		case "code":
				if code, ok := v.(int); ok {
						this.RespJson.Code = code
				}
				if num, ok := v.(Number); ok {
						this.RespJson.Code = num.Num()
				}
		case "message":
				fallthrough
		case "msg":
				if msg, ok := v.(string); ok {
						this.RespJson.Msg = msg
				}
				if str, ok := v.(fmt.Stringer); ok {
						this.RespJson.Msg = str.String()
				}
		case "err":
				fallthrough
		case "error":
				if err, ok := v.(error); ok {
						this.err = err
				}
		case "items":
				fallthrough
		case "data":
				this.RespJson.Data = v
		default:
				this.UpdateItem(key, v)
		}
}

// 更新返回数据
func (this *ServiceResultDto) UpdateItem(key string, v interface{}) {
		if key == "" || !strings.Contains(key, "data.") {
				return
		}
		set := setter.GetAnySetter(this.RespJson.Data)
		if set == nil {
				return
		}
		set.SetValue(key, v)
}
