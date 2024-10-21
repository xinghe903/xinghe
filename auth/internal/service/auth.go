package service

import (
	authpb "auth/api/auth/v1"
	"auth/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
	uc     *biz.AuthUsecase
	logger *log.Helper
}

func NewAuthService(uc *biz.AuthUsecase, logger log.Logger) *AuthService {
	return &AuthService{uc: uc, logger: log.NewHelper(logger)}
}

func (s *AuthService) Register(ctx context.Context, req *authpb.RegisterReq) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *AuthService) Login(ctx context.Context, req *authpb.LoginReq) (*authpb.LoginRsp, error) {
	return &authpb.LoginRsp{Token: "123"}, nil
}

func (s *AuthService) Logout(ctx context.Context, req *authpb.LogoutReq) (*emptypb.Empty, error) {
	return nil, nil
}
