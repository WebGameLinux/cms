package rabbitmq

import (
		"github.com/streadway/amqp"
)

type SimpleConnector struct {
		Connector
		QueueIns *amqp.Queue
		Listener *<-chan amqp.Delivery
}

type SimpleConnectorInterface interface {
		ConnectorInterface
		Push(string) bool                                      // 发送消息
		GetConsumer() (<-chan amqp.Delivery, error)            // 获取消费channel
		Queue() *amqp.Queue                                    // 创建获取队列
		CreateSimpleMessage(string, ...string) amqp.Publishing // 创建消息
}

func NewSimpleModeMessageQueue() SimpleConnectorInterface {
		var connector = new(SimpleConnector)
		connector.ConnInfo = new(ConnOptions)
		connector.Option = NewSimpleOptions()
		return connector
}

func NewSimpleClient() *SimpleClient {
		var connector = new(SimpleConnector)
		connector.ConnInfo = new(ConnOptions)
		connector.Option = NewSimpleOptions()
		return &SimpleClient{Connector: connector}
}

type SimpleClient struct {
		Connector *SimpleConnector
}

func (this *SimpleConnector) Push(message string) bool {
		if this.Queue() == nil {
				return false
		}
		if err := this.channel.Publish(
				this.GetExchange(),
				this.GetQueue(),
				this.Option.GetMandatory(),
				this.Option.GetImmediate(),
				this.CreateSimpleMessage(message),
		); err != nil {
				this.OnError(err)
				return false
		}
		return true
}

func (this *SimpleConnector) CreateSimpleMessage(message string, options ...string) amqp.Publishing {
		var (
				body            = []byte(message)
				argc            = len(options)
				contentType     = this.Option.GetContentType()
				contentEncoding = this.Option.GetContentEncoding()
		)
		if argc >= 1 {
				contentType = options[0]
		}
		if argc >= 2 {
				contentEncoding = options[1]
		}
		return amqp.Publishing{
				ContentType:     contentType,
				ContentEncoding: contentEncoding,
				Body:            body,
		}
}

func (this *SimpleConnector) Queue() *amqp.Queue {
		if this.QueueIns != nil {
				return this.QueueIns
		}
		if !this.Connected() {
				this.Connect()
		}
		var (
				err error
				q   amqp.Queue
		)
		q, err = this.channel.QueueDeclare(
				this.GetQueue(),
				this.Option.GetDurAble(),
				this.Option.GetAutoDelete(),
				this.Option.GetExclusive(),
				this.Option.GetNoWait(),
				this.Option.GetArgs(),
		)
		if err != nil {
				this.OnError(err)
				return nil
		}
		this.QueueIns = &q
		return this.QueueIns
}

func (this *SimpleConnector) GetKey() string {
		return ""
}

func (this *SimpleConnector) GetConsumer() (<-chan amqp.Delivery, error) {
		if this.Queue() == nil {
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
				this.GetQueue(),
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

func (this *SimpleConnector) Close() {
		this.Connector.Close()
		this.QueueIns = nil
		this.Listener = nil
}
