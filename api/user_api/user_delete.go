package user_api

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/models"
	"AbnormalPhoneBillWarning/routers/response"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 删除用户是只有管理员能删的
// 但是目前不打算做权限，所以先搁置

// UserDeleteView 删除用户
// @Tags 用户管理
// @Summary 删除用户
// @Description 删除用户
// @Param deleteRequest body models.DeleteRequest true  "删除用户的id列表"
// @Router /api/user/delete [delete]
// @Produce json
// @Success 200 {object} response.Response{}
func (UserAPI) userDeleteView(c *gin.Context) {
	var deleteRequest models.DeleteRequest
	// 接收删除的IDList
	err := c.ShouldBindJSON(&deleteRequest)
	if err != nil {
		response.FailedWithDetails(response.ArgumentsError, c)
		return
	}

	var userList []models.UserModel
	count := global.DB.Find(&userList, deleteRequest.IDList).RowsAffected
	if count == 0 {
		response.FailedWithMsg("要删除的用户不存在！", c)
		return
	}
	global.DB.Delete(&userList)
	response.OKWithMsg(fmt.Sprintf("删除%d个用户成功！", count), c)

}
