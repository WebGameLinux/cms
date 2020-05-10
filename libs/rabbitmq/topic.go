package rabbitmq

import (
		"errors"
		"fmt"
		"github.com/streadway/amqp"
)

type TopicConnector struct {
		SimpleConnector
		Exchanged bool
		Bended    bool
}

type TopicConnectorInterface interface {
		SimpleConnectorInterface
		IsTopic() bool
		Exchanger() bool
		BindQueueExchanger(args ...string)
}

func NewTopic(exchangeName string, routingKey string) *TopicConnector {
		var connector = new(TopicConnector)
		connector.ConnInfo = new(ConnOptions)
		connector.Option = NewSimpleOptions()
		connector.SetKey(routingKey)
		connector.SetExchange(exchangeName)
		connector.Option.Durable = true
		connector.Option.SetWorkMode(WorkModeTopic)
		return connector
}

func NewTopicConnector(exchangeName string, routingKey string, connJson string, optionJson string) TopicConnectorInterface {
		connector := NewTopic(exchangeName, routingKey)
		connector.ConnInfo.InitByJson([]byte(connJson))
		connector.Option.InitByJson([]byte(optionJson))
		return connector
}

func (this *TopicConnector) IsTopic() bool {
		return this.Option.GetWorkMode() == WorkModeTopic
}

func (this *TopicConnector) GetKey() string {
		if !this.IsTopic() {
				return this.Connector.GetKey()
		}
		if this.Key == "" {
				v, ok := this.Option.getInArgs("args.router")
				if !ok {
						panic(errors.New("topic empty router key"))
				}
				k, ok := v.(string)
				if ok {
						panic(errors.New(fmt.Sprintf("topic router key type error,%+v,%T", k, k)))
				}
				this.Key = k
		}
		return this.Key
}

func (this *TopicConnector) GetExchange() string {
		if !this.IsTopic() {
				return this.Connector.GetExchange()
		}
		if this.Exchange == "" {
				v, ok := this.Option.getInArgs("args.exchange")
				if !ok {
						panic(errors.New("topic empty exchange"))
				}
				e, ok := v.(string)
				if ok {
						panic(errors.New(fmt.Sprintf("topic exchange type error,%+v,%T", e, e)))
				}
				this.Exchange = e
		}
		return this.Exchange
}

func (this *TopicConnector) GetConsumer() (<-chan amqp.Delivery, error) {
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

func (this *TopicConnector) Queue(isConsumer ...bool) *amqp.Queue {
		if len(isConsumer) == 0 {
				isConsumer = append(isConsumer, false)
		}
		if !this.IsTopic() {
				return this.SimpleConnector.Queue(isConsumer...)
		}
		if !this.Connected() {
				this.Connect()
		}
		if !this.Exchanged {
				this.Exchanged = this.Exchanger()
		}
		this.SimpleConnector.Queue()
		if isConsumer[0] && !this.Bended && this.QueueIns != nil {
				this.Option.Consumer = ""
				this.BindQueueExchanger()
		}
		return this.QueueIns
}

func (this *TopicConnector) GetQueue() string {
		return ""
}

// 绑定队列 到 交互器
func (this *TopicConnector) BindQueueExchanger(args ...string) {
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

func (this *TopicConnector) Exchanger() bool {
		err := this.channel.ExchangeDeclare(
				this.GetExchange(),
				ExchangeTopic,
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

func (this *TopicConnector) getQueueName() string {
		if this.QueueIns != nil {
				return this.QueueIns.Name
		}
		this.SimpleConnector.Queue()
		return this.QueueIns.Name
}
