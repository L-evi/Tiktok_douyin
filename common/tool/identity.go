package tool

import (
	"context"
	"train-tiktok/common/errorx"
	"train-tiktok/service/identity/identityclient"
)

func CheckUserExist(ctx context.Context, identity identityclient.Identity, userId int64) (bool, error) {
	if userId == 0 {
		return false, nil
	}
	if _, err := identity.GetUserInfo(ctx, &identityclient.GetUserInfoReq{
		UserId: userId,
	}); errorx.IsRpcError(err, errorx.ErrUserNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
