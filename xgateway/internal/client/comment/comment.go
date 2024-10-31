package comment

import (
	"context"
	"log"

	xgatewaypb "xgateway/api/comment/v1"
	"xgateway/internal/conf"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Comment struct {
	gwmux *runtime.ServeMux
	conn  *grpc.ClientConn
}

func NewComment(c *conf.Config, gwmux *runtime.ServeMux) *Comment {
	conn, err := grpc.NewClient(
		c.Client.Comment.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	return &Comment{
		gwmux: gwmux,
		conn:  conn,
	}
}

func (c *Comment) Start() {
	err := xgatewaypb.RegisterCommentServiceHandler(context.Background(), c.gwmux, c.conn)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}
}

func (c *Comment) Stop() {
	c.conn.Close()
}
