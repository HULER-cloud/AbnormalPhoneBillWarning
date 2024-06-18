package email

import (
	"AbnormalPhoneBillWarning/abnormal_task"
	"encoding/json"
	"fmt"
	"log"
)

// 以recv开头的三个都没用了，只看send就行

func send(jsonStr string) {
	// json怎么过来的还有待具体实现，这里先写死在这里方便先写后续逻辑
	//jsonStr := `{
	//   "missionId": 1,
	//   "mission":{
	//	   "timeStamp":"2024-06-11 21:58:11",
	//	   "balance": 114.51,
	//	   "abnormal_consumption":[
	//		   {
	//			   "consumption_name": "业务1",
	//			   "consumption_amount": 23.66
	//		   },
	//			{
	//			   "consumption_name": "业务2",
	//			   "consumption_amount": 17.33
	//		   }
	//	   ]
	//   },
	//   "duixiangId": 1,
	//   "duixiang_target": "2799591178@qq.com"
	//}`
	// 从json反序列化为异常任务结构体对象
	task := abnormal_task.Task{}
	err := json.Unmarshal([]byte(jsonStr), &task)
	//fmt.Println(abnormal_task)

	// 发邮件功能的基础信息
	//email_info := EmailSendInfo{
	//	Host:             "smtp.qq.com",
	//	Port:             465,
	//	User:             "2799591178@qq.com", // 写这部分代码的人自己的QQ邮箱，暂且先这么用着
	//	Password:         "kenmohcdqtiydgac",
	//	DefaultFromEmail: "话费异常预警系统",
	//	UseSSL:           true,
	//	UseTLS:           true,
	//}

	// 实际上是要获取的，但是现在还不知道怎样获取过来，先写死
	// 要连数据库查，但是数据库现在还没整好，先mark一下
	balanceLine := 100

	if task.MissionID == 0 {
		main_text := fmt.Sprintf("您的话费余额低于%d元，请及时充值！", balanceLine)
		err = NewBalanceWarning().Send(
			task.DuixiangTarget,
			BalanceWarning,
			main_text,
		)
	} else {
		main_text := "检测到您近期消费情况存在异常，异常消费如下：<br>"
		for _, v := range task.Mission.AbnormalConsumption {
			main_text += v.ConsumptionName + "  " + fmt.Sprintf("%.2f", v.ConsumptionAmount) + "<br>"
		}
		err = NewConsumptionWarning().Send(
			task.DuixiangTarget,
			ConsumptionWarning,
			main_text,
		)
	}

	if err != nil {
		log.Println("邮件发送失败")
	}

}
