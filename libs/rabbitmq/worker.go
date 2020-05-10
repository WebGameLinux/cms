package rabbitmq

import (
		"encoding/json"
		"errors"
		"github.com/streadway/amqp"
		"regexp"
		"strings"
		"sync"
)

const HandlerNameKey = "_handler"

type QueueConsumerWorkerInterface interface {
		Consumer() ConsumerInterface
		Init() QueueConsumerWorkerInterface
		Register(name string, handle WatcherHandler)
		RegisterHandler(handler QueueHandler)
		Dispatcher(name string, message amqp.Delivery)
		GetError() error
		Stop() bool
		Run()
}

type ConsumerWorker struct {
		client     ClientInterface
		errs       []error
		On         <-chan bool // 只读
		State      int
		Config     string
		Ctr        chan<- bool // 只写
		WorkerNum  int
		WorkerPool WorkerPoolInterface
		Container  *HandlerContainer
}

type HandlerContainer struct {
		sync.RWMutex
		handlers map[string]WatcherHandler
}

func (this *HandlerContainer) Destroy() {
		this.Lock()
		defer this.Unlock()
		for v, _ := range this.handlers {
				delete(this.handlers, v)
		}
}

func (this *HandlerContainer) Len() int {
		this.Lock()
		defer this.Unlock()
		return len(this.handlers)
}

func (this *HandlerContainer) Register(name string, handle WatcherHandler) {
		this.Lock()
		defer this.Unlock()
		if _, ok := this.handlers[name]; ok {
				return
		}
		this.handlers[name] = handle
}

func (this *HandlerContainer) Resolver(msg amqp.Delivery) string {
		var (
				body = msg.Body
				m    = make(map[string]interface{})
		)
		if json.Unmarshal(body, &m) != nil {
				return ""
		}
		name := m[HandlerNameKey]
		if name == "" {
				return "*"
		}
		if n, ok := name.(string); ok {
				return n
		}
		return ""
}

func (this *HandlerContainer) Exec(name string, msg amqp.Delivery, ctx QueueConsumerWorkerInterface) {
		if name == "" {
				return
		}
		this.Lock()
		defer this.Unlock()
		if name == "*" {
				this.All(msg, ctx)
				return
		}
		if this.isGroupPatten(name) {
				this.Group(regexp.MustCompile(name), msg, ctx)
				return
		}
		handler, ok := this.handlers[name]
		if !ok {
				return
		}
		handler(msg.Body, msg, ctx)
}

func (this *HandlerContainer) All(msg amqp.Delivery, ctx QueueConsumerWorkerInterface) {
		for _, handler := range this.handlers {
				if !handler(msg.Body, msg, ctx) {
						break
				}
		}
}

func (this *HandlerContainer) Group(reg *regexp.Regexp, msg amqp.Delivery, ctx QueueConsumerWorkerInterface) {
		for n, handler := range this.handlers {
				if !reg.MatchString(n) {
						continue
				}
				if !handler(msg.Body, msg, ctx) {
						break
				}
		}
}

func (this *HandlerContainer) isGroupPatten(name string) bool {
		return strings.Contains(name, "*") || strings.Contains(name, "$") || strings.Contains(name, "^")
}

func NewHandlerContainer() *HandlerContainer {
		var container = new(HandlerContainer)
		container.handlers = make(map[string]WatcherHandler)
		return container
}

func NewConsumerWorker(config string) *ConsumerWorker {
		var (
				errs   []error
				worker = new(ConsumerWorker)
		)
		worker.State = 0
		worker.Config = config
		worker.errs = errs
		worker.Container = NewHandlerContainer()
		return worker
}

func (this *ConsumerWorker) Consumer() ConsumerInterface {
		if this.client == nil {
				this.initialize()
		}
		return NewConsumer(this.client)
}

func (this *ConsumerWorker) GetWorkNum() int {
		if this.WorkerNum == 0 {
				return 3
		}
		return this.WorkerNum
}

func (this *ConsumerWorker) Init() QueueConsumerWorkerInterface {
		var (
				errs   []error
				on     = make(chan bool, 2)
				config = NewWorkPoolConfigByJson([]byte(this.Config))
		)
		this.errs = errs
		this.On = on
		this.Ctr = on
		this.WorkerPool = NewWorkPool(config.MaxNum, config.Name)
		this.initialize()
		return this
}

func (this *ConsumerWorker) Register(name string, handle WatcherHandler) {
		this.Container.Register(name, handle)
}

func (this *ConsumerWorker) RegisterHandler(handler QueueHandler) {
		this.Container.Register(handler.Name(), handler.Resolver())
}

func (this *ConsumerWorker) GetError() error {
		if this.errs == nil || len(this.errs) == 0 {
				return nil
		}
		var (
				num = len(this.errs)
				err = this.errs[0]
		)
		if num == 1 {
				this.errs = this.errs[0:0]
		}
		if num > 1 {
				this.errs = this.errs[1:]
		}
		return err
}

func (this *ConsumerWorker) Stop() bool {
		this.Container.Destroy()
		this.client.Close()
		this.client = nil
		this.State = 4
		this.On = nil
		this.Config = ""
		this.WorkerNum = 0
		this.WorkerPool.Destroy()
		this.WorkerPool = nil
		return true
}

func (this *ConsumerWorker) Emit(state bool) {
		if this.Ctr == nil {
				this.Init()
		}
		go func(state bool) {
				this.Ctr <- state
		}(state)
}

func (this *ConsumerWorker) initialize() bool {
		if this.client != nil {
				this.client = NewRabbitmqClientByConfig(this.Config)
		}
		if this.Container == nil || this.Container.Len() == 0 {
				panic(errors.New("handler container is empty"))
		}
		this.WorkerPool.Init()
		return true
}

func (this *ConsumerWorker) Run() {
		if this.State == 0 {
				this.Init()
		}
		if this.State == 4 || this.State <= 0 {
				this.OnError(errors.New("worker state not support run"))
				return
		}
		defer this.Stop()
		consumer := this.Consumer()
		queue, ok := consumer.Consumer()
		if !ok || queue == nil {
				this.OnError(errors.New("worker get consumer queue failed"))
				return
		}
		this.start(queue)
}

func (this *ConsumerWorker) start(queue <-chan amqp.Delivery) {
		this.State = 2 // running
		for {
				select {
				case msg := <-queue:
						id := this.Container.Resolver(msg)
						// 回退不处理
						if this.State != 2 {
								this.client.Push(string(msg.Body))
								break
						}
						this.Dispatcher(id, msg)
				case on := <-this.On:
						if !on {
								goto EndWorker
						}
				}
		}
EndWorker:
}

func (this *ConsumerWorker) Dispatcher(name string, msg amqp.Delivery) {
		// 非启动状态不接受
		if this.State != 2 {
				return
		}
		if this.WorkerPool == nil {
				this.Init()
		}
		this.WorkerPool.Exec(func() {
				this.Container.Exec(name, msg, this)
		})
}

func (this *ConsumerWorker) OnError(err error) {
		if err == nil {
				return
		}
		this.errs = append(this.errs, err)
}

func (this *ConsumerWorker) GetState() int {
		return this.State
}
