// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"auth/internal/biz"
	"auth/internal/conf"
	"auth/internal/data"
	"auth/internal/server"
	"auth/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	userRepo := data.NewUserRepo(confServer, dataData, logger)
	authUsecase := biz.NewAuthUsecase(logger, userRepo)
	authService := service.NewAuthService(authUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, authService, logger)
	httpServer := server.NewHTTPServer(confServer, authService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
