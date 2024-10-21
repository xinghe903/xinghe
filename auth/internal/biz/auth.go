package biz

import (
	"auth/internal/biz/po"
	"auth/internal/biz/repo"
	"auth/internal/conf"
	"context"

	authpb "auth/api/auth/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/xinghe903/xinghe/pkg/encrypt"
)

type AuthUsecase struct {
	log   *log.Helper
	uRepo repo.UserRepo
	enc   *encrypt.EncryptAes
}

func NewAuthUsecase(c *conf.Config, logger log.Logger, u repo.UserRepo) *AuthUsecase {
	return &AuthUsecase{
		log:   log.NewHelper(logger),
		uRepo: u,
		enc:   encrypt.NewEncryptAes(c.EncryptKey),
	}
}

func (a *AuthUsecase) Register(ctx context.Context, info *po.User) (string, error) {
	users, err := a.uRepo.ListUser(ctx, &po.PageQuery[po.User]{Condition: &po.User{Name: info.Name}})
	if err != nil {
		return "", err
	}
	if len(users.Data) > 0 {
		return "", authpb.ErrorUsernameRepeat("用户名重复", info.Name)
	}
	info.Password, err = a.enc.Encrypt(info.Password)
	id, err := a.uRepo.CreateUser(ctx, info)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (a *AuthUsecase) Login(ctx context.Context, u *po.User) (string, error) {
	users, err := a.uRepo.ListUser(ctx, &po.PageQuery[po.User]{Condition: &po.User{Name: u.Name}})
	if err != nil || len(users.Data) == 0 {
		return "", authpb.ErrorUserOrPasswordInvalid("用户名或密码错误")
	}
	user := users.Data[0]
	pdText, err := a.enc.Decrypt(user.Password)
	if err != nil || pdText != u.Password {
		return "", authpb.ErrorUserOrPasswordInvalid("用户名或密码错误")
	}

	return "", nil
}

func (a *AuthUsecase) Logout(ctx context.Context, token string) error {
	return nil
}
