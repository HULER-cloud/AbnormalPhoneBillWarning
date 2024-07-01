package dataanalysis

//package main

import (
	"AbnormalPhoneBillWarning/DataAnalysis/RabbitMQ"
	"AbnormalPhoneBillWarning/abnormal_task"
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/models"
	"AbnormalPhoneBillWarning/utils/utils_spider"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sort"
	"time"

	"github.com/streadway/amqp"
)

/*
DataAnalysis负责处理爬虫进程返回的json结果，对其中的余额和消费进行分析，找出其中的异常，并提交异常任务给mq

******输入******
来自爬虫进程的json格式数据，通过mq进行传递，格式如下：

	{
	  "user_id": int
	  "balance": float32, //余额
	  "timeStamp": "年月日 时分秒",
	  "consumption_condition":
	  [
	    "consumption_amount": float32,
	    "consumption_set":
	    [
	        {
	          "subscription_name": "业务名称",
	          "subscription_amount": float32
	        },
	        {
	          "subscription_name": "业务名称2",
	          "subscription_amount": float32
	        },
	    ]
	  ]
	}

对应mq参数：
##使用work模式：queueName=PythonCrawlerResult

****************

******输出******
异常任务的json，格式如下：

	{
		"id":int
	    "missionId":0 or 1, //missionID的值为0或1，代表是余额异常还是消费异常

	    "mission":
	    {
	        "timeStamp":"年月日 时分秒",
	        "balance":float32,
			"cost":float32,
	        "abnormal_consumption" :[
				{
					"consumption_name": "业务名称",
	          		"consumption_amount": float32
				}
				...
			]

	    }

}

对应mq参数：
##使用work模式：queueName=AbnormalMission

****************

业务逻辑：
读取余额阈值和业务总消费阈值；
1、余额小于阈值，异常任务；
2、业务总消费大于阈值：消费top3业务提交至异常任务
*/

// 爬虫返回结构体
type PythonCrawlerResult struct {
	UserID               int                  `json:"user_id"`
	Balance              float32              `json:"balance"`
	TimeStamp            string               `json:"timeStamp"`
	ConsumptionCondition ConsumptionCondition `json:"consumption_condition"`
}

type ConsumptionCondition struct {
	ConsumptionAmount float32         `json:"consumption_amount"`
	ConsumptionArr    ConsumptionList `json:"consumption_set"`
}

type ConsumptionList []Consumption

type Consumption struct {
	SubscriptionName   string  `json:"subscription_name"`
	SubscriptionAmount float32 `json:"subscription_amount"`
}

func (c ConsumptionList) Len() int {
	return len(c)
}

func (c ConsumptionList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c ConsumptionList) Less(i, j int) bool {
	return c[i].SubscriptionAmount > c[j].SubscriptionAmount
}

func DataAnalysis() {
	mq_consumer := RabbitMQ.New_RabbitMQ_Work("PythonCrawlerResult")
	defer mq_consumer.Destroy()

	mq_publish := RabbitMQ.New_RabbitMQ_Work("AbnormalMission")
	defer mq_publish.Destroy()

	map_publish_RabbitMQ := make(map[string]*RabbitMQ.RabbitMQ)
	map_publish_RabbitMQ["AbnormalMission"] = mq_publish
	var list_publish_name []string
	list_publish_name = append(list_publish_name, "AbnormalMission")

	mq_consumer.Consume_Work(map_publish_RabbitMQ, list_publish_name, handler_DataAnalysis)

}

func initdb() {
	// 设置数据库连接的dsn
	dsn := "root:123456@tcp(127.0.0.1:3306)/base_db?charset=utf8&parseTime=true&loc=Local"
	// 根据dsn连接数据库，并设置数据库操作的日志系统为自定义logger
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	fmt.Println(err)
	if err != nil {
		log.Fatalf("初始化数据库[%s]失败！%s\n", dsn, err)
	}
	// 获取到数据库对象进行一些设置
	settingDB, _ := db.DB()
	settingDB.SetMaxIdleConns(10)                 // 最大空闲连接数
	settingDB.SetMaxOpenConns(100)                // 最大总连接数
	settingDB.SetConnMaxLifetime(time.Hour * 100) // 单个连接最大持续时间
	global.DB = db

	//defer settingDB.Close()
}

func handler_DataAnalysis(delivery amqp.Delivery, map_pub map[string]*RabbitMQ.RabbitMQ, pub_name []string) error {

	log.Println("已经接收到爬虫返回的json，开始处理...")
	// 读取脚本结果
	var task utils_spider.SpiderInfo
	err := json.Unmarshal(delivery.Body, &task)
	if err != nil {
		log.Fatalf("Error decoding json: %s", err)
		return err
	}

	// 根据脚本更新数据库
	go utils_spider.JSONProcess(delivery.Body, task.UserID)
	//fmt.Println(string(delivery.Body))
	//fmt.Println(task)
	/* 使用id查表，得到预设的余额阈值和业务消费阈值 */
	/* 查邮箱，创建协程发送 */
	var userModel models.UserModel
	err = global.DB.Where("id = ?", task.UserID).Take(&userModel).Error
	if err != nil {
		log.Println(err)
		return err
	}

	var balanceLimit float32 = userModel.BalanceThreshold
	var consumptionLimit float32 = userModel.BusinessThreshold
	var emailAddress string = userModel.Email

	//fmt.Println(balanceLimit, consumptionLimit, emailAddress)

	consumptionList := task.ConsumptionCondition.ConsumptionSet
	sort.Sort(consumptionList)

	//余额过低报警
	if balance := task.Balance; balance < balanceLimit {

		var am abnormal_task.Task
		am.UserID = task.UserID
		am.Email = emailAddress
		am.MissionID = 0
		am.Mission.Balance = balance
		am.Mission.Cost = 0
		am.Mission.TimeStamp = task.TimeStamp
		//jsonAm, err := json.Marshal(am)
		//if err != nil {
		//	log.Fatalf("Error encoding json: %s", err)
		//}
		fmt.Println("已向", emailAddress, "余额异常警告！")
		go sendEmail(am)
	}

	if consumptionAmount := task.ConsumptionCondition.ConsumptionAmount; consumptionAmount > consumptionLimit {
		var am abnormal_task.Task
		am.UserID = task.UserID
		am.Email = emailAddress
		am.MissionID = 1
		am.Mission.Balance = 0
		am.Mission.Cost = consumptionAmount
		am.Mission.TimeStamp = task.TimeStamp

		cntConsumption := 0
		for _, v := range consumptionList {
			if cntConsumption >= 3 {
				break
			}
			consumption := v
			ac := abnormal_task.Consumption{
				ConsumptionName:   consumption.SubscriptionName,
				ConsumptionAmount: consumption.SubscriptionAmount,
			}
			am.Mission.AbnormalConsumption = append(am.Mission.AbnormalConsumption, ac)
			cntConsumption++
		}

		//jsonAm, err := json.Marshal(am)
		//if err != nil {
		//	log.Fatalf("Error encoding json: %s", err)
		//}
		fmt.Println("已向", emailAddress, "发出消费异常警告！")
		go sendEmail(am)
	}

	return nil
}

func test() {
	task := abnormal_task.Task{
		UserID:    1,
		Email:     "2799591178@qq.com",
		MissionID: 1,
		Mission:   abnormal_task.Mission{},
	}
	sendEmail(task)
}

func sendEmail(task abnormal_task.Task) {
	///* 这里调用你的包里面的发送函数，我不细写了重复 */
	//Message := abnormal_task.Task{
	//	Province: province,
	//	UserID:   userID,
	//	PhoneNum: phoneNum,
	//	Password: pwd,
	//}
	log.Printf("准备执行id=%d的邮件发送，消息已经发送到队列\n", task.UserID)
	jsonStr, err := json.Marshal(task)
	if err != nil {
		log.Println(err)
		return
	}
	mq := RabbitMQ.New_RabbitMQ_Work("AbnormalMission")
	mq.Publish_Work_JSON_String(string(jsonStr))

	defer mq.Destroy()
	//email.AbnormalTaskSend(jsonStr, email_address)
}
