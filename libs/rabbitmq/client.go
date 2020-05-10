package rabbitmq

import "github.com/streadway/amqp"

type Client struct {
		connector SimpleConnectorInterface
}

type Producer struct {
		client ClientInterface
}

type Consumer struct {
		client ClientInterface
}

type ClientInterface interface {
		Push(string) bool
		Consumer() (<-chan amqp.Delivery, bool)
		Close()
		Connector() SimpleConnectorInterface
}

// 外部参数扩展键
const (
		ArgsRouteKey    = "args.route_key"
		ArgsExchangeKey = "args.exchange"
)

func NewClient(connector SimpleConnectorInterface) ClientInterface {
		return &Client{connector: connector}
}

func (this *Client) Push(msg string) bool {
		return this.Connector().Push(msg)
}

func (this *Client) Close() {
		this.Connector().Close()
}

func (this *Client) Consumer() (<-chan amqp.Delivery, bool) {
		if ch, err := this.Connector().GetConsumer(); err == nil {
				return ch, true
		}
		return nil, false
}

func (this *Client) Connector() SimpleConnectorInterface {
		return this.connector
}

func (this *Producer) Push(msg string) bool {
		return this.client.Push(msg)
}

func (this *Producer) Close() {
		this.client.Close()
}

func (this *Consumer) Consumer() (<-chan amqp.Delivery, bool) {
		return this.client.Consumer()
}

func (this *Consumer) Close() {
		this.client.Close()
}

func NewProducer(client ClientInterface) ProducerInterface {
		return &Producer{client: client}
}

func NewConsumer(client ClientInterface) ConsumerInterface {
		return &Consumer{client: client}
}

func NewConnectorAuto(mode string, connJson string, optionJson string) SimpleConnectorInterface {
		var (
				queue   string
				options = NewOptionsByJson(optionJson)
		)
		if options.Queue != "" {
				queue = options.Queue
				if buf, err := options.MarshalJSON(); err == nil {
						optionJson = string(buf)
				}
		}
		switch mode {
		case WorkModeSimpler:
				return NewSimpleConnector(queue, connJson, optionJson)
		case WorkModeWorker:
				return NewWorkConnector(queue, connJson, optionJson)
		case WorkModePublisher:
				v := options.GetDefault(ArgsExchangeKey)
				name, ok := v.(string)
				if !ok {
						return nil
				}
				return NewPublishConnector(name, connJson, optionJson)
		case WorkModeTopic:
				v := options.GetDefault(ArgsExchangeKey)
				name, ok := v.(string)
				if !ok {
						return nil
				}
				r := options.GetDefault(ArgsRouteKey)
				route, ok := r.(string)
				return NewTopicConnector(name, route, connJson, optionJson)
		case WorkModeRouting:
				v := options.GetDefault(ArgsExchangeKey)
				name, ok := v.(string)
				if !ok {
						return nil
				}
				r := options.GetDefault(ArgsRouteKey)
				route, ok := r.(string)
				if options.Queue != "" {
						if buf, err := options.MarshalJSON(); err == nil {
								optionJson = string(buf)
						}
				}
				return NewRouterConnector(name, route, connJson, optionJson)
		default:
				return NewSimpleConnector(queue, connJson, optionJson)
		}
}

func NewRabbitmqClientByConfig(config string) ClientInterface {
		object := CreateConfigObjectByJson([]byte(config))
		return NewClient(NewConnectorAuto(object.Mode, object.Conn, object.Options))
}
