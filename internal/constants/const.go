// 定义一些常量
package constants

import (
	"time"
)

// 两次查询之间的间隔
const QueryInterval = 30 * time.Minute

// 一次查询所需的时间的上界，暂时设为0，因为没想好怎么处理查询数据库很快导致临界用户查询请求重复发送的问题
const QueryTime = 0 * time.Second

const UpdateTimeTableInterval = 24 * time.Hour

const DefaultExpireInterval = 1 * UpdateTimeTableInterval
