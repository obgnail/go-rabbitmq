package handler

import (
	"fmt"
	"github.com/obgnail/go-rabbitmq/rabbitmq"
)

const (
	MySimpleExchange = "MySimpleExchange"
	MySimpleQueue    = "MySimpleQueue"
	MySimpleRouteKey = "MySimpleRouteKey"

	BusinessExchange = "BusinessExchange"
	BusinessQueue    = "BusinessQueue"
	BusinessRouteKey = "BusinessRouteKey"
	DelayExchange    = "MyDurableDelayExchange"
	DelayQueue       = "MyDurableDelayQueue"
	DelayRouteKey    = BusinessRouteKey

	SimplifyBusinessExchange = "SimplifyBusinessExchange"
	SimplifyBusinessQueue    = "SimplifyBusinessQueue"
	SimplifyBusinessRouteKey = "SimplifyBusinessRouteKey"
)

func MyHandler(body []byte) (requeue bool, err error) {
	fmt.Println("--- MyHandler", string(body))
	return
}

func MyDurableDelayHandler(body []byte) (requeue bool, err error) {
	fmt.Println("--- MyDurableDelayHandler", string(body))
	return
}

func MySimplifyDurableDelayHandler(body []byte) (requeue bool, err error) {
	fmt.Println("--- MySimplifyDurableDelayHandler", string(body))
	return
}

type MyMessage struct {
	Name    string
	Content string
}

func NewMyMessage(name, content string) *MyMessage {
	return &MyMessage{name, content}
}

func PublishMessageToMQ(msg *MyMessage) error {
	err := rabbitmq.PublishSampleDurableMessage(
		MySimpleExchange,
		rabbitmq.ExchangeTypeDirect,
		MySimpleRouteKey,
		rabbitmq.DeliveryModePersistent,
		rabbitmq.NoExpiration(),
		msg,
	)
	return err
}

func PublishMessageToDelayMQ(msg *MyMessage) error {
	err := rabbitmq.PublishSampleDurableMessage(
		DelayExchange,
		rabbitmq.ExchangeTypeDirect,
		DelayRouteKey,
		rabbitmq.DeliveryModePersistent,
		rabbitmq.ExpirationSecond(10),
		msg,
	)
	return err
}

func PublishSimplifyMessageToDelayMQ(msg *MyMessage) error {
	delayExchangeName := fmt.Sprintf("%s-%s", "delay", SimplifyBusinessExchange)
	err := rabbitmq.PublishSampleDurableMessage(
		delayExchangeName,
		rabbitmq.ExchangeTypeDirect,
		SimplifyBusinessRouteKey,
		rabbitmq.DeliveryModePersistent,
		rabbitmq.ExpirationSecond(10),
		msg,
	)
	return err
}
