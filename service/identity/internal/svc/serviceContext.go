package svc

import (
	"gorm.io/gorm"
	"log"
	"os"
	"train-tiktok/common/dbutil"
	"train-tiktok/service/identity/internal/config"
	"train-tiktok/service/identity/models"
)

type ServiceContext struct {
	Config        config.Config
	Db            *gorm.DB
	JwtSigningKey []byte
}

func NewServiceContext(c config.Config) *ServiceContext {

	// Gorm
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = c.Mysql.DataSource
	}

	_db, err := dbutil.New(dsn, os.Getenv("DEBUG"))
	if err != nil {
		log.Panicf("failed to connect to mysql: %v", err)
	}
	// 自动生成表结构
	if err := _db.AutoMigrate(models.User{}); err != nil {
		log.Panicf("failed to autoMigrate: %v", err)
	}

	// Get Jwt SigningKey
	_jwtSigningKey := os.Getenv("JWT_SIGNING_KEY")
	if _jwtSigningKey == "" {
		_jwtSigningKey = c.Jwt.SigningKey
	}

	return &ServiceContext{
		Config:        c,
		Db:            _db,
		JwtSigningKey: []byte(_jwtSigningKey),
	}
}
