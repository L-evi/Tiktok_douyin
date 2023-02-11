package svc

import (
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"log"
	"os"
	"train-tiktok/common/dbutil"
	"train-tiktok/common/redisutil"
	"train-tiktok/service/video/internal/config"
	"train-tiktok/service/video/models"
)

type ServiceContext struct {
	Config         config.Config
	Db             *gorm.DB
	StorageBaseUrl config.StorageStruct
	Rdb            *redis.Client
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

	//	自动生成表结构
	if err := _db.AutoMigrate(models.Video{}); err != nil {
		log.Panicf("failed to autoMigrate: %v", err)
	}

	// 读取 LocalBaseUrl
	if local, ok := os.LookupEnv("STORAGE_BASE_URL_LOCAL"); ok {
		c.StorageBaseUrl.Local = local
	}
	if err := _db.AutoMigrate(models.Comment{}); err != nil {
		log.Panicf("failed to autoMigrate: %v", err)
	}

	// redis
	_rdb := redisutil.New(redisutil.RedisConf{
		Addr:        c.Redis.Addr,
		Password:    c.Redis.Password,
		DB:          c.Redis.DB,
		MinIdle:     c.Redis.MinIdle,
		PoolSize:    c.Redis.PoolSize,
		MaxLifeTime: c.Redis.MaxLifeTime,
	})

	return &ServiceContext{
		Config: c,
		Db:     _db,
		StorageBaseUrl: config.StorageStruct{
			Local: c.StorageBaseUrl.Local,
			Cos:   c.StorageBaseUrl.Cos,
		},
		Rdb: _rdb,
	}
}
