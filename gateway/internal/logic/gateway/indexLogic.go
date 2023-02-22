package gateway

import (
	"context"

	"train-tiktok/gateway/internal/svc"
	"train-tiktok/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index() (resp *types.IndexResp, err error) {
	return &types.IndexResp{
		Resp: types.Resp{
			Code: 0,
			Msg:  "ok",
		},
		Github: "https://github.com/L-evi/Tiktok_douyin",
		Author: []string{"L-evi <Levitang@126.com>", "xcsoft <contact@xcsoft.top>"},
	}, nil
}
