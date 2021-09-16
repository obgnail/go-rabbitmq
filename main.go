package main

import (
	"fmt"
	"github.com/obgnail/go-rabbitmq/handler"
	"time"
)

func main() {
	myMsg := handler.NewMyMessage("hyl", "this is content")
	//err := handler.PublishMessageToMQ(myMsg)
	//err := handler.PublishMessageToDelayMQ(myMsg)
	err := handler.PublishSimplifyMessageToDelayMQ(myMsg)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(20 * time.Second)
}
