package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

const TOKEN = "access_token"

func LogMiddleware(l log.Logger) middleware.Middleware {
	logger := log.NewHelper(l)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			start := time.Now()
			var operation string
			if info, ok := transport.FromServerContext(ctx); ok {
				operation = info.Operation()
			}
			logger.WithContext(ctx).Debugf("operation: %s, request: %+v\n", operation, req)
			rsp, err := handler(ctx, req)
			end := time.Now()
			latency := end.Sub(start)
			logger.WithContext(ctx).Debugf("latency: %v response %+v\n", latency, rsp)
			return rsp, err
		}
	}
}

func HeaderMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if info, ok := transport.FromServerContext(ctx); ok {
				header := info.RequestHeader()
				if token := header.Get(TOKEN); token != "" {
					ctx = context.WithValue(ctx, TOKEN, token)
				}
			}
			rsp, err := handler(ctx, req)
			return rsp, err
		}
	}
}

func MiddlewareCors() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			fmt.Printf("req %+v\n", req)
			if ht, ok := transport.FromServerContext(ctx); ok {
				ht.ReplyHeader().Set("Access-Control-Allow-Origin", "*")
				ht.ReplyHeader().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,PATCH,DELETE")
				ht.ReplyHeader().Set("Access-Control-Allow-Credentials", "true")
				ht.ReplyHeader().Set("Access-Control-Allow-Headers", "Content-Type,"+
					"X-Requested-With,Access-Control-Allow-Credentials,User-Agent,Content-Length,Authorization")
			}
			return handler(ctx, req)
		}
	}
}
