package auth

import (
	"context"
	"log"

	xgatewaypb "xgateway/api/auth/v1"
	"xgateway/internal/conf"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Auth struct {
	gwmux *runtime.ServeMux
	conn  *grpc.ClientConn
}

func NewAuth(c *conf.Config, gwmux *runtime.ServeMux) *Auth {
	conn, err := grpc.NewClient(
		c.Client.Auth.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	return &Auth{
		gwmux: gwmux,
		conn:  conn,
	}
}

func (c *Auth) Start() {
	err := xgatewaypb.RegisterAuthServiceHandler(context.Background(), c.gwmux, c.conn)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}
}

func (c *Auth) Stop() {
	c.conn.Close()
}
