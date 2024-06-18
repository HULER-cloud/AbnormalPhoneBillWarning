package config

import "strconv"

type Mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DB       string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	LogLevel string `yaml:"log_level"`
	Config   string `yaml:"config"` // 数据库连接的高级配置
}

func (ms Mysql) Dsn() string {
	return ms.User + ":" + ms.Password + "@tcp(" + ms.Host + ":" + strconv.Itoa(ms.Port) + ")/" + ms.DB + ms.Config
}
