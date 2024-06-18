package main

import (
	"AbnormalPhoneBillWarning/command"
	"AbnormalPhoneBillWarning/core"
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/routers"
	"log"
)

func main() {
	core.InitConf()
	core.InitGorm()
	core.InitRedis()

	db, user := command.ParseCommand()
	if *db == true {
		log.Println("正在初始化数据库表信息...")
		command.MakeMigrations()
		return
	}
	if *user == "admin" || *user == "user" {
		log.Println("正在创建用户：", *user)
		command.CreateUser(*user)
		return
	}

	routers.InitRouter()
	addr := global.Config.System.Addr()
	log.Printf("服务器运行在[%s]\n", addr)
	err := global.Router.Run(addr)
	if err != nil {
		log.Fatal("服务器启动失败！")
	}

	// 接收邮件的内容暂时没想好该塞在哪，先mark一下

	//// 获取到邮箱模块的基础信息（包括发送和接收）
	//Cfg = InitConf()
	//fmt.Printf("%v", Cfg)
	//
	//// 然后启动接收邮件的任务，这块在启动后应当是一直运行的
	//email.Recv(Cfg)

}
