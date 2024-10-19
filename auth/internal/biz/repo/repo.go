package repo

import (
	"auth/internal/biz/po"
	"context"
)

type UserRepo interface {
	CreateUser(ctx context.Context, source *po.User) (string, error)
	UpdateUser(ctx context.Context, source *po.User) error
	GetUser(ctx context.Context, id string) (*po.User, error)
	DeleteUser(ctx context.Context, id string) error
	ListUser(ctx context.Context, cond *po.PageQuery[po.User]) (*po.SearchList[po.User], error)
}
