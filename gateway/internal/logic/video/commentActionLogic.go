package video

import (
	"context"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/service/user/types/user"
	"train-tiktok/service/video/types/video"

	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentActionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentActionLogic {

	return &CommentActionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentActionLogic) CommentAction(req *types.CommentActionReq) (resp *types.CommentActionResp, err error) {
	_userId := l.ctx.Value("user_id").(int64)

	// sent to rpc to consult
	rpcResp, err := l.svcCtx.VideoRpc.CommentAction(l.ctx, &video.CommentActionReq{
		VideoId:     req.VideoId,
		ActionType:  req.ActionType,
		CommentText: req.CommentText,
		CommentId:   req.CommentId,
		UserId:      _userId,
	})

	// consult failed
	if err != nil {
		return &types.CommentActionResp{
			Resp: errx.HandleRpcErr(err),
		}, nil
	}

	if req.ActionType == 1 {
		// add comment 时才需要返回 评论内容
		userRpcResp, err := l.svcCtx.UserRpc.User(l.ctx, &user.UserReq{
			UserId:   l.ctx.Value("user_id").(int64),
			TargetId: _userId,
		})

		// failed to get user information
		if err != nil {
			return &types.CommentActionResp{
				Resp: errx.HandleRpcErr(err),
			}, nil
		}

		// consult success
		return &types.CommentActionResp{
			Resp: errx.SUCCESS_RESP,
			Comment: types.Comment{
				Id: rpcResp.Comment.Id,
				User: types.User{
					Id:            _userId,
					Name:          userRpcResp.Name,
					FollowerCount: *userRpcResp.FollowerCount,
					FollowCount:   *userRpcResp.FollowCount,
					IsFollow:      userRpcResp.IsFollow,
				},
				Content:    rpcResp.Comment.Content,
				CreateDate: rpcResp.Comment.Content,
			},
		}, nil
	}

	// consult success
	return &types.CommentActionResp{
		Resp: errx.SUCCESS_RESP,
	}, nil
}
