package video

import (
	"context"
	"train-tiktok/gateway/common/errx"
	"train-tiktok/gateway/common/tool/rpcutil"
	"train-tiktok/service/video/types/video"

	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {

	return &CommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommentListLogic) CommentList(req *types.CommentListReq) (resp *types.CommentListResp, err error) {
	var _userId int64
	if l.ctx.Value("is_login").(bool) {
		_userId = l.ctx.Value("user_id").(int64)
	}
	// sent to rpc to consult
	rpcResp, err := l.svcCtx.VideoRpc.CommentList(l.ctx, &video.CommentListReq{
		VideoId: req.VideoId,
	})
	// consult failed
	if err != nil {
		logx.Errorf("get comment list failed: %v", err)

		return &types.CommentListResp{
			Resp: errx.HandleRpcErr(err),
		}, nil
	}
	// consult success
	var commentList []types.Comment
	for _, v := range rpcResp.CommentList {
		// get user information from user-rpc-service
		userRpcResp, err := rpcutil.GetUserInfo(l.svcCtx, l.ctx, _userId, v.UserId)
		if err != nil {
			return &types.CommentListResp{
				Resp: errx.HandleRpcErr(err),
			}, nil
		}

		commentList = append(commentList, types.Comment{
			Id:         v.Id,
			User:       userRpcResp,
			Content:    v.Content,
			CreateDate: v.CreateDate,
		})
	}

	return &types.CommentListResp{
		Resp:        errx.SUCCESS_RESP,
		CommentList: commentList,
	}, nil
}
