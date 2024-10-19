package service

import (
	authpb "auth/api/auth/v1"
	"auth/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
	uc     *biz.C1Usecase
	logger *log.Helper
}

func NewAuthService(uc *biz.C1Usecase, logger log.Logger) *AuthService {
	return &AuthService{uc: uc, logger: log.NewHelper(logger)}
}

func (s *AuthService) Login(ctx context.Context, req *authpb.LoginReq) (*authpb.LoginRsp, error) {
	return nil, nil
}

func (s *AuthService) Logout(ctx context.Context, req *authpb.LogoutReq) (*authpb.LogoutRsp, error) {
	return nil, nil
}
