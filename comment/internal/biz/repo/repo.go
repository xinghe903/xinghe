package repo

import (
	"comment/internal/biz/po"
	"context"
)

type CommentRepo interface {
	CreateComment(ctx context.Context, source *po.Comment) (string, error)
	UpdateComment(ctx context.Context, source *po.Comment) error
	GetComment(ctx context.Context, id string) (*po.Comment, error)
	DeleteComment(ctx context.Context, id string) error
	ListComment(ctx context.Context, cond *po.PageQuery[po.Comment]) (*po.SearchList[po.Comment], error)
}
