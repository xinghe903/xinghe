package biz

import (
	"comment/internal/biz/repo"
	"comment/internal/conf"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/xinghe903/xinghe/pkg/encrypt"
)

type CommentUsecase struct {
	log   *log.Helper
	cRepo repo.CommentRepo
	enc   *encrypt.EncryptAes
}

func NewCommentUsecase(c *conf.Config, logger log.Logger, u repo.CommentRepo) *CommentUsecase {
	return &CommentUsecase{
		log:   log.NewHelper(logger),
		cRepo: u,
		enc:   encrypt.NewEncryptAes(c.EncryptKey),
	}
}

func (a *CommentUsecase) Logout(ctx context.Context, token string) error {
	return nil
}
