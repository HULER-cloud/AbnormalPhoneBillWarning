package config

type JWT struct {
	Secret     string `yaml:"secret_key" json:"secret_key"`   // 密钥
	ExpireTime int    `yaml:"expire_time" json:"expire_time"` // 过期时间（小时）
	Issuer     string `yaml:"issuer" json:"issuer"`           // 颁发人
}
