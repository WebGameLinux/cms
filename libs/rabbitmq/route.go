package rabbitmq

import "github.com/streadway/amqp"

type RouterConnector struct {
		SimpleConnector
		Exchanged bool
		Bended    bool
}

type RouterConnectorInterface interface {
		SimpleConnectorInterface
		IsRouter() bool
		Exchanger() bool
		BindQueueExchanger(args ...string)
}

func NewRouterConnector(exchangeName string, routingKey string, connJson string, optionJson string) RouterConnectorInterface {
		connector := NewRoute(exchangeName, routingKey)
		connector.ConnInfo.InitByJson([]byte(connJson))
		connector.Option.InitByJson([]byte(optionJson))
		return connector
}

func NewRoute(exchangeName string, routingKey string) *RouterConnector {
		var connector = new(RouterConnector)
		connector.ConnInfo = new(ConnOptions)
		connector.Option = NewSimpleOptions()
		connector.SetKey(routingKey)
		connector.SetExchange(exchangeName)
		connector.Option.Durable = true
		connector.Option.SetWorkMode(WorkModeRouting)
		return connector
}

func (this *RouterConnector) IsRouter() bool {
		return this.Option.GetWorkMode() == WorkModeRouting
}

func (this *RouterConnector) Exchanger() bool {
		err := this.channel.ExchangeDeclare(
				this.GetExchange(),
				ExchangeRouting,
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

func (this *RouterConnector) BindQueueExchanger(args ...string) {
		if this.Bended {
				return
		}
		var (
				argc     = len(args)
				key      = this.GetKey()
				exchange = this.GetExchange()
				queue    = this.getQueueName()
		)
		if argc > 0 {
				switch argc {
				case 1:
						queue = args[0]
				case 2:
						queue = args[0]
						key = args[1]
				case 3:
						queue = args[0]
						key = args[1]
						exchange = args[2]
				}
		}
		//绑定队列到 exchange 中
		err := this.channel.QueueBind(
				queue,
				key, //在pub/sub模式下，这里的key要为空
				exchange,
				this.Option.GetNoWait(),
				this.Option.GetArgs())
		if err == nil {
				this.Bended = true
		}
		return
}

func (this *RouterConnector) GetConsumer() (<-chan amqp.Delivery, error) {
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

func (this *RouterConnector) Queue(isConsumer ...bool) *amqp.Queue {
		if len(isConsumer) == 0 {
				isConsumer = append(isConsumer, false)
		}
		if !this.IsRouter() {
				return this.SimpleConnector.Queue()
		}
		if !this.Connected() {
				this.Connect()
		}
		if !this.Exchanged {
				this.Exchanged = this.Exchanger()
		}
		this.SimpleConnector.Queue()
		if isConsumer[0] && !this.Bended && this.QueueIns != nil {
				this.BindQueueExchanger()
		}
		return this.QueueIns
}

func (this *RouterConnector) getQueueName() string {
		if this.QueueIns != nil {
				return this.QueueIns.Name
		}
		this.SimpleConnector.Queue()
		return this.QueueIns.Name
}
