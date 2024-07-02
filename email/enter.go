package email

import (
	"AbnormalPhoneBillWarning/DataAnalysis/RabbitMQ"
	"AbnormalPhoneBillWarning/abnormal_task"
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/models"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
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
	//TaskSend(task)
	return nil
}

func TaskSend(task abnormal_task.Task) {

	// 从json反序列化为异常任务结构体对象
	var userModel models.UserModel
	err := global.DB.Where("id = ?", task.UserID).Take(&userModel).Error
	if err != nil {
		log.Println(err)
		return
	}

	// 余额异常
	if task.MissionID == 0 {
		log.Println("已向", task.Email, "发出余额异常提醒")
		main_text := fmt.Sprintf("您的当前话费余额为%.2f元，已低于%.2f元，请及时充值！", task.Mission.Balance, userModel.BalanceThreshold)
		SendEmail("3285215326@qq.com", "话费预警小帮手", task.Email, BalanceWarning+"提醒", main_text, "nunufvstvjwtdagf", "smtp.qq.com", 465)
	}
	if task.MissionID == 1 {
		log.Println("已向", task.Email, "发出消费异常提醒")
		// 消费异常
		main_text := fmt.Sprintf("检测到您本月消费额度为%.2f元，已高于%.2f元<br>情况存在异常，可能的异常消费如下，请进一步核查：<br><br>", task.Mission.Cost, userModel.BusinessThreshold)
		for _, v := range task.Mission.AbnormalConsumption {
			main_text += v.ConsumptionName + "  " + fmt.Sprintf("%.2f", v.ConsumptionAmount) + "<br>"
		}
		SendEmail("3285215326@qq.com", "话费预警小帮手", task.Email, ConsumptionWarning+"提醒", main_text, "nunufvstvjwtdagf", "smtp.qq.com", 465)
	}

	if err != nil {
		log.Println("邮件发送失败")
	}

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
