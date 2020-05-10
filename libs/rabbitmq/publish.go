package rabbitmq

import (
		"errors"
		"fmt"
		"github.com/streadway/amqp"
)

type PublishConnector struct {
		SimpleConnector
		Exchanged bool
		Bended    bool
}

type PublishConnectorInterface interface {
		SimpleConnectorInterface
		IsPublish() bool
		Exchanger() bool
		BindQueueExchanger(args ...string)
}

func NewPublisher(exchangeName string) *PublishConnector {
		var connector = new(PublishConnector)
		connector.ConnInfo = new(ConnOptions)
		connector.Option = NewSimpleOptions()
		connector.SetExchange(exchangeName)
		connector.Option.Durable = true
		connector.Option.SetWorkMode(WorkModePublisher)
		return connector
}

func NewPublishConnector(exchangeName string, connJson string, optionJson string) PublishConnectorInterface {
		connector := NewPublisher(exchangeName)
		connector.ConnInfo.InitByJson([]byte(connJson))
		connector.Option.InitByJson([]byte(optionJson))
		return connector
}

func (this *PublishConnector) GetQueue() string {
		if this.IsPublish() {
				return ""
		}
		return this.SimpleConnector.GetQueue()
}

func (this *PublishConnector) IsPublish() bool {
		return this.Option.GetWorkMode() == WorkModePublisher
}

func (this *PublishConnector) GetExchange() string {
		if !this.IsPublish() {
				return this.SimpleConnector.GetExchange()
		}
		if this.Exchange == "" {
				exchange, ok := this.Option.getInArgs("args.exchange")
				if !ok || exchange == nil {
						panic("exchange is empty")
				}
				v, ok := exchange.(string)
				if ! ok {
						panic(errors.New(fmt.Sprintf("exchange type error: %+v, %T", v, v)))
				}
				this.Exchange = v
		}
		return this.Exchange
}

func (this *PublishConnector) GetConsumer() (<-chan amqp.Delivery, error) {
		this.Queue(true)
		if len(this.GetErrors()) > 0 {
				return nil, this.GetError()
		}
		if this.Listener != nil {
				return *this.Listener, nil
		}
		var (
				err     error
				msgChan <-chan amqp.Delivery
		)
		if msgChan, err = this.channel.Consume(
				this.QueueIns.Name,
				this.Option.GetConsumer(),
				this.Option.GetAutoAck(),
				this.Option.GetExclusive(),
				this.Option.GetNoLocal(),
				this.Option.GetNoLocal(),
				this.Option.GetArgs(),
		); err != nil {
				return nil, err
		}
		this.Listener = &msgChan
		return msgChan, nil
}

func (this *PublishConnector) Queue(isConsumer ...bool) *amqp.Queue {
		if len(isConsumer) == 0 {
				isConsumer = append(isConsumer, false)
		}
		if !this.IsPublish() {
				return this.SimpleConnector.Queue(isConsumer...)
		}
		if !this.Connected() {
				this.Connect()
		}
		if !this.Exchanged {
				this.Exchanged = this.Exchanger()
		}
		if !isConsumer[0] {
				return this.QueueIns
		}
		this.SimpleConnector.Queue()
		if isConsumer[0] && !this.Bended && this.QueueIns != nil {
				this.Option.Consumer = ""
				this.BindQueueExchanger()
		}
		return this.QueueIns
}

func (this *PublishConnector) Exchanger() bool {

		err := this.channel.ExchangeDeclare(
				this.GetExchange(),
				ExchangePublish,
				this.Option.GetDurAble(),
				this.Option.GetAutoDelete(),
				this.Option.GetInternal(),
				this.Option.GetNoWait(),
				this.Option.GetArgs())

		if err == nil {
				return true
		}
		this.OnError(err)
		return false
}

// 绑定队列 到 交互器
func (this *PublishConnector) BindQueueExchanger(args ...string) {
		if this.Bended {
				return
		}
		var (
				argc     = len(args)
				exchange = this.GetExchange()
				queue    = this.getQueueName()
		)
		if argc > 0 {
				switch argc {
				case 1:
						queue = args[0]
				case 2:
						queue = args[0]
						exchange = args[2]
				}
		}
		//绑定队列到 exchange 中
		err := this.channel.QueueBind(
				queue,
				"", //在pub/sub模式下，这里的key要为空
				exchange,
				this.Option.GetNoWait(),
				this.Option.GetArgs())
		if err == nil {
				this.Bended = true
		}
		return
}

func (this *PublishConnector) getQueueName() string {
		if this.QueueIns != nil {
				return this.QueueIns.Name
		}
		this.SimpleConnector.Queue()
		return this.QueueIns.Name
}
