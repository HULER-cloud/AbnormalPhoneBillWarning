package utils_spider

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/models"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
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

func TTT() {
	Spider("山东", 1, "17353825020", "153231")
	//jsonStr := "{\n  \"balance\": 56.82,\n  \"timeStamp\": \"2024/06/25 19:12:31\",\n  \"consumption_condition\": {\n    \"consumption_amount\": 3.15,\n    \"consumption_set\": [\n      {\n        \"subscription_name\": \"20元包20G国内流量包\",\n        \"subscription_amount\": 0.0\n      },\n      {\n        \"subscription_name\": \"iFree卡（2016版）\",\n        \"subscription_amount\": 3.0\n      },\n      {\n        \"subscription_name\": \"国内通话费\",\n        \"subscription_amount\": 0.15\n      }\n    ]\n  }\n}"
	//jsonStr := "{\n  \"user_id\": 1,\n  \"balance\": 56.82,\n  \"timeStamp\": \"2024/06/25 19:12:31\",\n  \"consumption_condition\": {\n    \"consumption_amount\": 3.15,\n    \"consumption_set\": [\n      {\n        \"subscription_name\": \"20元包20G国内流量包\",\n        \"subscription_amount\": 0.0\n      },\n      {\n        \"subscription_name\": \"iFree卡（2016版）\",\n        \"subscription_amount\": 3.0\n      },\n      {\n        \"subscription_name\": \"国内通话费\",\n        \"subscription_amount\": 0.15\n      }\n    ]\n  }\n}"
	//JSONProcess([]byte(jsonStr), 1)
}

func Spider(province string, userID uint, phoneNum string, pwd string) {
	var targetFile string
	if province == "山东" {
		targetFile = "./utils/utils_spider/slide.py"
	} else if province == "广东" {
		targetFile = "./utils/utils_spider/slide_gd.py"
	}

	// 一层传参
	cmd := exec.Command("python", "./utils/utils_spider/execute.py", targetFile, strconv.Itoa(int(userID)), phoneNum, pwd)
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
	//fmt.Println(spiderInfo)

	// 更新用户的余额
	// 尝试获取目标用户信息
	//fmt.Println(123)
	var userModel models.UserModel
	count := global.DB.Where("id = ?", userID).
		Take(&userModel).RowsAffected
	if count == 0 {
		return
	}
	//fmt.Println(456)

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
			//fmt.Println("新增", userID, businessModel.ID)
			err = global.DB.Create(&models.UserBusinessModel{
				MODEL:      models.MODEL{},
				UserID:     userID,
				BusinessID: businessModel.ID,
				Spending:   v.SubscriptionAmount,
			}).Error
		} else {
			//fmt.Println("更新", userID, businessModel.ID)
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

func TestFun(province string, userID uint, phoneNum string, pwd string) {
	fmt.Println(province, userID, phoneNum, pwd)
}
