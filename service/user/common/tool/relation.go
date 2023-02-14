package tool

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"train-tiktok/common/errorx"
	"train-tiktok/service/user/models"
)

// IsFollowing
// 获取 userId 是否关注 targetId
// if userId == 0 || targetId == 0: return false
// if userId == targetId: return true
func IsFollowing(ctx context.Context, db *gorm.DB, userId int64, targetId int64) (followed bool, err error) {
	if userId == 0 || targetId == 0 {
		return false, nil
	} else if userId == targetId { // prevent self follow
		return true, nil
	}

	if err = db.Model(&models.Follow{}).Where(&models.Follow{
		UserId:   userId,
		TargetId: targetId,
	}).First(&models.Follow{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		followed = false
	} else if err != nil {
		logx.WithContext(ctx).Errorf("failed to query isFollowed: %v", err)

		return false, errorx.ErrDatabaseError
	}

	return true, nil
}
