package command

import "flag"

func ParseCommand() (*bool, *string) {
	db := flag.Bool("db", false, "初始化数据库")
	user := flag.String("u", "", "创建用户")
	flag.Parse()
	return db, user
}

//// 好像没用，先留着
//func IsWebStop(db *bool) bool {
//	if *db == true {
//		return true
//	}
//	return false
//}
//
//func SwitchOption(is_db bool) {
//	if is_db == true {
//		MakeMigrations()
//	}
//}
