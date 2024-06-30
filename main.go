package main

import (
	dataanalysis "AbnormalPhoneBillWarning/DataAnalysis"
	"AbnormalPhoneBillWarning/command"
	"AbnormalPhoneBillWarning/core"
	"AbnormalPhoneBillWarning/email"
	"AbnormalPhoneBillWarning/global"
	"AbnormalPhoneBillWarning/internal/app"
	"AbnormalPhoneBillWarning/routers"
	"context"
	"log"
)

func main() {
	core.InitConf()
	core.InitGorm()
	core.InitRedis()

	db, _ := command.ParseCommand()
	if *db == true {
		log.Println("正在初始化数据库表信息...")
		command.MakeMigrations()
		return
	}

	//utils_spider.TTT()
	//go app.InitTimeTable()
	//go app.UpdateDefaultAccessTimer(utils_spider.Spider)
	// 启动定时器
	app.InitDBandTable(context.Background(), global.Redis, global.DB)
	// 启动分析器
	go dataanalysis.DataAnalysis()

	// 多开几个邮件服务，竞争消费队列提高效率
	go email.MultipleSend()
	go email.MultipleSend()
	go email.MultipleSend()

	// 启动服务器主体监听请求
	routers.InitRouter()
	addr := global.Config.System.Addr()
	log.Printf("服务器运行在[%s]\n", addr)
	err := global.Router.Run(addr)
	if err != nil {
		log.Fatal("服务器启动失败！")
	}

}

//func main() {
//	go email.MultipleSend()
//	time.Sleep(time.Second * 5)
//}
