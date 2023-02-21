package logic

import (
	"context"
	"gorm.io/gorm"
	"train-tiktok/common/errorx"
	tool2 "train-tiktok/common/tool"
	"train-tiktok/service/user/common/errx"
	"train-tiktok/service/user/common/tool"
	"train-tiktok/service/user/internal/svc"
	"train-tiktok/service/user/models"
	"train-tiktok/service/user/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationActLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRelationActLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationActLogic {
	return &RelationActLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RelationAct 关注/取消关注
// 需要登录
func (l *RelationActLogic) RelationAct(in *user.RelationActReq) (*user.RelationActResp, error) {
	// prevent self follow
	if in.UserId == 0 {
		return &user.RelationActResp{}, errx.ErrLoginRequired
	} else if in.UserId == in.TargetId {
		return &user.RelationActResp{}, errx.ErrSelfFollow
	}

	// check if targetId exists
	if exists, err := tool2.CheckUserExist(l.ctx, l.svcCtx.IdentityRpc, in.TargetId); err != nil {
		logx.WithContext(l.ctx).Errorf("failed to query user: %v", err)

		return nil, errorx.ErrDatabaseError
	} else if !exists {
		return &user.RelationActResp{}, errorx.ErrUserNotFound
	}

	// gorm
	_db := l.svcCtx.Db

	switch in.Action {
	case 1:
		// prevent repeat follow
		if isFollowed, err := tool.IsFollowing(l.ctx, l.svcCtx.Db, in.UserId, in.TargetId); err != nil {
			logx.WithContext(l.ctx).Errorf("failed to query isFollowed: %v", err)

			return &user.RelationActResp{}, errorx.ErrDatabaseError
		} else if isFollowed {
			// 防止重复关注
			return &user.RelationActResp{
				Success: true,
			}, nil
		}

		// 关注
		// models.Fans(targetId, userId) targetId 新增 userId 为粉丝 // uid 是 tid 的粉丝
		// models.Follow(userId, targetId) userId 新增 targetId 为关注对象 // uid 关注 tid
		if err := _db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&models.Follow{
				UserId:   in.UserId,
				TargetId: in.TargetId,
			}).Error; err != nil {
				return err
			}
			if err := tx.Create(&models.Fans{
				UserId:   in.TargetId,
				TargetId: in.UserId,
			}).Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			logx.WithContext(l.ctx).Errorf("failed to create follow: %v", err)

			return &user.RelationActResp{}, errorx.ErrDatabaseError
		}
		break
	case 2:
		// 取消关注
		// models.Fans(targetId, userId) targetId 删除 userId 为粉丝
		// models.Follow(userId, targetId) userId 删除 targetId 为关注对象
		if err := _db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where(&models.Follow{
				UserId:   in.UserId,
				TargetId: in.TargetId,
			}).Delete(&models.Follow{}).Error; err != nil {
				return err
			}
			if err := tx.Where(&models.Fans{
				UserId:   in.TargetId,
				TargetId: in.UserId,
			}).Delete(&models.Fans{}).Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			logx.WithContext(l.ctx).Errorf("failed to delete follow: %v", err)

			return &user.RelationActResp{}, errorx.ErrDatabaseError
		}
		break
	default:
		return &user.RelationActResp{}, errorx.ErrInvalidParameter
	}

	return &user.RelationActResp{Success: true}, nil
}
