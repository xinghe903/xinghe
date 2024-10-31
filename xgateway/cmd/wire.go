//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"xgateway/internal/client/auth"
	"xgateway/internal/client/comment"
	"xgateway/internal/conf"

	"github.com/google/wire"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func wireApp(c *conf.Config) (*Server, error) {
	panic(wire.Build(comment.NewComment,
		auth.NewAuth,
		NewServer, newGwMux))
}

type Client interface {
	Start()
	Stop()
}

type Server struct {
	clients []Client
	gw      *runtime.ServeMux
}

func NewServer(gw *runtime.ServeMux, co *comment.Comment, au *auth.Auth) *Server {
	clients := []Client{co, au}
	return &Server{
		clients: clients,
		gw:      gw,
	}
}

func (s *Server) Start() {
	for _, c := range s.clients {
		c.Start()
	}
}

func (s *Server) Stop() {
	for _, c := range s.clients {
		c.Stop()
	}
}
