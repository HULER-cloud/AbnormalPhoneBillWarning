package config

type Config struct {
	Email  Email  `yaml:"email"`
	System System `yaml:"system"`
	Mysql  Mysql  `yaml:"mysql"`
	Redis  Redis  `yaml:"redis"`
	JWT    JWT    `yaml:"jwt"`
	Expire Expire `yaml:"expire"`
}
