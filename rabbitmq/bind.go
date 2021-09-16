package rabbitmq

import (
	"fmt"
	"github.com/obgnail/go-rabbitmq/conn"
	"github.com/streadway/amqp"
)

func BindSampleDurableQueue(exchangeName string, exchangeType string, queueName string, bindingKey string) error {
	exchange := newSampleDurableExchange(exchangeName, exchangeType)
	queue := newSampleDurableQueue(queueName)
	queueBind := newSampleQueueBind(exchangeName, queueName, bindingKey)
	return BindQueue(exchange, queue, queueBind)
}

// 延迟队列和延迟交换器统一添加delay后缀
func BindDurableQueue(exchangeName string, exchangeType string, queueName string, bindingKey string) error {
	delayExchangeName := fmt.Sprintf("%s-%s", "delay", exchangeName)
	delayQueueName := fmt.Sprintf("%s-%s", "delay", queueName)
	return BindSampleDurableDelayQueue(
		exchangeName, exchangeType, queueName, bindingKey,
		delayExchangeName, exchangeType, delayQueueName, bindingKey,
	)
}

// 绑定延时队列 (延时时间，由业务决定)
// 所有发送到 delayQueueName 的消息， 满足死信条件（比如消息超过 Expiration过期时间时）下， 将会把消息转发到 exchangeName 上
func BindSampleDurableDelayQueue(exchangeName string, exchangeType string, queueName string, bindingKey string,
	delayExchangeName, delayExchangeType string, delayQueueName, delayBindingKey string) error {
	err := BindSampleDurableQueue(exchangeName, exchangeType, queueName, bindingKey)
	if err != nil {
		return err
	}
	delayExchange := newSampleDurableExchange(delayExchangeName, delayExchangeType)
	delayQueue := &Queue{
		Name:       delayQueueName,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args: amqp.Table{
			// 当消息过期时把消息发送到 exchangeName 这个交换器上
			"x-dead-letter-exchange": exchangeName,
		},
	}
	delayQueueBind := newSampleQueueBind(delayExchangeName, delayQueueName, delayBindingKey)
	return BindQueue(delayExchange, delayQueue, delayQueueBind)
}

func BindQueue(exchange *Exchange, queue *Queue, queueBing *QueueBind) (err error) {
	channel, close, err := conn.GetAmqpChannel()
	if err != nil {
		return
	}
	defer close()
	// 声明 exchange
	err = exchangeDeclare(channel, exchange)
	if err != nil {
		return
	}

	// 声明 queue
	err = queueDeclare(channel, queue)
	if err != nil {
		return
	}

	// 绑定
	err = bindQueue(channel, queueBing)
	if err != nil {
		return
	}
	return nil
}

func newSampleDurableQueue(queueName string) *Queue {
	return &Queue{
		Name:       queueName,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
}

func newSampleQueueBind(exchangeName string, queueName string, bindKey string) *QueueBind {
	return &QueueBind{
		Name:     queueName,
		Key:      bindKey,
		Exchange: exchangeName,
		NoWait:   false,
		Args:     nil,
	}
}
