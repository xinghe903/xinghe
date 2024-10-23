package service

import (
	commentpb "comment/api/comment/v1"
	"comment/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type CommentService struct {
	commentpb.UnimplementedCommentServiceServer
	cu     *biz.CommentUsecase
	logger *log.Helper
}

func NewCommentService(cu *biz.CommentUsecase, logger log.Logger) *CommentService {
	return &CommentService{cu: cu, logger: log.NewHelper(logger)}
}

func (c *CommentService) CreateComment(ctx context.Context, req *commentpb.CreateCommentReq) (*commentpb.CreateCommentRsp, error) {
	return &commentpb.CreateCommentRsp{CommentId: "123"}, nil
}
func (c *CommentService) GetComment(ctx context.Context, req *commentpb.GetCommentReq) (*commentpb.GetCommentRsp, error) {
	return &commentpb.GetCommentRsp{Comment: &commentpb.Comment{Content: "水电费"}}, nil
}
