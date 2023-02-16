package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"log"
	"os"
	"train-tiktok/common/dbutil"
	"train-tiktok/common/redisutil"
	"train-tiktok/service/identity/identityclient"
	"train-tiktok/service/video/internal/config"
	"train-tiktok/service/video/models"
)

type ServiceContext struct {
	Config         config.Config
	Db             *gorm.DB
	StorageBaseUrl config.StorageStruct
	Rdb            *redis.Client
	IdentityRpc    identityclient.Identity
}

func NewServiceContext(c config.Config) *ServiceContext {

	// release
	var debug = false
	if isDebug, ok := os.LookupEnv("DEBUG"); ok {
		if isDebug == "true" {
			debug = true
			c.Log.Level = "debug"
		} else {
			c.Log.Level = "info"
			c.Log.Mode = "file"
			c.Log.KeepDays = 60
			c.Log.Rotation = "daily"
			c.Log.Encoding = "json"
		}
	}
	logx.MustSetup(c.Log)

	// Gorm
	if dsn, ok := os.LookupEnv("MYSQL_DSN"); ok {
		c.Mysql.DataSource = dsn
	}

	_db, err := dbutil.New(c.Mysql.DataSource, debug)
	if err != nil {
		log.Panicf("failed to connect to mysql: %v", err)
	}

	//	自动生成表结构
	if err := _db.AutoMigrate(models.Video{}); err != nil {
		log.Panicf("failed to autoMigrate: %v", err)
	}
	if err := _db.AutoMigrate(models.Comment{}); err != nil {
		log.Panicf("failed to autoMigrate: %v", err)
	}

	// 读取 LocalBaseUrl
	if local, ok := os.LookupEnv("STORAGE_BASE_URL_LOCAL"); ok {
		c.StorageBaseUrl.Local = local
	}

	// redis
	if rdbAddr, ok := os.LookupEnv("REDIS_ADDR"); ok {
		c.RedisConf.Addr = rdbAddr
	}
	if rdbPwd, ok := os.LookupEnv("REDIS_PASSWD"); ok {
		c.RedisConf.Passwd = rdbPwd
	}
	if rdbDb, ok := os.LookupEnv("REDIS_DB"); ok {
		c.RedisConf.Addr = rdbDb
	}
	if rdbPrefix, ok := os.LookupEnv("REDIS_PREFIX"); ok {
		c.RedisConf.Prefix = rdbPrefix
	}

	_rdb := redisutil.New(c.RedisConf)

	// set Etcd Host
	if etcdEndpoint, ok := os.LookupEnv("ETCD_ENDPOINT"); ok {
		c.Etcd.Hosts = []string{etcdEndpoint}
		c.IdentityRpcConf.Etcd.Hosts = []string{etcdEndpoint}
	}

	return &ServiceContext{
		Config: c,
		Db:     _db,
		StorageBaseUrl: config.StorageStruct{
			Local: c.StorageBaseUrl.Local,
			Cos:   c.StorageBaseUrl.Cos,
		},
		Rdb:         _rdb,
		IdentityRpc: identityclient.NewIdentity(zrpc.MustNewClient(c.IdentityRpcConf))}
}
