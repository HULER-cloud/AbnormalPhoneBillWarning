package email

import (
	"AbnormalPhoneBillWarning/DataAnalysis/RabbitMQ"
	"AbnormalPhoneBillWarning/abnormal_task"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

func MultipleSend() {
	mq_consumer := RabbitMQ.New_RabbitMQ_Work("AbnormalMission")
	defer mq_consumer.Destroy()

	var map_publish_RabbitMQ map[string]*RabbitMQ.RabbitMQ
	var list_publish_name []string

	mq_consumer.Consume_Work(map_publish_RabbitMQ, list_publish_name, handler_MultipleSend)
}

func handler_MultipleSend(delivery amqp.Delivery, map_pub map[string]*RabbitMQ.RabbitMQ, pub_name []string) error {

	var task abnormal_task.Task
	err := json.Unmarshal(delivery.Body, &task)
	if err != nil {
		log.Fatalf("Error decoding json: %s", err)
		return err
	}
	AbnormalTaskSend(task)
	return nil
}
