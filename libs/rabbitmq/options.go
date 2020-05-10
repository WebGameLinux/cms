package rabbitmq

import (
		"encoding/json"
		"regexp"
		"strings"
		"time"
)

const (
		WorkModeSimpler   = "simpler"      // 单队列-单消费
		WorkModeWorker    = "worker"       // 单队列-多消费,每个消息只能被消费端消费一次
		WorkModePublisher = "publisher"    //  订阅模式
		WorkModeTopic     = "topic"        //  话题模式
		WorkModeRouting   = "routing"      //  路由模式
		ExchangeFanOut    = "fanout"       // 广播
		ExchangePublish   = ExchangeFanOut // 订阅
		ExchangeDirect    = "direct"       // 直接
		ExchangeSimple    = ExchangeDirect // 简单
		ExchangeRouting   = ExchangeDirect // 路由模式
		ExchangeTopic     = "topic"        // 话题
)

// rabbitmq 工作模式
var Modes = []string{
		WorkModeSimpler,
		WorkModeWorker,
		WorkModePublisher,
		WorkModeTopic,
		WorkModeRouting,
}

// 参数
type Options struct {
		Queue      string                 `json:"queue"`       // 消息队列名
		Consumer   string                 `json:"consumer"`    // 指定消费者
		AutoAck    bool                   `json:"auto_ack"`    // 是否自动回复
		Exclusive  bool                   `json:"exclusive"`   // 排他性, 仅创建进程可见
		NoLocal    bool                   `json:"no_local"`    // 是否允许用同一个connection 生成消息也消费消息
		NoWait     bool                   `json:"no_wait"`     // 队列消费是否阻塞
		Durable    bool                   `json:"durable"`     // 是否消息持久化
		AutoDelete bool                   `json:"auto_delete"` // 是否自动删除,当最后一个消费者断开后
		Args       map[string]interface{} `json:"table"`       // 其他可选参数
		Mandatory  bool                   `json:"mandatory"`   // 队列不存在,消息回退给发送端
		Immediate  bool                   `json:"immediate"`   // 消费者不存在, 消息回退
		Internal   bool                   `json:"internal"`    //  true: exchange不可以被client用来推送消息,仅用来进行exchange和exchange之间的绑定
		WorkMode   string                 `json:"work_mode"`   // 工作模式 simpler,worker,publisher
}

// 链接信息
type ConnOptions struct {
		Username    string `json:"username"`     // 用户名  [手动在rabbitmq中创建用户]
		Password    string `json:"password"`     // 密码    [手动在rabbitmq中创建的用户对应密码]
		Host        string `json:"host"`         // 主机    [rabbitmq 服务器地址]
		Port        int    `json:"port"`         // 端口    [rabbitmq 对外服务的端口]
		VirtualHost string `json:"virtual_host"` // 虚拟主机 [手动创建]
}

var defaultConnInfo = new(ConnOptions)
var defaultOptions = NewSimpleOptions()

const (
		DefaultHost = "127.0.0.1"
		DefaultPort = 5672
)

func NewOptions() *Options {
		options := new(Options)
		options.Args = make(map[string]interface{})
		return options
}

func NewSimpleOptions() *Options {
		opts := NewOptions()
		FormatterSimpleOption(opts)
		return opts
}

func FormatterSimpleOption(options *Options) bool {
		options.Consumer = ""
		options.NoLocal = false
		options.NoWait = false
		options.AutoAck = true
		options.Exclusive = false
		return true
}

func (this *ConnOptions) GetUserName() string {
		return this.Username
}

func (this *ConnOptions) GetPassword() string {
		return this.Password
}

func (this *ConnOptions) GetHost() string {
		if this.Host == "" {
				return DefaultHost
		}
		return this.Host
}

func (this *ConnOptions) GetVirtualHost() string {
		return this.VirtualHost
}

func (this *ConnOptions) GetPort() int {
		if this.Port == 0 {
				return DefaultPort
		}
		return this.Port
}

func (this *ConnOptions) InitByMapper(m map[string]interface{}) *ConnOptions {
		if len(m) == 0 {
				return this
		}
		data, err := json.Marshal(m)
		if err == nil {
				return this.InitByJson(data)
		}
		return this
}

func (this *ConnOptions) InitByJson(data []byte) *ConnOptions {
		if len(data) == 0 {
				return this
		}
		_ = this.UnmarshalJSON(data)
		return this
}

func (this *ConnOptions) Copy(other *ConnOptions) *ConnOptions {
		if other == nil || other == this {
				return this
		}
		this.Username = other.Username
		this.Password = other.Password
		this.Host = other.Host
		this.Port = other.Port
		this.VirtualHost = other.VirtualHost
		return this
}

func (this *Options) Get(key string) (interface{}, bool) {
		switch key {
		case "queue":
				fallthrough
		case "Queue":
				return this.Queue, true
		case "consumer":
				fallthrough
		case "Consumer":
				return this.Consumer, true
		case "autoAck":
				fallthrough
		case "AutoAck":
				return this.AutoAck, true
		case "exclusive":
				fallthrough
		case "Exclusive":
				return this.Exclusive, true
		case "noLocal":
				fallthrough
		case "NoLocal":
				return this.NoLocal, true
		case "NoWait":
				fallthrough
		case "noWait":
				return this.NoWait, true
		case "durable":
				fallthrough
		case "Durable":
				return this.Durable, true
		case "autoDelete":
				fallthrough
		case "AutoDelete":
				return this.AutoDelete, true
		case "Args":
				fallthrough
		case "args":
				return this.Args, true
		case "mode":
				fallthrough
		case "Mode":
				return this.WorkMode, true
		case "internal":
				fallthrough
		case "Internal":
				return this.Internal, true
		}
		if strings.Contains(key, ".") && regexp.MustCompile(`^(args|Args)\.\w+`).MatchString(key) {
				return this.getInArgs(key)
		}
		return nil, false
}

func (this *Options) getInArgs(key string) (interface{}, bool) {
		var scopes = strings.SplitN(key, ".", -1)
		if len(scopes) <= 1 {
				return nil, false
		}
		return ArgsGet(this.Args, scopes[1:])
}

func (this *Options) GetContentType(def ...string) string {
		if len(def) == 0 {
				def = append(def, "text/plain")
		}
		v, ok := this.Get("args.content_type")
		if !ok {
				return def[0]
		}
		if t, ok := v.(string); ok {
				return t
		}
		return def[0]
}

func (this *Options) GetContentEncoding(def ...string) string {
		if len(def) == 0 {
				def = append(def, "utf8")
		}
		v, ok := this.Get("args.content_encoding")
		if !ok {
				return def[0]
		}
		if t, ok := v.(string); ok {
				return t
		}
		return def[0]
}

func (this *Options) GetDurAble() bool {
		return this.Durable
}

func (this *Options) GetNoWait() bool {
		return this.NoWait
}

func (this *Options) GetAutoDelete() bool {
		return this.AutoDelete
}

func (this *Options) GetInternal() bool {
		return this.Internal
}

func (this *Options) GetWorkMode() string {
		if this.WorkMode == "" {
				return WorkModeSimpler
		}
		return this.WorkMode
}

func (this *Options) SetWorkMode(mode string) bool {
		if mode == "" {
				return false
		}
		for _, v := range Modes {
				if v == mode {
						this.WorkMode = mode
						return true
				}
		}
		return false
}

func (this *Options) GetExclusive() bool {
		return this.Exclusive
}

func (this *Options) GetAutoAck() bool {
		return this.AutoAck
}

func (this *Options) GetNoLocal() bool {
		return this.NoLocal
}

func (this *Options) GetArgs() map[string]interface{} {
		if this.Args == nil || len(this.Args) == 0 {
				return nil
		}
		return this.Args
}

func (this *Options) GetQueue() string {
		return this.Queue
}

func (this *Options) GetMandatory() bool {
		return this.Mandatory
}

func (this *Options) GetImmediate() bool {
		return this.Immediate
}

func (this *Options) GetConsumer() string {
		return this.Consumer
}

func (this *Options) InitByJson(data []byte) *Options {
		if len(data) == 0 {
				return this
		}
		_ = this.UnmarshalJSON(data)
		return this
}

func (this *Options) InitByMapper(data map[string]interface{}) *Options {
		buf, err := json.Marshal(data)
		if err != nil {
				return this
		}
		return this.InitByJson(buf)
}

func (this *Options) Copy(opt *Options) *Options {
		if opt == nil || opt == this {
				return this
		}
		this.Queue = opt.Queue
		this.Consumer = opt.Consumer
		this.AutoAck = opt.AutoAck
		this.Exclusive = opt.Exclusive
		this.NoLocal = opt.NoLocal
		this.NoWait = opt.NoWait
		this.Durable = opt.Durable
		this.AutoDelete = opt.AutoDelete
		this.Args = opt.Args
		this.Mandatory = opt.Mandatory
		this.Immediate = opt.Immediate
		this.WorkMode = opt.WorkMode
		return this
}

func (this *Options) GetDefault(key string, def ...interface{}) interface{} {
		if len(def) == 0 {
				def = append(def, nil)
		}
		v, ok := this.Get(key)
		if ok {
				return v
		}
		return def[0]
}

func ArgsGet(ctx map[string]interface{}, keys []string) (interface{}, bool) {
		var times = len(keys)
		if times == 0 {
				return nil, false
		}
		for i, key := range keys {
				v, ok := ctx[key]
				if !ok {
						return nil, false
				}
				if i+1 >= times {
						return v, true
				}
				if m, ok := v.(map[string]interface{}); ok {
						return ArgsGet(m, keys[i+1:])
				}
				return nil, false
		}
		return nil, false
}

func GetDefaultOption() *Options {
		return defaultOptions
}

func GetDefaultConnInfo() *ConnOptions {
		return defaultConnInfo
}

func NewConnByJson(conn string) *ConnOptions {
		opts := new(ConnOptions)
		return opts.InitByJson([]byte(conn))
}

func NewOptionsByJson(opt string) *Options {
		opts := new(Options)
		return opts.InitByJson([]byte(opt))
}

type ConfigObject struct {
		Mode    string `json:"mode"`
		Conn    string `json:"conn"`
		Options string `json:"options"`
}

func NewConfigObject() *ConfigObject {
		object := new(ConfigObject)
		return object
}

func CreateConfigObjectByJson(data []byte) *ConfigObject {
		object := NewConfigObject()
		_ = object.UnmarshalJSON(data)
		return object
}

type WorkPoolConfig struct {
		Name               string        `json:"name"`
		MaxNum             int           `json:"max_num"`
		Interval           time.Duration `json:"interval"`
		CacheMaxNum        int           `json:"cache_max_num"`
		CheckCacheInterval time.Duration `json:"check_cache_interval"`
}

func NewWorkPoolConfig() *WorkPoolConfig {
		config := new(WorkPoolConfig)
		return config
}

func NewWorkPoolConfigByJson(data []byte) *WorkPoolConfig {
		config := NewWorkPoolConfig()
		_ = config.UnmarshalJSON(data)
		return config
}
