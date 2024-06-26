package RabbitMQ

import (
	"log"

	"github.com/streadway/amqp"
)

const url = "amqp://guest:guest@localhost:5672/"

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel

	QueueName string
	Exchange  string
	Key       string
	Mqurl     string
}

type JSON_Format interface {
	JSON_Marshal(any) ([]byte, error)
}

func New_RabbitMQ_Struct(queueName string, exchange string, key string) *RabbitMQ {
	return &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: url}
}

func (mq *RabbitMQ) Destroy() {
	mq.Channel.Close()
	mq.Connection.Close()
}

func (mq *RabbitMQ) ErrorCatch(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
	}
}

//目前只使用work模式，对work模式的mq实例和生产者消费者实现

func New_RabbitMQ_Work(queueName string) *RabbitMQ {
	rabbitmq := New_RabbitMQ_Struct(queueName, "", "")
	var err error

	rabbitmq.Connection, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.ErrorCatch(err, "failed to connect RabbitMq")

	rabbitmq.Channel, err = rabbitmq.Connection.Channel()
	rabbitmq.ErrorCatch(err, "failed to open Channel")

	return rabbitmq
}

// any 需要为json格式的结构体
func (mq *RabbitMQ) Publish_Work_JSON(format JSON_Format, message any) {
	_, err := mq.Channel.QueueDeclare(
		mq.QueueName,
		true,  //是否持久化
		false, //是否自动删除
		false, //是否具有排他性
		false, //是否阻塞处理
		nil,   //额外属性
	)
	mq.ErrorCatch(err, "Failed to declare a queue")

	body, err := format.JSON_Marshal(message)
	mq.ErrorCatch(err, "Failed to encode task to JSON")

	err = mq.Channel.Publish(
		mq.Exchange,
		mq.QueueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	mq.ErrorCatch(err, "Failed to publish message")

}

func (mq *RabbitMQ) Publish_Work_JSON_String(json_str string) {
	_, err := mq.Channel.QueueDeclare(
		mq.QueueName,
		true,  //是否持久化
		false, //是否自动删除
		false, //是否具有排他性
		false, //是否阻塞处理
		nil,   //额外属性
	)
	mq.ErrorCatch(err, "Failed to declare a queue")

	body := []byte(json_str)

	err = mq.Channel.Publish(
		mq.Exchange,
		mq.QueueName,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	mq.ErrorCatch(err, "Failed to publish message")

}

// 自定义函数对数据进行处理，函数参数为amqp.Delivery
// 为函数中可能会使用到的publish行为界定参数：
// map_publish_RabbitMQ map[string]*RabbitMQ 存储需要用到的mq；list_publish_name []string 存储mq对应的名称
func (mq *RabbitMQ) Consume_Work(map_publish_RabbitMQ map[string]*RabbitMQ, list_publish_name []string,
	handler func(amqp.Delivery, map[string]*RabbitMQ, []string) error) {
	q, err := mq.Channel.QueueDeclare(
		mq.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	mq.ErrorCatch(err, "failed to declare a queue")

	msgs, err := mq.Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	mq.ErrorCatch(err, "failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			err := handler(d, map_publish_RabbitMQ, list_publish_name)
			mq.ErrorCatch(err, "error from handler")
		}
	}()

	<-forever

}
