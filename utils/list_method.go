package utils

import (
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/models"
)

func ListMethod[T any](id uint, pi models.PageInfo) (int64, []T) {
	var list []T
	// 统计总共有多少条数据
	count := global.DB.
		Select("id").Find(&list).RowsAffected

	// 默认按时间从新到旧排
	if pi.Sort == "" {
		pi.Sort = "created_at desc"
	}

	if pi.Page == 0 && pi.Limit == 0 {
		// 一次性找出来所有数据
		//fmt.Println(pi)
		global.DB.Where("user_id = ?", id).Order(pi.Sort).Find(&list)
	} else {
		// 只把目标页数的全找出来
		global.DB.Where("user_id = ?", id).Limit(pi.Limit).
			Offset((pi.Page - 1) * pi.Limit).
			Order(pi.Sort).Find(&list)
	}

	//fmt.Println(imageList)
	// 返回结果
	return count, list
}
