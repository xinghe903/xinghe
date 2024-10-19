package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type AuthUsecase struct {
	log *log.Helper
}

func NewAuthUsecase(logger log.Logger) *AuthUsecase {
	return &AuthUsecase{c2Repo: c2Repo, log: log.NewHelper(logger)}
}

func (uc *AuthUsecase) Conversations(ctx context.Context) {

}
