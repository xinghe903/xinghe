package repo

import (
	"auth/internal/biz/po"
	"context"
)

type UserRepo interface {
	Create(ctx context.Context, source *po.User) (string, error)
	Get(ctx context.Context, id string) (*po.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, cond *po.PageQuery[po.User]) (*po.SearchList[po.User], error)
}

type PermissionRepo interface {
	Create(ctx context.Context, source *po.Permission) (string, error)
	Get(ctx context.Context, id string) (*po.Permission, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, cond *po.PageQuery[po.Permission]) (*po.SearchList[po.Permission], error)
}

type RoleRepo interface {
	Create(ctx context.Context, source *po.Role) (string, error)
	Get(ctx context.Context, id string) (*po.Role, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, cond *po.PageQuery[po.Role]) (*po.SearchList[po.Role], error)
}

type RolePermissionRepo interface {
	Create(ctx context.Context, source *po.RolePermission) (string, error)
	Get(ctx context.Context, id string) (*po.RolePermission, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, cond *po.PageQuery[po.RolePermission]) (*po.SearchList[po.RolePermission], error)
}
