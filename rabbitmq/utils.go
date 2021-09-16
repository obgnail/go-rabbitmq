package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
)

const (
	ExchangeTypeDirect  = "direct"
	ExchangeTypeFanout  = "fanout"
	ExchangeTypeTopic   = "topic"
	ExchangeTypeHeaders = "headers"

	// DeliveryModeTransient 消息非持久化
	DeliveryModeTransient = amqp.Transient
	// DeliveryModePersistent 消息持久化
	DeliveryModePersistent = amqp.Persistent
)

func ExpirationMillisecond(millisecond int64) string {
	return fmt.Sprintf("%d", millisecond)
}

func ExpirationSecond(second int64) string {
	return fmt.Sprintf("%d", second*1000)
}

// NoExpiration 无过期时间
func NoExpiration() string {
	return ""
}

func queueDeclare(channel *amqp.Channel, queue *Queue) error {
	_, err := channel.QueueDeclare(
		queue.Name,       // name of the queue
		queue.Durable,    // durable
		queue.AutoDelete, // delete when unused
		queue.Exclusive,  // exclusive
		queue.NoWait,     // noWait
		queue.Args,       // arguments
	)
	if err != nil {
		return fmt.Errorf("queue Declare: %s", err)
	}
	return nil
}

func bindQueue(channel *amqp.Channel, queueBind *QueueBind) error {
	if err := channel.QueueBind(
		queueBind.Name,
		queueBind.Key,
		queueBind.Exchange,
		queueBind.NoWait,
		queueBind.Args,
	); err != nil {
		return fmt.Errorf("queue Bind: %s", err)
	}
	return nil
}

func exchangeDeclare(channel *amqp.Channel, exchange *Exchange) error {
	if err := channel.ExchangeDeclare(
		exchange.Name,
		exchange.Type,
		exchange.Durable,
		exchange.AutoDelete,
		exchange.Internal,
		exchange.NoWait,
		exchange.Args,
	); err != nil {
		return fmt.Errorf("exchange declare: %s", err)
	}
	return nil
}

func channelConsume(channel *amqp.Channel, consume *Consume) (<-chan amqp.Delivery, error) {
	deliveries, err := channel.Consume(
		consume.Name,
		consume.Tag,
		consume.AutoAck,
		consume.Exclusive,
		consume.NoLocal,
		consume.NoWait,
		consume.Args,
	)
	if err != nil {
		return nil, fmt.Errorf("queue consume: %s", err)
	}
	return deliveries, nil
}
