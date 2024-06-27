package main

import (
	dataanalysis "AbnormalPhoneBillWarning/DataAnalysis"
	"AbnormalPhoneBillWarning/command"
	"AbnormalPhoneBillWarning/core"
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
	app.InitDBandTable(context.Background(), global.Redis, global.DB)
	go dataanalysis.DataAnalysis()

	routers.InitRouter()
	addr := global.Config.System.Addr()
	log.Printf("服务器运行在[%s]\n", addr)
	err := global.Router.Run(addr)
	if err != nil {
		log.Fatal("服务器启动失败！")
	}

}
