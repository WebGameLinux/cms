package services

import (
		"encoding/json"
		"fmt"
		utils "github.com/WebGameLinux/cms/utils/beego"
		"github.com/WebGameLinux/cms/utils/getter"
		"github.com/WebGameLinux/cms/utils/mapper"
		"github.com/WebGameLinux/cms/utils/setter"
		string2 "github.com/WebGameLinux/cms/utils/string"
		"reflect"
		"strings"
)

type Map mapper.Mapper

type Decoder interface {
		Unmarshal(v interface{}) error
}

type Encoder interface {
		Marshal() string
}

type ResultStruct interface {
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
}

type Number interface {
		Num() int
}

// 基础数据对象
type BaseResultDto struct {
		ErrNo     int         `json:"code"`
		ErrMsg    string      `json:"message"`
		ValueItem interface{} `json:"data"`
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

func NewResult(init ...bool) ResultStruct {
		if len(init) == 0 || !init[0] {
				return new(ServiceResultDto)
		}
		return new(ServiceResultDto).Init()
}

func defaultBootResult(dto *ServiceResultDto) {
		dto.ErrNo = 0
		dto.ErrMsg = ""
		dto.err = nil
		dto.ValueItem = nil
}

func (dto *ServiceResultDto) boot() ResultStruct {
		dto.GetBoot()(dto)
		return dto
}

// 初始化相关处理器
func (dto *ServiceResultDto) GetBoot() func(v interface{}) {
		if dto.bootHandler == nil {
				dto.bootHandler = func(v interface{}) {
						if dto, ok := v.(*ServiceResultDto); ok {
								defaultBootResult(dto)
						}
				}
		}
		return dto.bootHandler
}

// 初始化
func (dto *ServiceResultDto) Init() ResultStruct {
		return dto.boot()
}

// 序列化 json
func (dto *ServiceResultDto) Marshal() string {
		if dto.Encoder != nil {
				return dto.Encoder()
		}
		return dto.Json()
}

// 获取 业务码
func (dto *ServiceResultDto) Code() int {
		return dto.ErrNo
}

// 获取业务提示信息
func (dto *ServiceResultDto) Message() string {
		return dto.ErrMsg
}

// 获取业务异常
func (dto *ServiceResultDto) Error() error {
		return dto.err
}

// 获取业务值
func (dto *ServiceResultDto) Value(key string, def ...interface{}) interface{} {
		var (
				err error
				v   interface{}
		)
		if len(def) == 0 {
				def = append(def, nil)
		}
		if key == "" {
				return dto.ValueItem
		}
		g := getter.GetAnyGetter(dto.ValueItem)
		if g != nil {
				if v, err = g.GetValue(key); err == nil {
						return v
				}
				utils.Onerror(err)
		}
		return def[0]
}

// 获取可以导出的数据
func (dto *ServiceResultDto) Export() *BaseResultDto {
		var data = &BaseResultDto{}
		data.ValueItem = dto.ValueItem
		data.ErrNo = dto.ErrNo
		data.ErrMsg = dto.ErrMsg
		return data
}

// 获取纯数据结果
func (dto *ServiceResultDto) GetDtoStruct() interface{} {
		return dto.Export()
}

// 获取返回内容体
func (dto *ServiceResultDto) Item() interface{} {
		return dto.ValueItem
}

// 字符串输出
func (dto *ServiceResultDto) String() string {
		return dto.Json()
}

// 转为 mapper
func (dto *ServiceResultDto) Mapper() map[string]interface{} {
		var v = make(map[string]interface{})
		if dto.Encoder != nil {
				if err := json.Unmarshal([]byte(dto.Encoder()), &v); err != nil {
						utils.Onerror(err)
				}
		}
		data := dto.GetDtoStruct()
		if data == nil {
				return v
		}
		ty := reflect.TypeOf(data)
		val := reflect.ValueOf(data)
		if val.Kind() != reflect.Struct {
				if val.Elem().Kind() != reflect.Struct {
						return v
				}
				ty = ty.Elem()
				val = val.Elem()
		}
		nums := val.NumField()
		for i := 0; i < nums; i++ {
				name := ty.Field(i).Tag.Get("json")
				if name == "" {
						name = ty.Field(i).Name
				} else {
						name = string2.StrSplitFirst(name, ";", ",", "|")
				}
				v[name] = val.Field(i).Interface()
		}
		return v
}

// 输出json
func (dto *ServiceResultDto) Json() string {
		var (
				buf []byte
				err error
		)
		if dto.Encoder != nil {
				return dto.Encoder()
		}
		if buf, err = json.Marshal(dto.GetDtoStruct()); err != nil {
				utils.Onerror(err)
				return ""
		}
		return string(buf)
}

// 反序列化为为对应结构对象
func (dto *ServiceResultDto) Unmarshal(v interface{}) error {
		if dto.Decoder != nil {
				return dto.Decoder(v)
		}
		return json.Unmarshal([]byte(dto.Json()), v)
}

// 业务结果是否成功
func (dto *ServiceResultDto) IsSuccess() bool {
		if dto.Error() != nil {
				return false
		}
		return true
}

// 设置返回结果体
func (dto *ServiceResultDto) Set(key string, v interface{}) {
		switch key {
		case "code":
				if code, ok := v.(int); ok {
						dto.ErrNo = code
				}
				if num, ok := v.(Number); ok {
						dto.ErrNo = num.Num()
				}
		case "message":
				fallthrough
		case "msg":
				if msg, ok := v.(string); ok {
						dto.ErrMsg = msg
				}
				if str, ok := v.(fmt.Stringer); ok {
						dto.ErrMsg = str.String()
				}
		case "err":
				fallthrough
		case "error":
				if err, ok := v.(error); ok {
						dto.err = err
				}
		case "items":
				fallthrough
		case "data":
				dto.ValueItem = v
		default:
				dto.UpdateItem(key, v)
		}
}

// 更新返回数据
func (dto *ServiceResultDto) UpdateItem(key string, v interface{}) {
		if key == "" || !strings.Contains(key, "data.") {
				return
		}
		set := setter.GetAnySetter(dto.ValueItem)
		if set == nil {
				return
		}
		set.SetValue(key, v)
}
