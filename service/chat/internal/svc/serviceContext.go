package svc

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"os"
	"train-tiktok/common/dbutil"
	"train-tiktok/common/logset"
	"train-tiktok/service/chat/internal/config"
	"train-tiktok/service/chat/models"
)

type ServiceContext struct {
	Config config.Config
	Db     *gorm.DB
	Rdb    *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {

	var debug = false
	if isDebug, ok := os.LookupEnv("DEBUG"); ok {
		if isDebug == "true" {
			debug = true
		}
	}
	logset.Handler(debug, c.Log)

	// Gorm
	if dsn, ok := os.LookupEnv("MYSQL_DSN"); ok {
		c.Mysql.DataSource = dsn
	}

	_db, err := dbutil.New(c.Mysql.DataSource, debug)
	if err != nil {
		log.Panicf("failed to connect to mysql: %v", err)
	}

	// auto migrate table

	if err := _db.AutoMigrate(&models.Chat{}); err != nil {
		log.Panicf("failed to autoMigrate: %v", err)
	}

	if etcdEndpoint, ok := os.LookupEnv("ETCD_ENDPOINT"); ok {
		c.RpcServerConf.Etcd.Hosts = []string{etcdEndpoint}
	}

	// return init
	return &ServiceContext{
		Config: c,
		Db:     _db,
	}
}
