package config

type Email struct {
	SendHost         string `yaml:"send_host" json:"send_host"` // 邮箱服务器地址
	SendPort         int    `yaml:"send_port" json:"send_port"`
	RecvHost         string `yaml:"recv_host" json:"recv_host"`                   // 邮箱服务器地址
	RecvPort         int    `yaml:"recv_port" json:"recv_port"`                   // 邮箱服务端口
	User             string `yaml:"user" json:"user"`                             // 用户
	Password         string `yaml:"password" json:"password"`                     // 密码
	DefaultFromEmail string `yaml:"default_from_email" json:"default_from_email"` // 默认发件人名字
	UseSSL           bool   `yaml:"use_ssl" json:"use_ssl"`                       // 是否使用SSL
	UseTLS           bool   `yaml:"use_tls" json:"use_tls"`                       // 是否使用TLS
}
