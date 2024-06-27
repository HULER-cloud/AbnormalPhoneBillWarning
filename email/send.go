package email

import (
	"AbnormalPhoneBillWarning/abnormal_task"
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/models"
	"encoding/json"
	"fmt"
	"log"
)

// 以recv开头的三个都没用了，只看send就行

func AbnormalTaskSend(jsonStr string, email_address string) {
	fmt.Println(jsonStr)
	// 从json反序列化为异常任务结构体对象
	task := abnormal_task.Task{}
	err := json.Unmarshal([]byte(jsonStr), &task)
	//fmt.Println(abnormal_task)

	var userModel models.UserModel
	err = global.DB.Where("id = ?", task.UserID).Take(&userModel).Error
	if err != nil {
		log.Println(err)
		return
	}

	// 余额异常
	if task.MissionID == 0 {
		main_text := fmt.Sprintf("您的当前话费余额为%.2f元，已低于%.2f元，请及时充值！", task.Mission.Balance, userModel.BalanceThreshold)
		err = NewBalanceWarning().Send(
			email_address,
			BalanceWarning,
			main_text,
		)
	}
	if task.MissionID == 1 {
		// 消费异常
		main_text := fmt.Sprintf("检测到您本月消费额度为%.2f元，已高于%.2f元<br>情况存在异常，可能的异常消费如下，请进一步核查：<br><br>", task.Mission.Cost, userModel.BusinessThreshold)
		for _, v := range task.Mission.AbnormalConsumption {
			main_text += v.ConsumptionName + "  " + fmt.Sprintf("%.2f", v.ConsumptionAmount) + "<br>"
		}
		err = NewConsumptionWarning().Send(
			email_address,
			ConsumptionWarning,
			main_text,
		)
	}

	if err != nil {
		log.Println("邮件发送失败")
	}

}
