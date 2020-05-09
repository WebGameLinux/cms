package test

import (
		"fmt"
		"github.com/WebGameLinux/cms/libs/rabbitmq"
		"testing"
		"time"
)

func TestNewSimpleModeMessageQueue(t *testing.T) {
		var mq = rabbitmq.NewSimpleClient()
		var queue = "my_rabbitmq"
		mq.Connector.ConnInfo.Username = "dev"
		mq.Connector.ConnInfo.Password = "dev123"
		mq.Connector.ConnInfo.VirtualHost = "DevGolang"
		mq.Connector.SetQueue(queue)
		ch, err := mq.Connector.GetConsumer()
		if err != nil {
				fmt.Println(err)
				return
		}
		defer  mq.Connector.Close()
		var i = 0
		for {
				select {
				case v := <-ch:
						fmt.Println(string(v.Body))
				case <-time.NewTicker(2 * time.Second).C:
						mq.Connector.Push(fmt.Sprintf("msg:%d,%s", i, time.Now().Format(time.RFC1123Z)))
						i++
				case <-time.NewTimer(10 * time.Minute).C:
						break
				}
		}
}
