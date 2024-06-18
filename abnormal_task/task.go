package abnormal_task

// 这个应该没有用

type Task struct {
	MissionID      int     `json:"missionId"`       // 异常任务队列任务ID(实际上叫类别更合适)：0-->余额异常  1-->消费异常
	Mission        Mission `json:"mission"`         // 异常任务具体内容
	DuixiangID     int     `json:"duixiangId"`      // 预警对象ID(实际上叫类别更合适)：0-->发短信  1-->发邮件  2-->都发
	DuixiangTarget string  `json:"duixiang_target"` // 预警对象地址(手机号/邮箱地址)
}
type Mission struct {
	TimeStamp           string        `json:"timeStamp"`            // 异常任务时间戳(年月日时分秒)
	Balance             float32       `json:"balance"`              // 余额
	AbnormalConsumption []Consumption `json:"abnormal_consumption"` // 异常业务统计
}
type Consumption struct {
	ConsumptionName   string  `json:"consumption_name"`   // 异常业务名称
	ConsumptionAmount float32 `json:"consumption_amount"` // 异常任务量(？)
}
