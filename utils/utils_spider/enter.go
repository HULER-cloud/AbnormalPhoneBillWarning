package main

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/models"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

type SpiderInfo struct {
	UserID               uint                 `json:"user_id"`
	Balance              float32              `json:"balance"`
	TimeStamp            string               `json:"timeStamp"`
	ConsumptionCondition ConsumptionCondition `json:"consumption_condition"`
}

type ConsumptionCondition struct {
	ConsumptionAmount float32        `json:"consumption_amount"`
	ConsumptionSet    []Subscription `json:"consumption_set"`
}

type Subscription struct {
	SubscriptionName   string  `json:"subscription_name"`
	SubscriptionAmount float32 `json:"subscription_amount"`
}

func main() {
	Spider(1, "山东", "17353825020", "153231")
	//jsonStr := "{\n  \"balance\": 56.82,\n  \"timeStamp\": \"2024/06/25 19:12:31\",\n  \"consumption_condition\": {\n    \"consumption_amount\": 3.15,\n    \"consumption_set\": [\n      {\n        \"subscription_name\": \"20元包20G国内流量包\",\n        \"subscription_amount\": 0.0\n      },\n      {\n        \"subscription_name\": \"iFree卡（2016版）\",\n        \"subscription_amount\": 3.0\n      },\n      {\n        \"subscription_name\": \"国内通话费\",\n        \"subscription_amount\": 0.15\n      }\n    ]\n  }\n}"
	//JSONProcess([]byte(jsonStr), 1)
}

func testfun(phoneNum string, pwd string) {
	// 一层传参
	cmd := exec.Command("python", "execute.py", phoneNum, pwd)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 打印爬取结果
	fmt.Println(string(output))
}

func Spider(userID uint, province string, phoneNum string, pwd string) {
	var targetFile string
	if province == "山东" {
		targetFile = "slide.py"
	} else if province == "广东" {
		targetFile = "slide_gd.py"
	}

	// 一层传参
	cmd := exec.Command("python", "execute.py", strconv.Itoa(int(userID)), targetFile, phoneNum, pwd)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 打印爬取结果
	fmt.Println(string(output))

	//JSONProcess(output, userID)

}

func JSONProcess(output []byte, userID uint) {
	var spiderInfo SpiderInfo
	err := json.Unmarshal(output, &spiderInfo)

	//spiderInfo, err := JSONProcess(string(output))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(spiderInfo)

	// 更新用户的余额
	// 尝试获取目标用户信息
	fmt.Println(123)
	var userModel models.UserModel
	count := global.DB.Where("id = ?", userID).
		Take(&userModel).RowsAffected
	if count == 0 {
		return
	}
	fmt.Println(456)

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
			Take(&businessModel).RowsAffected
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
}
