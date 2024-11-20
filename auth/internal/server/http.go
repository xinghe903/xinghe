package server

import (
	authpb "auth/api/auth/v1"
	"auth/internal/conf"
	"auth/internal/service"

	"auth/internal/server/middleware"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, auth *service.AuthService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			// middleware.MiddlewareCors(),
			middleware.LogMiddleware(logger),
			middleware.HeaderMiddleware(),
			metadata.Server(),
		),
		// http.Filter(handlers.CORS(
		// 	handlers.AllowedOrigins([]string{"http://172.18.20.117:5173"}),
		// 	handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		// )),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	authpb.RegisterAuthServiceHTTPServer(srv, auth)
	return srv
}
