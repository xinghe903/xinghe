// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.0
// - protoc             v5.28.2
// source: auth/v1/auth.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationAuthServiceLogin = "/auth.v1.AuthService/Login"
const OperationAuthServiceLogout = "/auth.v1.AuthService/Logout"
const OperationAuthServiceRegister = "/auth.v1.AuthService/Register"

type AuthServiceHTTPServer interface {
	Login(context.Context, *LoginReq) (*LoginRsp, error)
	Logout(context.Context, *LogoutReq) (*emptypb.Empty, error)
	Register(context.Context, *RegisterReq) (*emptypb.Empty, error)
}

func RegisterAuthServiceHTTPServer(s *http.Server, srv AuthServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/auth/register", _AuthService_Register0_HTTP_Handler(srv))
	r.POST("/auth/login", _AuthService_Login0_HTTP_Handler(srv))
	r.POST("/auth/logout", _AuthService_Logout0_HTTP_Handler(srv))
}

func _AuthService_Register0_HTTP_Handler(srv AuthServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in RegisterReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthServiceRegister)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Register(ctx, req.(*RegisterReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _AuthService_Login0_HTTP_Handler(srv AuthServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in LoginReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthServiceLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Login(ctx, req.(*LoginReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*LoginRsp)
		return ctx.Result(200, reply)
	}
}

func _AuthService_Logout0_HTTP_Handler(srv AuthServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in LogoutReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthServiceLogout)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Logout(ctx, req.(*LogoutReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

type AuthServiceHTTPClient interface {
	Login(ctx context.Context, req *LoginReq, opts ...http.CallOption) (rsp *LoginRsp, err error)
	Logout(ctx context.Context, req *LogoutReq, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	Register(ctx context.Context, req *RegisterReq, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
}

type AuthServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewAuthServiceHTTPClient(client *http.Client) AuthServiceHTTPClient {
	return &AuthServiceHTTPClientImpl{client}
}

func (c *AuthServiceHTTPClientImpl) Login(ctx context.Context, in *LoginReq, opts ...http.CallOption) (*LoginRsp, error) {
	var out LoginRsp
	pattern := "/auth/login"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthServiceLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *AuthServiceHTTPClientImpl) Logout(ctx context.Context, in *LogoutReq, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/auth/logout"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthServiceLogout))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *AuthServiceHTTPClientImpl) Register(ctx context.Context, in *RegisterReq, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/auth/register"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAuthServiceRegister))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}