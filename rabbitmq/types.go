package rabbitmq


import "github.com/streadway/amqp"

type Exchange struct {
	Name       string     // name of the exchange
	Type       string     // type
	Durable    bool       // durable
	AutoDelete bool       // delete when complete
	Internal   bool       // internal
	NoWait     bool       // noWait
	Args       amqp.Table // arguments
}

type Queue struct {
	Name       string     // name of the queue
	Durable    bool       // durable
	AutoDelete bool       // delete when unused
	Exclusive  bool       // exclusive
	NoWait     bool       // noWait
	Args       amqp.Table // arguments
}

type QueueBind struct {
	Name     string     // name of the queue
	Key      string     // bindingKey
	Exchange string     // sourceExchange
	NoWait   bool       // noWait
	Args     amqp.Table // arguments
}

type Consume struct {
	Name      string     // name
	Tag       string     // consumerTag
	AutoAck   bool       // noAck
	Exclusive bool       // exclusive
	NoLocal   bool       // noLocal
	NoWait    bool       // noWait
	Args      amqp.Table // arguments
}
