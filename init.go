package main

import (
	"fmt"
	"github.com/obgnail/go-rabbitmq/amqp"
	"github.com/obgnail/go-rabbitmq/config"
	"github.com/obgnail/go-rabbitmq/conn"
)

func init() {
	onStart(config.LoadConfigs)
	onStart(conn.InitAmqpConn)
	onStart(amqp.InitAmqp)
}

func onStart(fn func() error) {
	if err := fn(); err != nil {
		panic(fmt.Sprintf("Error at onStart: %s\n", err))
	}
}
