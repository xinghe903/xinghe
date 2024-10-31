package biz

import (
	"auth/internal/biz/po"
	"auth/internal/biz/repo"
	"auth/internal/conf"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/xinghe903/xinghe/pkg/encrypt"
)

type RolePermissionUsecase struct {
	log   *log.Helper
	uRepo repo.UserRepo
	enc   *encrypt.EncryptAes
}

func NewRolePermissionUsecase(c *conf.Config, logger log.Logger, u repo.UserRepo) *RolePermissionUsecase {
	return &RolePermissionUsecase{
		log:   log.NewHelper(logger),
		uRepo: u,
		enc:   encrypt.NewEncryptAes(c.EncryptKey),
	}
}

func (a *RolePermissionUsecase) CreateRole(ctx context.Context, role *po.Role) (string, error) {

	return "", nil
}

func (a *RolePermissionUsecase) UpdateRole(ctx context.Context, role *po.Role) error {

	return nil
}

func (a *RolePermissionUsecase) DeleteRole(ctx context.Context, role *po.Role) error {

	return nil
}

func (a *RolePermissionUsecase) GetRole(ctx context.Context, id string) (*po.Role, error) {

	return nil, nil
}

func (a *RolePermissionUsecase) ListRole(ctx context.Context, cond *po.PageQuery[po.Role]) (*po.SearchList[po.Role], error) {

	return nil, nil
}

func (a *RolePermissionUsecase) CreatePermission(ctx context.Context, role *po.Permission) (string, error) {

	return "", nil
}

func (a *RolePermissionUsecase) UpdatePermission(ctx context.Context, role *po.Permission) error {

	return nil
}

func (a *RolePermissionUsecase) DeletePermission(ctx context.Context, role *po.Permission) error {

	return nil
}

func (a *RolePermissionUsecase) GetPermission(ctx context.Context, id string) (*po.Permission, error) {

	return nil, nil
}

func (a *RolePermissionUsecase) ListPermission(ctx context.Context, cond *po.PageQuery[po.Permission]) (*po.SearchList[po.Permission], error) {

	return nil, nil
}
