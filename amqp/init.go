package amqp

import (
	"github.com/obgnail/go-rabbitmq/handler"
	"github.com/obgnail/go-rabbitmq/rabbitmq"
	"log"
	"time"
)

const (
	workCount = 3
)

func InitAmqp() error {
	initCustomers()
	repeatInitQueueBindingUntilSuccess()
	return nil
}

func initCustomers() {
	rabbitmq.ConsumeSampleMessage(handler.MySimpleQueue, handler.MyHandler, workCount)
	rabbitmq.ConsumeSampleMessage(handler.BusinessQueue, handler.MyDurableDelayHandler, workCount)
	rabbitmq.ConsumeSampleMessage(handler.SimplifyBusinessQueue, handler.MySimplifyDurableDelayHandler, workCount)
}

func repeatInitQueueBindingUntilSuccess() {
	go func() {
		displayLog := true
		for {
			err := initQueueBinding()
			if err != nil {
				if !displayLog {
					displayLog = true
					log.Printf("bind queue failed, %+v", err)
				}
				time.Sleep(20 * time.Second)
				continue
			}
			if displayLog {
				log.Println("bind queue success")
			}
			return
		}
	}()
}

func initQueueBinding() (err error) {
	// 普通持久化队列
	err = rabbitmq.BindSampleDurableQueue(
		handler.MySimpleExchange,
		rabbitmq.ExchangeTypeDirect,
		handler.MySimpleQueue,
		handler.MySimpleRouteKey,
	)
	if err != nil {
		return err
	}

	// 普通延时持久化队列
	err = rabbitmq.BindSampleDurableDelayQueue(
		handler.BusinessExchange,
		rabbitmq.ExchangeTypeDirect,
		handler.BusinessQueue,
		handler.BusinessRouteKey,

		handler.DelayExchange,
		rabbitmq.ExchangeTypeDirect,
		handler.DelayQueue,
		handler.DelayRouteKey,
	)
	if err != nil {
		return err
	}

	// 简化版的延时持久化队列
	err = rabbitmq.BindDurableQueue(
		handler.SimplifyBusinessExchange,
		rabbitmq.ExchangeTypeDirect,
		handler.SimplifyBusinessQueue,
		handler.SimplifyBusinessRouteKey,
	)
	if err != nil {
		return err
	}
	return
}
