package rabbitmq

import (
		"fmt"
		"github.com/streadway/amqp"
)

type WatcherHandler func(data []byte, args ...interface{}) bool

type QueueHandler interface {
		Name() string
		Resolver() WatcherHandler
		Handler() bool
		Call(data amqp.Delivery)
}

type QueueJobRabbitmqInterface interface {
		ConnInfo() *ConnOptions
		Options() *Options
		Dispatch(data fmt.Stringer, async bool) bool
		Push(string, async bool) bool
}

type JobDto struct {
		ConnInfoJson string
		OptionJson   string
		Async        bool
}

type SyncJob struct {
		JobDto
		Conn   *ConnOptions
		Option *Options
		Mode   string
		Client ClientInterface
}

type AsyncJob struct {
		JobDto
		Conn   *ConnOptions
		Option *Options
		Mode   string
		Client ClientInterface
}

func NewSync(mode string, connJson string, optionJson string) *SyncJob {
		job := new(SyncJob)
		job.Async = false
		job.ConnInfoJson = connJson
		job.OptionJson = optionJson
		job.Mode = mode
		return job
}

func (this *SyncJob) ConnInfo() *ConnOptions {
		if this.Conn == nil {
				this.Conn = new(ConnOptions)
				this.Conn.InitByJson([]byte(this.ConnInfoJson))
		}
		return this.Conn
}

func (this *SyncJob) Options() *Options {
		if this.Option == nil {
				this.Option = NewOptions()
				this.Option.InitByJson([]byte(this.OptionJson))
		}
		return this.Option
}

func (this *SyncJob) Dispatch(data fmt.Stringer, async bool) bool {
		if async != this.Async {
				return false
		}
		return this.client().Push(data.String())
}

func (this *SyncJob) Push(data string, async bool) bool {
		if async != this.Async {
				return false
		}
		return this.client().Push(data)
}

func (this *SyncJob) client() ClientInterface {
		if this.Client == nil {
				this.Client = NewClient(NewConnectorAuto(this.Mode, this.ConnInfoJson, this.OptionJson))
		}
		return this.Client
}

func NewAsync(mode string, connJson string, optionJson string) *AsyncJob {
		job := new(AsyncJob)
		job.Async = true
		job.ConnInfoJson = connJson
		job.OptionJson = optionJson
		job.Mode = mode
		return job
}

func (this *AsyncJob) ConnInfo() *ConnOptions {
		if this.Conn == nil {
				this.Conn = new(ConnOptions)
				this.Conn.InitByJson([]byte(this.ConnInfoJson))
		}
		return this.Conn
}

func (this *AsyncJob) Options() *Options {
		if this.Option == nil {
				this.Option = NewOptions()
				this.Option.InitByJson([]byte(this.OptionJson))
		}
		return this.Option
}

func (this *AsyncJob) Dispatch(data fmt.Stringer, async bool) bool {
		if async != this.Async {
				return false
		}
		go func(job string) {
				this.client().Push(job)
		}(data.String())
		return true
}

func (this *AsyncJob) Push(data string, async bool) bool {
		if async != this.Async {
				return false
		}
		go func(job string) {
				this.client().Push(job)
		}(data)
		return true
}

func (this *AsyncJob) client() ClientInterface {
		if this.Client == nil {
				this.Client = NewClient(NewConnectorAuto(this.Mode, this.ConnInfoJson, this.OptionJson))
		}
		return this.Client
}
