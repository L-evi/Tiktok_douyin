package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"log"
	"os"
	"train-tiktok/common/dbutil"
	"train-tiktok/common/redisutil"
	"train-tiktok/service/chat/internal/config"
	"train-tiktok/service/chat/models"
)

type ServiceContext struct {
	Config          config.Config
	Db              *gorm.DB
	StrorageBaseUrl config.StorageStruct
	Rdb             *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	logx.MustSetup(c.Log)

	// Gorm
	if dsn, ok := os.LookupEnv("MYSQL_DSN"); ok {
		c.Mysql.DataSource = dsn
	}

	_db, err := dbutil.New(c.Mysql.DataSource, os.Getenv("DEBUG"))
	if err != nil {
		log.Panicf("failed to connect to mysql: %v", err)
	}

	// auto migrate table

	if err := _db.AutoMigrate(&models.Chat{}); err != nil {
		log.Panicf("failed to autoMigrate: %v", err)
	}

	// get LocalBaseUrl
	if local, ok := os.LookupEnv("STORAGE_BASE_URL_LOCAL"); ok {
		c.StorageBaseUrl.Local = local
	}

	// redis
	if rdbAddr, ok := os.LookupEnv("REDIS_ADDR"); ok {
		c.Redis.Addr = rdbAddr
	}
	if rdbPwd, ok := os.LookupEnv("REDIS_PASSWD"); ok {
		c.Redis.Passwd = rdbPwd
	}
	if rdbDb, ok := os.LookupEnv("REDIS_DB"); ok {
		c.Redis.Addr = rdbDb
	}
	if rdbPrefix, ok := os.LookupEnv("REDIS_PREFIX"); ok {
		c.Redis.Prefix = rdbPrefix
	}
	_rdb := redisutil.New(c.Redis)

	// return init
	return &ServiceContext{
		Config: c,
		Db:     _db,
		Rdb:    _rdb,
		StrorageBaseUrl: config.StorageStruct{
			Local: c.StorageBaseUrl.Local,
			Cos:   c.StorageBaseUrl.Cos,
		},
	}
}
