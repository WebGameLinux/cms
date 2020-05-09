package rabbitmq

import (
		"regexp"
		"strings"
)

type Options struct {
		Queue      string                 `json:"queue"`
		Consumer   string                 `json:"consumer"`
		AutoAck    bool                   `json:"auto_ack"`
		Exclusive  bool                   `json:"exclusive"`
		NoLocal    bool                   `json:"no_local"`
		NoWait     bool                   `json:"no_wait"`
		Durable    bool                   `json:"durable"`
		AutoDelete bool                   `json:"auto_delete"`
		Args       map[string]interface{} `json:"table"`
		Mandatory  bool                   `json:"mandatory"`
		Immediate  bool                   `json:"immediate"`
}

// 链接信息
type ConnOptions struct {
		Username    string `json:"username"`     // 用户名
		Password    string `json:"password"`     // 密码
		Host        string `json:"host"`         // 主机
		Port        int    `json:"port"`         // 端口
		VirtualHost string `json:"virtual_host"` // 虚拟主机
}

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
		}
		if strings.Contains(key, ".") && regexp.MustCompile(`/^(args|Args)\.\w+/`).MatchString(key) {
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
