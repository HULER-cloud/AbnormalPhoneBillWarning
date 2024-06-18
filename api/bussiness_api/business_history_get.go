package bussiness_api

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/middleware/mdw_jwt"
	"AbnormalPhoneBillWarning/models"
	"AbnormalPhoneBillWarning/routers/response"
	"AbnormalPhoneBillWarning/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type BusinessHistory struct {
	QueryDate    time.Time `json:"query_date"`
	BusinessName string    `json:"business_name"`
	Spending     float32   `json:"spending"`
}

func (BusinessAPI) BusinessHistoryGetView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*mdw_jwt.MyClaims) // 断言一下，因为拿过来的是any

	var pageInfo models.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		//fmt.Println(err.Error())
		response.FailedWithDetails(response.ArgumentsError, c)
		return
	}

	count, userBusinessHistoryList := utils.ListMethod[models.UserBusinessHistoryModel](claims.UserID, pageInfo)
	var returnList []BusinessHistory

	for _, userBusinessHistory := range userBusinessHistoryList {
		var businessModel models.BusinessModel
		err = global.DB.Where("id = ?", userBusinessHistory.BusinessID).Take(&businessModel).Error
		if err != nil {
			response.FailedWithMsg("查询用户当前业务失败，请重试！", c)
			return
		}

		resp := BusinessHistory{
			QueryDate:    userBusinessHistory.QueryDate,
			BusinessName: businessModel.BusinessName,
			Spending:     userBusinessHistory.Spending,
		}
		returnList = append(returnList, resp)
	}

	response.OKWithList(count, returnList, c)
}
