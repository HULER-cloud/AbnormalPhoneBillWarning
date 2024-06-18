package core

import (
	"AbnormalPhoneBillWarning/config"
	"AbnormalPhoneBillWarning/global"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

const YamlPath = "./settings.yaml"

// 读取yaml配置文件并初始化
func InitConf() {

	c := &config.Config{}
	// 读出配置文件为字节流
	confBytes, err := os.ReadFile(YamlPath)
	if err != nil {
		log.Fatalf("读取配置文件失败：%s\n", err)
	}
	// 反序列化为配置项结构体
	err = yaml.Unmarshal(confBytes, c)
	if err != nil {
		log.Fatalf("初始化配置项失败！%s\n", err)
	}

	// 赋值给全局配置项变量
	global.Config = c
	log.Println("初始化配置项成功...")
}

// 更新yaml配置文件
func SetConf() error {
	conf, err := yaml.Marshal(global.Config)
	if err != nil {
		return err
	}
	err = os.WriteFile(YamlPath, conf, 0777)
	if err != nil {
		return err
	}
	return nil
}
