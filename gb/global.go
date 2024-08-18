package gb

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"vpn-web.funcworks.net/config"
	"xorm.io/xorm"
)

// 日志
var Logger *zap.SugaredLogger

// 配置读取
var Viper = viper.New()

// 配置结构
var Config config.StructConfig

// db
var DB *xorm.Engine

// redis
var RedisProxy redisProxy
var RedisClient redis.UniversalClient

// 定时任务
var Sched SchedManager
