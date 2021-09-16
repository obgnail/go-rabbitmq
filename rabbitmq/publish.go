package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/obgnail/go-rabbitmq/conn"
	"github.com/streadway/amqp"
)

// SamplePublishDurableMessage 使用最少的参数,发布简单的持久化消息
// 该方法的业务场景是 durable queue , durable message
// message 为可以序列为json, 并且确保消息已持久化成功, 能满足大部分消息持久化的需求
// exchangeType
func PublishSampleDurableMessage(
	exchangeName string,
	exchangeType string,
	routerKey string,
	deliveryMode uint8,
	expiration string,
	message interface{},
) error {
	exchange := newSampleDurableExchange(exchangeName, exchangeType)
	return PublishMessage(exchange, routerKey, deliveryMode, expiration, message)
}

// PublishMessage  使用较为复杂的参数，发布消息
// 该方法可以支持大多数的业务场景
// 其中 mandatory=false, immediate=false
// 如果想支持 mandatory=true的功能, 建议使用备用交换器方法替代
// 如果想支持 immediate=true的功能, 建议使用TTL和DLX方法代替
func PublishMessage(
	exchange *Exchange,
	routerKey string,
	deliveryMode uint8,
	expiration string,
	message interface{},
) (err error) {
	// 序列化
	bodyJson, err := json.Marshal(message)
	if err != nil {
		return err
	}
	amqpMsg := &amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "application/json",
		ContentEncoding: "",
		Body:            bodyJson,
		DeliveryMode:    amqp.Persistent,
		Priority:        0,
		Expiration:      expiration,
	}
	return PublishAmqpMessage(exchange, routerKey, deliveryMode, amqpMsg)
}

func PublishAmqpMessage(
	exchange *Exchange,
	routerKey string,
	deliveryMode uint8,
	msg *amqp.Publishing,
) (err error) {
	// exchange
	channel, closeChannel, err := conn.GetAmqpChannel()
	if err != nil {
		return
	}
	defer closeChannel()
	err = exchangeDeclare(channel, exchange)
	if err != nil {
		return
	}

	// 事务, 确保数据已经写入队列
	if err = channel.Confirm(false); err != nil {
		err = fmt.Errorf("channel could not be put into confirm mode: %s", err)
		return
	}
	confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	defer func() {
		e := confirmOne(confirms)
		if err == nil {
			err = e
		}
	}()

	// 消息推送
	if err = channel.Publish(
		exchange.Name,
		routerKey,
		false,
		false,
		*msg,
	); err != nil {
		return fmt.Errorf("exchange publish: %s", err)
	}
	return nil
}

func newSampleDurableExchange(exchangeName string, exchangeType string) *Exchange {
	return &Exchange{
		Name:       exchangeName,
		Type:       exchangeType,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
}

func confirmOne(confirms <-chan amqp.Confirmation) error {
	if confirmed := <-confirms; !confirmed.Ack {
		return fmt.Errorf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
	return nil
}

type DurationMessage struct {
	RetryCount int
	RetryMax   int
	Body       []byte
}

func getMessageBody(body []byte) ([]byte, error) {
	durationMessage := new(DurationMessage)
	err := json.Unmarshal(body, durationMessage)
	if err != nil {
		return nil, err
	}
	return durationMessage.Body, nil
}

func PublishDurationMessage(exchangeName, routerKey string, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	durationMessage := new(DurationMessage)
	durationMessage.RetryMax = 10
	durationMessage.Body = body

	msgBody, err := json.Marshal(durationMessage)
	if err != nil {
		return err
	}
	exchange := newSampleDurableExchange(exchangeName, ExchangeTypeDirect)
	return PublishAmqpMessage(exchange, routerKey, DeliveryModePersistent, buildPublishing(msgBody, 0))
}

func buildPublishing(body []byte, expiration int64) *amqp.Publishing {
	var e string
	if expiration != 0 {
		e = ExpirationSecond(expiration)
	}
	return &amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "application/json",
		ContentEncoding: "",
		Body:            body,
		DeliveryMode:    DeliveryModePersistent,
		Priority:        0,
		Expiration:      e,
	}
}
