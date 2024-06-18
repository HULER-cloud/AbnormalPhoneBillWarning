package dataanalysis

import (
	"AbnormalPhoneBillWarning/DataAnalysis/RabbitMQ"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

/*
DataAnalysis负责处理爬虫进程返回的json结果，对其中的余额和消费进行分析，找出其中的异常，并提交异常任务给mq

******输入******
来自爬虫进程的json格式数据，通过mq进行传递，格式如下：

	{
	  "id": int
	  "balance": float32, //余额
	  "timeStamp": "年月日 时分秒",
	  "consumption_conditon":
	  {
	    "consumption_amount": float32,
	    "consumption_arr":
	    [
	        "consumption_1": {
	          "consumption_name": "业务名称",
	          "consumption_amount": float32
	        },
	        "consumption_2": {
	          "consumption_name": "业务名称2",
	          "consumption_amount": float32
	        },
	        "consumption_3": ...
	    ]
	  }
	}

对应mq参数：
##使用work模式：queueName=PythonCrawlerResult

****************

******输出******
异常任务的json，格式如下：

	{
	    "missionId":0 or 1, //missionID的值为0或1，代表是余额异常还是消费异常

	    "mission":
	    {
	        "timeStamp":"年月日 时分秒",
	        "balance":float32,
	        "abnormal_consumption" :[
				"consumption_1": {
					"consumption_name": "业务名称",
	          		"consumption_amount": float32
				}
				...
			]

	    },
	    "duixiangId":0 or 1 or 2, //对象id代表发送的对象，0=短信，1=邮箱，2=都发送

	    "duixiang_target":"短信 / 邮箱"

}

对应mq参数：
##使用work模式：queueName=AbnormalMission

****************

业务逻辑：
读取余额阈值和业务总消费阈值；
1、余额小于阈值，异常任务；
2、业务总消费大于阈值：

	1）若有新业务，将所有新业务按照消费从大到小提交至异常任务
	2）若无新业务，将消费top3旧业务提交至异常任务
*/
func DataAnalysis() {
	mq_consumer := RabbitMQ.New_RabbitMQ_Struct("PythonCrawlerResult", "", "")
	defer mq_consumer.Destroy()

	mq_publish := RabbitMQ.New_RabbitMQ_Struct("AbnormalMission", "", "")
	defer mq_publish.Destroy()

	var map_publish_RabbitMQ map[string]*RabbitMQ.RabbitMQ
	map_publish_RabbitMQ["AbnormalMission"] = mq_publish
	var list_publish_name []string
	list_publish_name = append(list_publish_name, "AbnormalMission")
	mq_consumer.Consume_Work(map_publish_RabbitMQ, list_publish_name, handler_DataAnalysis)

}

func handler_DataAnalysis(delivery amqp.Delivery, map_pub map[string]*RabbitMQ.RabbitMQ, pub_name []string) error {
	//读取脚本结果
	var task map[string]interface{}
	err := json.Unmarshal(delivery.Body, &task)
	if err != nil {
		log.Fatalf("Error decoding json: %s", err)
		return err
	}

	/* 使用id查表，得到预设的余额阈值和业务消费阈值 */
	var balanceLimit float32 = 10.0
	// 未使用，但是未来可能会使用
	//var consumptionLimit float32 = 5.0

	//检查key是否存在&&类型是否正确
	var balance float32
	if _, exists := task["balance"]; exists {
		if balance_read, ok := task["balance"].(float32); ok {
			balance = balance_read
		} else {
			balance = float32(balance_read)
		}
	} else {
		return fmt.Errorf("balance not found in json")
	}

	if balance < balanceLimit {

	}

	return nil
}
