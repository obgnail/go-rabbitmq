# 队列持久化

* durable: 持久化, 每个队列和交换机的durable属性。该属性默认情况是false，它决定了RabbitMQ是否需要在崩溃或者重启之后重新创建队列（或者交换机）,将交换机和队列的durable属性设置为true，这样你就不需要在服务器断电后重新创建队列和交换机了

# 消息持久化

投递消息时, 也有个属性: Delivery mode , 可以为 Persistent 和 Transient
* Persistent: 只要消息一达到队列, 就会立即写到文件来持久化
* Transient: 不持久化.

四种组织情况及实际效果

* durable queue , durable message 队列持久化及消息持久化
* durable queue, non-persistent message 队列持久化, 消息非持久化.
* transient queue, non-persistent message 队列非持久化, 消息非持久化
* transient queue, persistent message 队列非持久化, 消息非持久化

# 交换器类型

* direct: 如果 routekey匹配的话, 消息就会被投递到对应的队列. 当声明一个队列的时候, 它会自动绑定到默认交换器,并以队列名称作为路由键.
* fanout: 这种类型的交换器会将收到的消息广播到绑定的队列上. 使用场景: 比如 充值完一笔话费之后, 需要对话费账号增加余额,同时该账号应该得到某些积分奖励的时,可以将两个队列绑定到充值完一笔话费的交换器上. 一个用于增加余额,一个用于增加积分, 后续如果有更多的其他需求,只需增加新的队列并绑定到充值完一笔话费的交换器上就可以了.
* topic: 能够使得来自不同源头的消息到达同一个队列.(routekey 有通配符, 单个 "." 把路由键分为几个部分, "*" 匹配特定位置的任意文本, "#" 匹配所有规则)
* headers: 允许你匹配AMQP消息的header 而非路由键. 除此之外, headers 交换器和direct 交换完全一致, 但性能会差很多(不太实用,不建议使用)

# 队列设置参数

* internal: 是否是内部专用exchange.
* nowait: 如果设置为true，服务器将不会对方法作出回应。客户端只能在非事务性通道上使用此方法.
* exclusive: 如果设置为true, 队列将变为私有的, 此时只有你的应用程序才能够消费队列消息. 使用场景:想要限制一个队列只有一个消费者的情况.
* auto-delete: 如果设置为true,当最后一个消费者取消订阅的时候,队列就会自动删除, 使用场景: 如果你需要临时队列只为一个消费者服务的时候, 结合 auto-delete 和 exclusive, 当消费者断开连接时,队列就被移除了
* args： 是AMQP协议留给AMQP实现做扩展使用的。

# 消息设置参数

## mandatory

* true: 当消息无法通过交换器匹配到队列时，会调用Channel.NotifyReturn通知生产者
* false时, 当消息无法通过交换器匹配到队列时，会丢弃消息.

注意: 不建议设置为 true，因会使程序逻辑变得复杂, 可以通过备用交换机来实现类似的功能。

通过设置参数，可以设置Ex的备用交换器ErrEx
创建Exchange时，指定Ex的Args:
* "alternate-exchange":"ErrEx" 
其中ErrEx为备用交换器名称.

## immediate

* true: 当消息到达Queue后，发现队列上无消费者时，通过Channel.NotifyReturn返回给生产者。 
* false时,消息一直缓存在队列中，等待生产者。 

注意: 不建议设置为true，遇到这种情况，可用TTL和DLX方法代替

通过设置参数，可以设置Queue的DLX交换机DlxEX
创建Queue时，指定Q的Args参数: 
* "x-message-ttl":0 //msg超时时间，单位毫秒 
* "x-dead-letter-exchange":"dlxExchange" //DlxEx名称 
* "x-dead-letter-routing-key":"dlxQueue" //DlxEx路由键


# 死信交换机(Dead Letter Exchanges)

有三种情况可能进死信交换机

1. 被reject或者nack，并且requeue设置为false
2. 消息最大存活时间（TTL）超时
3. 消息数量超过最大队列长度

参数:
x-dead-letter-exchange 将死信消息发送到指定的 exchange 中
x-dead-letter-routing-key 将死信消息发送到自定的 route当中

# 队列限制

RabbitMQ有两种对队列长度的限制方式:

* 对队列中消息的条数进行限制  x-max-length
* 对队列中消息的总量进行限制  x-max-length-bytes

# 其他参数

* x-message-ttl: msg超时时间，单位毫秒