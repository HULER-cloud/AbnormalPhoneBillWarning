package command

import "flag"

func ParseCommand() *bool {
	db := flag.Bool("db", false, "初始化数据库")
	flag.Parse()
	return db
}
