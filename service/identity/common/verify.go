package common

import (
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"regexp"
	"train-tiktok/common/response"
	"train-tiktok/service/identity/internal/svc"
	"train-tiktok/service/identity/models"
)

// VerifyUsername 验证用户名是否合法
func VerifyUsername(username string) error {
	if len(username) > 32 || len(username) < 5 {
		return ErrInvalidUsername
	}

	if b, err := regexp.MatchString("^[0-9a-zA-Z]{5,30}$", username); err != nil || !b {
		return ErrInvalidUsername
	}

	return nil
}

// VerifyPwd 验证密码是否合法
func VerifyPwd(pwd string) error {
	if len(pwd) > 32 || len(pwd) < 5 {
		return ErrInvalidPassword
	}

	return nil
}

// IsUsernameExists 判断用户名是否存在
func IsUsernameExists(c *svc.ServiceContext, username string) error {
	if err := c.Db.First(models.User{Username: username}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	} else if err != nil {
		logx.Errorf("failed to query user: %v", err)
		return response.ErrDatabaseError
	}
	return ErrUsernameExists
}
