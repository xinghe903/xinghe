package repo

import (
	"auth/internal/data/po"
	"context"
	"time"

	"github.com/xinghe903/xinghe/pkg/bo"
)

type UserRepo interface {
	Create(ctx context.Context, source *po.User) (string, error)
	Update(ctx context.Context, source *po.User) error
	Get(ctx context.Context, id string) (*po.User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, cond *bo.PageQuery[po.User], username string) (*bo.SearchList[po.User], error)
}

type AuthRepo interface {
	Create(ctx context.Context, source *po.Auth) (string, error)
	Update(ctx context.Context, source *po.Auth) error
	UpdateByCode(ctx context.Context, source *po.Auth) error
	Get(ctx context.Context, id string) (*po.Auth, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, cond *bo.PageQuery[po.Auth]) (*bo.SearchList[po.Auth], error)
	GetByToken(ctx context.Context, token string) (*po.Auth, error)
	UpdateByToken(ctx context.Context, source *po.Auth) error
	ClearExpiredToken(ctx context.Context, date time.Time) error
}

type PermissionRepo interface {
	Create(ctx context.Context, source *po.Permission) (string, error)
	Update(ctx context.Context, source *po.Permission) error
	Get(ctx context.Context, id string) (*po.Permission, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, cond *bo.PageQuery[po.Permission]) (*bo.SearchList[po.Permission], error)
}

type RoleRepo interface {
	Create(ctx context.Context, source *po.Role) (string, error)
	Update(ctx context.Context, source *po.Role) error
	Get(ctx context.Context, id string) (*po.Role, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, cond *bo.PageQuery[po.Role]) (*bo.SearchList[po.Role], error)
}

type RolePermissionRepo interface {
	Create(ctx context.Context, source *po.RolePermission) (string, error)
	Update(ctx context.Context, source *po.RolePermission) error
	Get(ctx context.Context, id string) (*po.RolePermission, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, cond *bo.PageQuery[po.RolePermission]) (*bo.SearchList[po.RolePermission], error)
	CoverRelations(ctx context.Context, id string, data []string) error
}

type UserRoleRepo interface {
	Create(ctx context.Context, source *po.UserRole) (string, error)
	Update(ctx context.Context, source *po.UserRole) error
	Get(ctx context.Context, id string) (*po.UserRole, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, cond *bo.PageQuery[po.UserRole]) (*bo.SearchList[po.UserRole], error)
	CoverRelations(ctx context.Context, id string, data []string) error
}
