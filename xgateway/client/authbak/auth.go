package auth

import (
	"context"
	"log"

	// xgatewaypb "xgateway/api/auth/v1"
	xgatewaybakpb "xgateway/api/authbak/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Register(gwmux *runtime.ServeMux) {
	conn, err := grpc.NewClient(
		"0.0.0.0:9001",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	err = xgatewaybakpb.RegisterAuthServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}
}