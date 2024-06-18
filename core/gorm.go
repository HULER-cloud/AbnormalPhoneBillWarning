package core

import (
	"AbnormalPhoneBillWarning/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

func InitGorm() {
	// 设置数据库连接的dsn
	dsn := global.Config.Mysql.Dsn()
	// 根据dsn连接数据库，并设置数据库操作的日志系统为自定义logger
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("初始化数据库[%s]失败！%s\n", dsn, err)
	}
	// 获取到数据库对象进行一些设置
	settingDB, _ := db.DB()
	settingDB.SetMaxIdleConns(10)                 // 最大空闲连接数
	settingDB.SetMaxOpenConns(100)                // 最大总连接数
	settingDB.SetConnMaxLifetime(time.Hour * 100) // 单个连接最大持续时间

	// 赋值给全局数据库变量
	global.DB = db
	log.Printf("初始化数据库[%s]成功...\n", dsn)
}
