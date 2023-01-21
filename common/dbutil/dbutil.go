package dbutil

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func New(dsn string, debug string) (*gorm.DB, error) {

	var logMode = logger.Warn
	if debug == "true" {
		logMode = logger.Info
	}

	sqlLogger := logger.New(
		log.New(os.Stderr, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Millisecond * 200, // 慢 SQL 阈值
			LogLevel:                  logMode,                // 日志级别
			IgnoreRecordNotFoundError: true,                   // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,                   // 禁用彩色打印
		},
	)

	var err error
	D, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: sqlLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDb, err := D.DB()
	sqlDb.SetMaxOpenConns(60)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetConnMaxLifetime(15 * time.Minute)
	if err := sqlDb.Ping(); err != nil {
		return nil, err
	}
	return D, nil
}
