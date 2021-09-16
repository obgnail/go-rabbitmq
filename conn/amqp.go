package conn

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/obgnail/go-rabbitmq/config"
	"github.com/streadway/amqp"
)

var (
	amqpConnLock sync.Mutex
	amqpConn     *amqp.Connection
)

// GetAmqpChannel 用完channel后请调用 close 关闭
func GetAmqpChannel() (channel *amqp.Channel, closeChannel func(), err error) {
	if amqpConn == nil {
		return tryGetAmqpChannelAgain()
	}
	channel, err = amqpConn.Channel()
	if err != nil {
		return tryGetAmqpChannelAgain()
	}
	closeChannel = func() {
		channel.Close()
	}
	return
}

func tryGetAmqpChannelAgain() (channel *amqp.Channel, close func(), err error) {
	initAmqpConn()
	if amqpConn == nil {
		err = fmt.Errorf("amqp connetion pending")
		return
	}
	channel, err = amqpConn.Channel()
	if err != nil {
		err = fmt.Errorf("open channel :%v", err)
		return
	}
	close = func() {
		channel.Close()
	}
	return channel, close, nil
}

func InitAmqpConn() (err error) {
	go func() {
		for {
			// 长轮训确保连接不断开
			initAmqpConn()
			time.Sleep(10 * time.Second)
		}
	}()
	return nil
}

func initAmqpConn() {
	amqpConnLock.Lock()
	defer amqpConnLock.Unlock()
	if amqpConn != nil {
		if !amqpConn.IsClosed() {
			return
		}
		log.Println("rabbitMQ connection is closed")
		amqpConn.Close()
		amqpConn = nil
	}

	amqpURL := config.String("amqp_url", "")
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return
	}
	amqpConn = conn
	// log.Info("amqp connection successful")
}
