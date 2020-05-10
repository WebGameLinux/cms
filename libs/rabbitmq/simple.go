package rabbitmq

import (
		"errors"
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
		Queue(...bool) *amqp.Queue                             // 创建获取队列
		CreateSimpleMessage(string, ...string) amqp.Publishing // 创建消息
		IsSimple() bool
		Mode() string
}

func NewSimpleConnector(queue string, connJson string, optionJson string) SimpleConnectorInterface {
		var connector = NewSimple(queue)
		connector.ConnInfo.InitByJson([]byte(connJson))
		connector.Option.InitByJson([]byte(optionJson))
		return connector
}

func NewSimple(queue string) *SimpleConnector {
		var connector = new(SimpleConnector)
		connector.ConnInfo = new(ConnOptions)
		connector.Option = NewSimpleOptions()
		connector.SetQueue(queue)
		connector.Option.SetWorkMode(WorkModeSimpler)
		return connector
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

func (this *SimpleConnector) Queue(isConsumer ...bool) *amqp.Queue {
		if len(isConsumer) == 0 {
				isConsumer = append(isConsumer, false)
		}
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

func (this *SimpleConnector) IsSimple() bool {
		return this.Option.GetWorkMode() == WorkModeSimpler
}

func (this *SimpleConnector) Mode() string {
		return this.Option.GetWorkMode()
}

func (this *SimpleConnector) GetQueue() string {
		queue := this.Connector.GetQueue()
		if this.IsSimple() && queue == "" {
				panic(errors.New("simple mode queue not allow empty"))
		}
		return queue
}

func (this *SimpleConnector) GetKey() string {
		if this.IsSimple() {
				return ""
		}
		return this.Connector.GetKey()
}

func (this *SimpleConnector) GetConsumer() (<-chan amqp.Delivery, error) {
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

func (this *SimpleConnector) Close() {
		this.Connector.Close()
		this.QueueIns = nil
		this.Listener = nil
}

type SimpleClient struct {
		connector ConnectorInterface
}

func NewSimpleClient(queue string, connJson string, optionJson string) *SimpleClient {
		return &SimpleClient{connector: NewSimpleConnector(queue, connJson, optionJson)}
}

func (this *SimpleClient) SetConnector(connector ConnectorInterface) bool {
		if _, ok := connector.(SimpleConnectorInterface); !ok {
				return false
		}
		this.connector = connector
		return true
}

func (this *SimpleClient) GetConnector() ConnectorInterface {
		return this.connector
}

func (this *SimpleClient) GetSimpleConnector() SimpleConnectorInterface {
		return this.connector.(SimpleConnectorInterface)
}

func (this *SimpleClient) Push(msg string) bool {
		return this.GetSimpleConnector().Push(msg)
}

func (this *SimpleClient) Consumer() (<-chan amqp.Delivery, error) {
		return this.GetSimpleConnector().GetConsumer()
}

func (this *SimpleClient) GetErrors() []error {
		return this.connector.GetErrors()
}

func (this *SimpleClient) Close() {
		this.connector.Close()
}

type ConsumerInterface interface {
		Consumer() (<-chan amqp.Delivery, bool)
		Close()
}

type ProducerInterface interface {
		Push(string) bool
		Close()
}

type SimpleProducer struct {
		client *SimpleClient
}

type SimpleConsumer struct {
		client *SimpleClient
}

func NewSimpleProducer(client *SimpleClient) ProducerInterface {
		return &SimpleProducer{client: client}
}

func NewSimpleConsumer(client *SimpleClient) ConsumerInterface {
		return &SimpleConsumer{client: client}
}

func (this *SimpleProducer) Push(msg string) bool {
		return this.client.Push(msg)
}

func (this *SimpleProducer) Close() {
		this.client.Close()
}

func (this *SimpleConsumer) Consumer() (<-chan amqp.Delivery, bool) {
		if ch, err := this.client.Consumer(); err == nil {
				return ch, true
		}
		return nil, false
}

func (this *SimpleConsumer) Close() {
		this.client.Close()
}
