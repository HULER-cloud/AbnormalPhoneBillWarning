package models

import "time"

// 覆盖（替换）gorm.Model
type MODEL struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 高频使用，分页展示模型
type PageInfo struct {
	Page int `form:"page"` // 第几页
	//Key   string `form:"key"`   // 用来搜索的key
	Limit int    `form:"limit"` // 一页有多少条数据
	Sort  string `form:"sort"`  // 搜索的排序依据
}

// 高频使用，删除请求
type DeleteRequest struct {
	IDList []uint `form:"id_list" json:"id_list"` // 请求删除数据的ID列表
}
