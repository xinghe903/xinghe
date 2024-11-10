// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"xgateway/internal/client/auth"
	"xgateway/internal/client/comment"
	"xgateway/internal/conf"
)

// Injectors from wire.go:

func wireApp(c *conf.Config) (*Server, error) {
	serveMux := newGwMux()
	commentComment := comment.NewComment(c, serveMux)
	authAuth := auth.NewAuth(c, serveMux)
	server := NewServer(serveMux, commentComment, authAuth)
	return server, nil
}

// wire.go:

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
