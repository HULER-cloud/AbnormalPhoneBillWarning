package email

import (
	"AbnormalPhoneBillWarning/global"
	"gopkg.in/gomail.v2"
)

// 邮件主题（邮件类型/邮件标题）
const (
	BalanceWarning     = "余额异常"
	ConsumptionWarning = "消费异常"
	RegisterCode       = "注册验证码"
	ResetCode          = "密码重置验证码"
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
func NewResetCode() SendEmailAPI {
	return SendEmailAPI{Subject: ResetCode}
}

// Send 调用方式是 email.NewCode().Send(arg1目标邮箱地址, arg2邮件类型, arg3邮件正文, arg4业务基本信息)
func (SendEmailAPI) Send(name, subject, body string) error {
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
