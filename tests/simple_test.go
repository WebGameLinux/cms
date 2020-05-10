package test

import (
		"fmt"
		"github.com/WebGameLinux/cms/libs/rabbitmq"
		"testing"
		"time"
)

func TestNewSimpleModeMessageQueue(t *testing.T) {
		var (
				queue  = "my_rabbitmq"
				config = `{"username":"dev","password":"dev123","host":"127.0.0.1","virtual_host":"DevGolang"}`
				mq     = rabbitmq.NewSimpleClient(queue, config, "")
		)
		mq.GetSimpleConnector().GetConnInfo().InitByJson([]byte(config))
		mq.GetSimpleConnector().SetQueue(queue)
		ch, err := mq.Consumer()
		if err != nil {
				fmt.Println(err)
				return
		}
		defer mq.Close()
		var i = 0
		for {
				select {
				case v := <-ch:
						fmt.Println(string(v.Body))
				case <-time.NewTicker(2 * time.Second).C:
						mq.Push(fmt.Sprintf("msg:%d,%s", i, time.Now().Format(time.RFC1123Z)))
						i++
				case <-time.NewTimer(10 * time.Minute).C:
						break
				}
		}
}

func TestNewPublish(t *testing.T) {
		var ex = "test"
		var conn = `{"username":"dev","password":"dev123","host":"127.0.0.1","virtual_host":"DevGolang"}`
		var connector = rabbitmq.NewPublishConnector(ex, conn, "")
		consumer := rabbitmq.NewConsumer(rabbitmq.NewClient(connector))
		producer := rabbitmq.NewProducer(rabbitmq.NewClient(connector))
		defer producer.Close()
		defer consumer.Close()
		ch, _ := consumer.Consumer()
		fmt.Println(producer.Push("start publish"))
		var i = 0
		for {
				select {
				case v := <-ch:
						fmt.Println(string(v.Body))
				case <-time.NewTicker(2 * time.Second).C:
						producer.Push(fmt.Sprintf("msg:%d,%s", i, time.Now().Format(time.RFC1123Z)))
						i++
				case <-time.NewTimer(1 * time.Minute).C:
						break
				}
		}
}
