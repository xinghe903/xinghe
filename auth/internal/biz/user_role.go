package biz

import (
	authpb "auth/api/auth/v1"
	"auth/internal/biz/repo"
	"auth/internal/conf"
	"auth/internal/data/po"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/xinghe903/xinghe/pkg/bo"
)

type UserRoleUsecase struct {
	log *log.Helper
	ur  repo.UserRoleRepo
}

func NewUserRoleUsecase(c *conf.Config, logger log.Logger,
	r repo.UserRoleRepo) *UserRoleUsecase {
	return &UserRoleUsecase{
		log: log.NewHelper(logger),
		ur:  r,
	}
}

func (a *UserRoleUsecase) ListUserRole(ctx context.Context, userId string) ([]string, error) {
	list, err := a.ur.List(ctx, &bo.PageQuery[po.UserRole]{
		Condition: &po.UserRole{UserId: userId},
	})
	if err != nil {
		a.log.WithContext(ctx).Errorf("list user role: %v", err.Error())
		return nil, authpb.ErrorListUserrole("获取用户角色列表失败")
	}
	var roleIds []string
	for _, rp := range list.Data {
		roleIds = append(roleIds, rp.RoleId)
	}
	return roleIds, nil
}

func (a *UserRoleUsecase) UpdateUserRole(ctx context.Context, userId string, roleIds []string) error {
	err := a.ur.CoverRelations(ctx, userId, roleIds)
	if err != nil {
		a.log.WithContext(ctx).Errorf("update user role: %v", err.Error())
		return authpb.ErrorUpdateUserrole("更新用户角色列表失败")
	}
	return nil
}
