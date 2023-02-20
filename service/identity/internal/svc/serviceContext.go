package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"log"
	"os"
	"train-tiktok/common/dbutil"
	"train-tiktok/service/identity/internal/config"
	"train-tiktok/service/identity/models"
)

type ServiceContext struct {
	Config config.Config
	Db     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	// debug mode
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
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = c.Mysql.DataSource
	}

	_db, err := dbutil.New(dsn, debug)
	if err != nil {
		log.Panicf("failed to connect to mysql: %v", err)
	}
	// 自动生成表结构
	if err := _db.AutoMigrate(models.User{}); err != nil {
		log.Panicf("failed to autoMigrate: %v", err)
	}
	if err := _db.AutoMigrate(models.UserInformation{}); err != nil {
		log.Panicf("failed to autoMigrate: %v", err)
	}

	// Get Jwt SigningKey
	if _jwtSigningKey, ok := os.LookupEnv("JWT_SIGNING_KEY"); ok {
		c.Jwt.SigningKey = _jwtSigningKey
	}

	// Get Conf
	if gravatarBaseUrl, ok := os.LookupEnv("GRAVATAR_BASE_URL"); ok {
		c.Conf.GravatarBaseURL = gravatarBaseUrl
	}
	if _defaultBackground, ok := os.LookupEnv("DEFAULT_BACKGROUND_IMAGE"); ok {
		c.Conf.DefaultBackground = _defaultBackground
	}

	// etcd Endpoints
	if etcdEndpoint := os.Getenv("ETCD_ENDPOINT"); etcdEndpoint != "" {
		c.RpcServerConf.Etcd.Hosts = []string{etcdEndpoint}
	}

	return &ServiceContext{
		Config: c,
		Db:     _db,
	}
}
