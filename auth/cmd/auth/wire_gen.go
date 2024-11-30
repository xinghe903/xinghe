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
func wireApp(confServer *conf.Server, confData *conf.Data, config *conf.Config, logger log.Logger) (*kratos.App, func(), error) {
	db, err := data.NewGormClient(confData)
	if err != nil {
		return nil, nil, err
	}
	dataData, cleanup, err := data.NewData(confData, db, logger)
	if err != nil {
		return nil, nil, err
	}
	sonyflake := data.NewSnowflake()
	userRepo := data.NewUserRepo(confServer, dataData, logger, sonyflake)
	authRepo := data.NewAuthRepo(confServer, dataData, logger, sonyflake)
	authUsecase := biz.NewAuthUsecase(config, logger, userRepo, sonyflake, authRepo)
	roleRepo := data.NewRoleRepo(confServer, dataData, logger, sonyflake)
	permissionRepo := data.NewPermissionRepo(confServer, dataData, logger, sonyflake)
	rolePermissionRepo := data.NewRolePermissionRepo(confServer, dataData, logger, sonyflake)
	rolePermissionUsecase := biz.NewRolePermissionUsecase(config, logger, roleRepo, permissionRepo, rolePermissionRepo)
	userRoleRepo := data.NewUserRoleRepo(confServer, dataData, logger, sonyflake)
	userRoleUsecase := biz.NewUserRoleUsecase(config, logger, userRoleRepo)
	authService := service.NewAuthService(logger, authUsecase, rolePermissionUsecase, userRoleUsecase)
	grpcServer := server.NewGRPCServer(confServer, authService, logger)
	httpServer := server.NewHTTPServer(confServer, authService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
