package rabbitmq

import (
	"fmt"
	"github.com/obgnail/go-rabbitmq/conn"
	"log"
	"runtime/debug"
	"time"
)

// HandleMessage 处理消息，实际执行的业务逻辑
// 返回参数：
//    requeue: 消息是否重新入队
//    err: 错误信息
type HandleMessage func([]byte) (requeue bool, err error)

// ConsumeSampleMessage 使用最少的参数,对消息进行消费
// 该方法能处理比较简单的业务场景
// 参数:
//    queueName: 队列名称
//    handle: 执行实际的业务逻辑的函数
//    workCount: 该进程中。对该队列进行消费的消费者数量
func ConsumeSampleMessage(queueName string, handle HandleMessage, workCount int) {
	sampleConsume := newSampleConsume(queueName)
	ConsumeMessage(sampleConsume, handle, workCount)
}

// ConsumeMessage 使用较为复杂的参数,对消息进行消费
// 该方法能处理更多的业务场景
// 参数:
//    consume: 消费者配置信息
//    handle: 执行实际的业务逻辑的函数
//    workCount: 该进程中。对该队列进行消费的消费者数量
func ConsumeMessage(consume *Consume, handle HandleMessage, workCount int) {
	for i := 0; i < workCount; i++ {
		go func(i int) {
			for {
				err := consumeMessage(consume, handle)
				if i == 0 && err != nil {
					//  控制重复log的数量
					log.Printf("consume message : %+v", err)
				}
				time.Sleep(10 * time.Second)
			}
		}(i)
	}
}

func consumeMessage(consume *Consume, handle HandleMessage) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("%s", debug.Stack())
		}
	}()

	channel, closeChannel, err := conn.GetAmqpChannel()
	if err != nil {
		return err
	}
	defer closeChannel()

	deliveries, err := channelConsume(channel, consume)
	if err != nil {
		return err
	}
	for d := range deliveries {
		if len(d.Body) == 0 {
			return fmt.Errorf("miss body")
		}

		requeue, err := handle(d.Body)
		if err != nil {
			if !requeue {
				log.Printf("consumeMessage message error, message.body:%s", string(d.Body))
			} else {
				log.Println(err)
			}
			d.Reject(requeue)
		} else {
			d.Ack(true)
		}
	}
	return nil
}

func newSampleConsume(queueName string) *Consume {
	return &Consume{
		Name:      queueName,
		Tag:       "",
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}
}
