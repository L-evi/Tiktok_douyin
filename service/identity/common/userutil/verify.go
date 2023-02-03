package userutil

import (
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"regexp"
	"train-tiktok/common/errorx"
	"train-tiktok/service/identity/common/errx"
	"train-tiktok/service/identity/internal/svc"
	"train-tiktok/service/identity/models"
)

// VerifyUsername 验证用户名是否合法
func VerifyUsername(username string) error {
	if len(username) > 32 {
		return errx.ErrInvalidUsername
	}

	// username must be email
	if b, err := regexp.MatchString(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, username); err != nil || !b {
		return errx.ErrInvalidUsername
	}

	return nil
}

// VerifyPwd 验证密码是否合法
func VerifyPwd(pwd string) error {
	if len(pwd) > 32 || len(pwd) < 5 {
		return errx.ErrInvalidPassword
	}

	return nil
}

// IsUsernameExists 判断用户名是否存在
func IsUsernameExists(c *svc.ServiceContext, username string) error {
	if err := c.Db.Where(&models.User{Username: username}).First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	} else if err != nil {
		logx.Errorf("failed to query user: %v", err)
		return errorx.ErrDatabaseError
	}
	return errx.ErrUsernameExists
}
