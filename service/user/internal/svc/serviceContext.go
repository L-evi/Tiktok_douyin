package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"log"
	"os"
	"train-tiktok/common/dbutil"
	"train-tiktok/service/identity/identityclient"
	"train-tiktok/service/user/internal/config"
	"train-tiktok/service/user/models"
)

type ServiceContext struct {
	Config      config.Config
	Db          *gorm.DB
	IdentityRpc identityclient.Identity
}

func NewServiceContext(c config.Config) *ServiceContext {
	// release
	var debug = false
	if isDebug, ok := os.LookupEnv("DEBUG"); ok {
		if isDebug == "true" {
			debug = true
			c.Log.Level = "debug"
			c.Log.Mode = "console"
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

	// 自动生成表结构
	if err := _db.AutoMigrate(models.Fans{}); err != nil {
		log.Panicf("failed to autoMigrate: %v", err)
	}
	if err := _db.AutoMigrate(models.Follow{}); err != nil {
		log.Panicf("failed to autoMigrate: %v", err)
	}
	//if err := _db.AutoMigrate(models.UserFavorite{}); err != nil {
	//	log.Panicf("failed to autoMigrate: %v", err)
	//}
	// connect identityRpc
	// set Etcd Host
	if etcdEndpoint, ok := os.LookupEnv("ETCD_ENDPOINT"); ok {
		c.RpcServerConf.Etcd.Hosts = []string{etcdEndpoint}
		c.IdentityRpcConf.Etcd.Hosts = []string{etcdEndpoint}
	}

	return &ServiceContext{
		Config:      c,
		Db:          _db,
		IdentityRpc: identityclient.NewIdentity(zrpc.MustNewClient(c.IdentityRpcConf)),
	}
}
