// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package server

import (
	"context"

	"train-tiktok/service/user/internal/logic"
	"train-tiktok/service/user/internal/svc"
	"train-tiktok/service/user/types/user"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	user.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServer) User(ctx context.Context, in *user.UserReq) (*user.UserResp, error) {
	l := logic.NewUserLogic(ctx, s.svcCtx)
	return l.User(in)
}

func (s *UserServer) RelationAct(ctx context.Context, in *user.RelationActReq) (*user.RelationActResp, error) {
	l := logic.NewRelationActLogic(ctx, s.svcCtx)
	return l.RelationAct(in)
}
