package global

import (
	"AbnormalPhoneBillWarning/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Config *config.Config
	DB     *gorm.DB

	Router *gin.Engine
	Redis  *redis.Client
)
