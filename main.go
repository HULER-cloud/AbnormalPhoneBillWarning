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
	"os"
)

func main() {

	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	defer file.Close()
	// 设置日志输出到文件
	log.SetOutput(file)

	// 初始化配置项
	core.InitConf()
	core.InitGorm()
	core.InitRedis()

	// 判断要不要迁移（新建）表结构
	db := command.ParseCommand()
	if *db == true {
		log.Println("正在初始化数据库表信息...")
		command.MakeMigrations()
		return
	}

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
	err = global.Router.Run(addr)
	if err != nil {
		log.Fatal("服务器启动失败！")
	}

}
