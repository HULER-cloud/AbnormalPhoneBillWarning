package email

import (
	"AbnormalPhoneBillWarning/global"
	"gopkg.in/gomail.v2"
)

// 用于发送邮件业务的基本数据结构体
//type EmailSendInfo struct {
//	Host             string `yaml:"host" json:"host"`                             // 邮箱服务器地址
//	Port             int    `yaml:"port" json:"port"`                             // 邮箱服务端口
//	User             string `yaml:"user" json:"user"`                             // 用户
//	Password         string `yaml:"password" json:"password"`                     // 密码
//	DefaultFromEmail string `yaml:"default_from_email" json:"default_from_email"` // 默认发件人名字
//	UseSSL           bool   `yaml:"use_ssl" json:"use_ssl"`                       // 是否使用SSL
//	UseTLS           bool   `yaml:"use_tls" json:"use_tls"`                       // 是否使用TLS
//}

// 邮件主题（邮件类型/邮件标题）
const (
	BalanceWarning     = "余额异常"
	ConsumptionWarning = "消费异常"
	RegisterCode       = "注册验证码"
)

type SendEmailAPI struct {
	Subject string
}

func NewBalanceWarning() SendEmailAPI {
	return SendEmailAPI{Subject: BalanceWarning}
}
func NewConsumptionWarning() SendEmailAPI {
	return SendEmailAPI{Subject: ConsumptionWarning}
}
func NewRegisterCode() SendEmailAPI {
	return SendEmailAPI{Subject: RegisterCode}
}

// Send 调用方式是 email.NewCode().Send(arg1目标邮箱地址, arg2邮件类型, arg3邮件正文, arg4业务基本信息)
func (SendEmailAPI) Send(name, subject, body string) error {
	//e := info
	//fmt.Printf("%v", e)
	return SendEmail(
		global.Config.Email.User,
		global.Config.Email.DefaultFromEmail,
		name,
		subject+"提醒",
		body,
		global.Config.Email.Password,
		global.Config.Email.SendHost,
		global.Config.Email.SendPort)
}

// userName 发件人邮箱
// sendName 发件人昵称
// mailTo 收件人邮箱
// subject 邮件主题
// body 邮件内容
// authCode 邮箱服务器授权码
// host 邮箱服务器地址
// port 邮箱服务器端口

func SendEmail(userName, sendName, mailTo, subject, body string, authCode, host string, port int) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(userName, sendName))
	m.SetHeader("To", mailTo)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(host, port, userName, authCode)
	err := d.DialAndSend(m)
	return err
}
