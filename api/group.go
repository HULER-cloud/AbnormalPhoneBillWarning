package api

import (
	"AbnormalPhoneBillWarning/api/bussiness_api"
	"AbnormalPhoneBillWarning/api/user_api"
)

type APIGroup struct {
	UserAPI     user_api.UserAPI
	BusinessAPI bussiness_api.BusinessAPI
	//SettingsAPI settings_api.SettingsAPI
	//ImageAPI    image_api.ImageAPI
	//AdAPI       ad_api.AdAPI
	//MenuAPI     menu_api.MenuAPI
	//UserAPI     user_api.UserAPI
	//TagAPI      tag_api.TagAPI
	//MessageAPI  message_api.MessageAPI
	//ArticleAPI  article_api.ArticleAPI
	//CommentAPI  comment_api.CommentAPI
}

// 一个包含所有模块API的大的实例化对象
// 实际上是空的，为了方便调方法
var APIGroupAPP = new(APIGroup)
