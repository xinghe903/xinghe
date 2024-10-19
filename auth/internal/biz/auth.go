package biz

import (
	"auth/internal/biz/po"
	"auth/internal/biz/repo"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type AuthUsecase struct {
	log   *log.Helper
	uRepo repo.UserRepo
}

func NewAuthUsecase(logger log.Logger, u repo.UserRepo) *AuthUsecase {
	return &AuthUsecase{
		log:   log.NewHelper(logger),
		uRepo: u,
	}
}

func (a *AuthUsecase) Register(ctx context.Context, info *po.User) (string, error) {
	id, err := a.uRepo.CreateUser(ctx, info)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (a *AuthUsecase) Login(ctx context.Context, info *po.User) (string, error) {
	return nil, nil
}

func (a *AuthUsecase) Logout(ctx context.Context, token string) error {
	return nil, nil
}
