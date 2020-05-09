package rabbitmq

import (
		"fmt"
		log "github.com/sirupsen/logrus"
		"github.com/streadway/amqp"
		"reflect"
)

type Connector struct {
		channel   *amqp.Channel
		conn      *amqp.Connection
		QueueName string
		Exchange  string
		Key       string
		Option    *Options
		ConnInfo  *ConnOptions
		Errors    []error
}

type ConnectorInterface interface {
		SetOption(*Options) ConnectorInterface
		GetOption() *Options
		SetKey(key string) ConnectorInterface
		GetKey() string
		GetQueue() string
		SetQueue(string) ConnectorInterface
		GetExchange() string
		SetExchange(string) ConnectorInterface
		SetConnInfo(*ConnOptions) ConnectorInterface
		GetConnInfo() *ConnOptions
		Connected() bool
		ReConnect() ConnectorInterface
		Close()
		GetConnectionUrl() string
		Connect() ConnectorInterface
		GetErrors() []error
}

const UrlConn = "amqp://%s:%s@%s:%d/%s"

// 创建Connector
func NewConnector() *Connector {
		var connector = &Connector{}
		connector.Option = NewOptions()
		connector.ConnInfo = new(ConnOptions)
		return connector
}

func (this *Connector) SetOption(opts *Options) ConnectorInterface {
		this.Option = opts
		return this
}

func (this *Connector) GetOption() *Options {
		return this.Option
}

func (this *Connector) destroy() {
		var err error
		if this.channel != nil {
				if err = this.channel.Close(); err != nil {
						this.OnError(err)
				}
				this.channel = nil
		}
		if this.conn != nil {
				if err = this.conn.Close(); err != nil {
						this.OnError(err)
				}
				this.conn = nil
		}
}

func (this *Connector) GetErrors() []error {
		var err = this.Errors
		this.Errors = this.Errors[0:0]
		return err
}

func (this *Connector) Close() {
		this.destroy()
}

// 错误日志
func (this *Connector) failOnErr(err error) {
		if err != nil {
				log.Errorf(reflect.TypeOf(this).PkgPath()+"error: %s", err)
		}
}

// 获取链接 link
func (this *Connector) GetConnectUrl() string {
		return fmt.Sprintf(UrlConn,
				this.ConnInfo.GetUserName(),
				this.ConnInfo.GetPassword(),
				this.ConnInfo.GetHost(),
				this.ConnInfo.GetPort(),
				this.ConnInfo.GetVirtualHost(),
		)
}

func (this *Connector) Connect() ConnectorInterface {
		var (
				err  error
				link = this.GetConnectUrl()
		)

		if this.conn != nil && this.channel != nil {
				return this
		}
		if this.conn == nil {
				this.conn, err = amqp.Dial(link)
				this.OnError(err)
				if this.channel != nil {
						_ = this.channel.Close()
				}
		}
		this.channel, err = this.conn.Channel()
		this.OnError(err)
		return this
}

func (this *Connector) OnError(err error) {
		if err == nil {
				return
		}
		this.Errors = append(this.Errors, err)
		this.failOnErr(err)
}

func (this *Connector) ReConnect() ConnectorInterface {
		this.Close()
		return this.Connect()
}

func (this *Connector) Open() *Connector {
		if !this.Connected() {
				this.Connect()
		}
		return this
}

func (this *Connector) Connected() bool {
		return this.channel != nil && this.conn != nil
}

func (this *Connector) SetKey(key string) ConnectorInterface {
		this.Key = key
		return this
}

func (this *Connector) GetKey() string {
		return this.Key
}

func (this *Connector) GetQueue() string {
		return this.QueueName
}

func (this *Connector) SetQueue(queue string) ConnectorInterface {
		this.QueueName = queue
		return this
}

func (this *Connector) GetExchange() string {
		return this.Exchange
}

func (this *Connector) SetExchange(exchange string) ConnectorInterface {
		this.Exchange = exchange
		return this
}

func (this *Connector) SetConnInfo(info *ConnOptions) ConnectorInterface {
		if info == nil {
				return this
		}
		this.ConnInfo = info
		return this
}

func (this *Connector) GetConnInfo() *ConnOptions {
		return this.ConnInfo
}

func (this *Connector) GetConnectionUrl() string {
		return this.GetConnectUrl()
}

func (this *Connector) GetError() error {
		num := len(this.Errors)
		if num > 0 {
				err := this.Errors[0]
				if num > 1 {
						this.Errors = this.Errors[1:]
				} else {
						this.Errors = this.Errors[0:0]
				}
				return err
		}
		return nil
}
