package utils_spider

import (
	"AbnormalPhoneBillWarning/DataAnalysis/RabbitMQ"
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/models"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type SpiderInfo struct {
	UserID               uint                 `json:"user_id"`
	Balance              float32              `json:"balance"`
	TimeStamp            string               `json:"timeStamp"`
	ConsumptionCondition ConsumptionCondition `json:"consumption_condition"`
}

type ConsumptionCondition struct {
	ConsumptionAmount float32          `json:"consumption_amount"`
	ConsumptionSet    SubscriptionList `json:"consumption_set"`
}

type SubscriptionList []Subscription

func (s SubscriptionList) Len() int {
	return len(s)
}

func (s SubscriptionList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SubscriptionList) Less(i, j int) bool {
	return s[i].SubscriptionAmount > s[j].SubscriptionAmount
}

type Subscription struct {
	SubscriptionName   string  `json:"subscription_name"`
	SubscriptionAmount float32 `json:"subscription_amount"`
}

type MQMessage1 struct {
	Province string `json:"province"`
	UserID   uint   `json:"userid"`
	PhoneNum string `json:"phone_number"`
	Password string `json:"password"`
}

func Spider(province string, userID uint, phoneNum string, pwd string) {

	Message := MQMessage1{
		Province: province,
		UserID:   userID,
		PhoneNum: phoneNum,
		Password: pwd,
	}

	jsonStr, err := json.Marshal(Message)
	if err != nil {
		log.Println(err)
		return
	}
	mq := RabbitMQ.New_RabbitMQ_Work("TimerAwakeScript")
	mq.Publish_Work_JSON_String(string(jsonStr))

	defer mq.Destroy()

}

func JSONProcess(output []byte, userID uint) {
	var spiderInfo SpiderInfo
	err := json.Unmarshal(output, &spiderInfo)

	if err != nil {
		fmt.Println(err)
	}

	// 更新用户的余额
	// 尝试获取目标用户信息

	var userModel models.UserModel
	count := global.DB.Where("id = ?", userID).
		Take(&userModel).RowsAffected
	if count == 0 {
		return
	}

	// 更新信息入库
	err = global.DB.Model(&userModel).Updates(map[string]any{
		"default_query_time": userModel.DefaultQueryTime,
		"query_time":         userModel.QueryTime,
		"balance_threshold":  userModel.BalanceThreshold,
		"phone":              userModel.Phone,
		"phone_password":     userModel.PhonePassword,
		"province":           userModel.Province,
		"balance":            spiderInfo.Balance,
	}).Error
	if err != nil {
		log.Println("用户信息修改失败！")
		return
	}

	for _, v := range spiderInfo.ConsumptionCondition.ConsumptionSet {
		// 增加业务（如果有）
		var businessModel models.BusinessModel
		count = global.DB.Where("business_name = ?", v.SubscriptionName).
			Take(&businessModel).RowsAffected
		if count == 0 {
			err = global.DB.Create(&models.BusinessModel{
				MODEL:        models.MODEL{},
				BusinessName: v.SubscriptionName,
			}).Error
			if err != nil {
				return
			}
		}

		// 更新用户的历史业务
		global.DB.Where("business_name = ?", v.SubscriptionName).
			Take(&businessModel)
		err = global.DB.Create(&models.UserBusinessHistoryModel{
			MODEL:      models.MODEL{},
			UserID:     userID,
			BusinessID: businessModel.ID,
			Spending:   v.SubscriptionAmount,
			QueryDate:  spiderInfo.TimeStamp,
		}).Error

		// 更新用户的当前业务
		var userBusinessModel models.UserBusinessModel
		count = global.DB.Where("user_id = ? and business_id = ?", userID, businessModel.ID).
			Take(&userBusinessModel).RowsAffected
		// 没有的业务新增
		if count == 0 {
			err = global.DB.Create(&models.UserBusinessModel{
				MODEL:      models.MODEL{},
				UserID:     userID,
				BusinessID: businessModel.ID,
				Spending:   v.SubscriptionAmount,
			}).Error
		} else {
			// 有的业务更新
			err = global.DB.Model(&userBusinessModel).Updates(map[string]any{
				"user_id":     userID,
				"business_id": businessModel.ID,
				"spending":    v.SubscriptionAmount,
			}).Error
		}

	}

	// 最后删除过期的数据
	// 计算五秒钟前的时间
	fiveSecondsAgo := time.Now().Add(-5 * time.Second)
	// 删去过早的业务，太早的业务说明不是在本次查询中爬到的
	err = global.DB.Where("user_id = ? and updated_at < ?", userID, fiveSecondsAgo).Delete(&models.UserBusinessModel{}).Error
	if err != nil {
		log.Println(err)
		return
	}
}
