package rabbitmq

import "github.com/streadway/amqp"

type WorkersConnector struct {
		SimpleConnector
}

type WorkerConnectorInterface interface {
		SimpleConnectorInterface
		IsWork() bool
}

func NewWork(queue string) *WorkersConnector {
		var connector = new(WorkersConnector)
		connector.ConnInfo = new(ConnOptions)
		connector.Option = NewSimpleOptions()
		connector.SetQueue(queue)
		connector.Option.SetWorkMode(WorkModeWorker)
		return connector
}

func NewWorkConnector(queue string, connJson string, optionJson string) WorkerConnectorInterface {
		connector := NewWork(queue)
		connector.ConnInfo.InitByJson([]byte(connJson))
		connector.Option.InitByJson([]byte(optionJson))
		return connector
}

func (this *WorkersConnector) IsWork() bool {
		return this.Option.GetWorkMode() == WorkModeWorker
}

func (this *WorkersConnector) GetConsumer() (<-chan amqp.Delivery, error) {
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
				"",
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
