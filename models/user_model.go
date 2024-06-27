package models

type UserModel struct {
	MODEL    `json:"model"`
	Email    string `json:"email"`    // 用户账号（邮箱）
	Password string `json:"password"` // 用户密码

	DefaultQueryTime  string  `json:"default_query_time"`
	QueryTime         string  `json:"query_time"`         // 用户每日的查询时间
	Balance           float32 `json:"balance"`            // 用户当前余额
	BalanceThreshold  float32 `json:"balance_threshold"`  // 用户余额阈值
	BusinessThreshold float32 `json:"business_threshold"` // 用户业务阈值

	Phone         string `json:"phone"`          // 用户手机号
	PhonePassword string `json:"phone_password"` // 登运营商用的密码
	Province      string `json:"province"`       // 省份

	BusinessModels []BusinessModel `gorm:"many2many:user_business_models"` // 用于多对多关系表
}
