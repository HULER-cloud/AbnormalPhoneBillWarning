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

type BusinessInfo struct {
	BusinessName string    `json:"business_name"`
	Spending     float32   `json:"spending"`
	QueryDate    time.Time `json:"query_date"`
}

func (BusinessAPI) BusinessInfoGetView(c *gin.Context) {
	_claims, _ := c.Get("claims")
	claims := _claims.(*mdw_jwt.MyClaims) // 断言一下，因为拿过来的是any

	var pageInfo models.PageInfo
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		//fmt.Println(err.Error())
		response.FailedWithDetails(response.ArgumentsError, c)
		return
	}

	count, userBusinessList := utils.ListMethod[models.UserBusinessModel](claims.UserID, pageInfo)
	var returnList []BusinessInfo
	//var businessList []models.BusinessModel
	//err=global.DB.Where("id in ?",).Find(&businessList).Error
	//if err != nil {
	//
	//}

	for _, userBusiness := range userBusinessList {
		var businessModel models.BusinessModel
		err = global.DB.Where("id = ?", userBusiness.BusinessID).Take(&businessModel).Error
		if err != nil {
			response.FailedWithMsg("查询用户当前业务失败，请重试！", c)
			return
		}

		resp := BusinessInfo{
			BusinessName: businessModel.BusinessName,
			Spending:     userBusiness.Spending,
			QueryDate:    userBusiness.UpdatedAt,
		}
		// 只有新查询的才加进去
		// 因为会自动执行爬虫脚本来查询，所以不用担心查到过早的数据
		// 当然一定的延迟是不可避免的
		if time.Now().Sub(resp.QueryDate).Hours() < float64(36) {
			returnList = append(returnList, resp)
		}
		//returnList = append(returnList, resp)
	}

	response.OKWithList(count, returnList, c)
}
