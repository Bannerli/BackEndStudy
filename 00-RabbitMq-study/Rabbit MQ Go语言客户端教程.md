# Rabbit MQ Go语言客户端教程

# 1 介绍

RabbitMq是一个消息代理，它接受并转发消息。（接受、存储和转发二进制数据块（消息））

生产（P）：发送

![image-20220410154719026](Rabbit MQ Go语言客户端教程.assets/image-20220410154719026.png)

队列（queue）：消息存储在队列中。队列本质上是一个大的消息缓冲区。许多生产者可以向一个队列发送消息，而许多消费者可以尝试从一个队列接受数据，

![image-20220410154741054](Rabbit MQ Go语言客户端教程.assets/image-20220410154741054.png)

消费与接收具有相似的含义。消费者是一个主要等待接收消息的程序

生产者，消费者和代理（broker）不必位于同一个主机上。实际上，在大多数应用中他们都不是，一个应用程序既可以是生产者又可以是消费者。

# 2 Hello World

![image-20220410155000465](Rabbit MQ Go语言客户端教程.assets/image-20220410155000465.png)

“ P”是我们的生产者，“ C”是我们的消费者。中间的框是一个队列——RabbitMQ代表消费者保存的消息缓冲区

Go RabbitMQ客户端库

RabbitMQ讲多种协议。本教程使用`amqp0-9-1`，这是一个开放的、通用的消息传递协议。RabbitMQ有[许多不同语言的客户端](http://rabbitmq.com/devtools.html)。在本教程中，我们将使用Go amqp客户端。

首先，使用`go get`安装amqp

```bash
go get github.com/streadway/amqp
```

## 2.1 发送消息

![image-20220410155126011](Rabbit MQ Go语言客户端教程.assets/image-20220410155126011.png)

我们将消息发布者（发送者）称为 `send.go`，将消息消费者（接收者）称为`receive.go`。发布者将连接到RabbitMQ，发送一条消息，然后退出。

**send.go**

```go
package main

import (
  "log"

  "github.com/streadway/amqp"
)

// 检查每个amqp调用的返回值
func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}
// 1. 尝试连接RabbitMQ，建立连接
// 该连接抽象了套接字连接，并为我们处理协议版本协商和认证等。
conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
failOnError(err, "Failed to connect to RabbitMQ")
defer conn.Close()

// 2. 接下来，我们创建一个通道，大多数API都是用过该通道操作的。
ch, err := conn.Channel()
failOnError(err, "Failed to open a channel")
defer ch.Close()

// 3. 声明消息要发送到的队列
q, err := ch.QueueDeclare(
  "hello", // name
  false,   // durable
  false,   // delete when unused
  false,   // exclusive
  false,   // no-wait
  nil,     // arguments
)
// 3.1 检测是否出错
failOnError(err, "Failed to declare a queue")

// 4.将消息发布到声明的队列
body := "Hello World!"
err = ch.Publish(
  "",     // exchange
  q.Name, // routing key
  false,  // mandatory
  false,  // immediate
  amqp.Publishing {
    ContentType: "text/plain",
    Body:        []byte(body),
  })
failOnError(err, "Failed to publish a message")
```





## 2.2 接收消息

![image-20220410155712178](Rabbit MQ Go语言客户端教程.assets/image-20220410155712178.png)

**消费者监听来自RabbitMQ的消息**（不是生产者），因此与发布单个消息的发布者不同，我们将使消费者保持运行状态以监听消息并打印出来。

**receive.go**

```go
package main

import (
  "log"

  "github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}

// 建立连接
conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
failOnError(err, "Failed to connect to RabbitMQ")
defer conn.Close()

// 获取channel
ch, err := conn.Channel()
failOnError(err, "Failed to open a channel")
defer ch.Close()

// 声明队列:因为我们可能在发布者之前启动消费者,所以要在消费者里面也声明队列，
// 用于保证使用队列之前中的消息之前确保队列存在
q, err := ch.QueueDeclare(
  "hello", // name
  false,   // durable
  false,   // delete when unused
  false,   // exclusive
  false,   // no-wait
  nil,     // arguments
)
failOnError(err, "Failed to declare a queue")

// 获取接收消息的Delivery通道
msgs, err := ch.Consume(
  q.Name, // queue
  "",     // consumer
  true,   // auto-ack
  false,  // exclusive
  false,  // no-local
  false,  // no-wait
  nil,    // args
)
failOnError(err, "Failed to register a consumer")

forever := make(chan bool)

// 因为服务器是异步的发送消息给消费者，所以使用goroutine从通道中读取消息。
go func() {
  for d := range msgs {
    log.Printf("Received a message: %s", d.Body)
  }
}()

log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
<-forever
```













































































































































































































































