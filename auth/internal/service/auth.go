package service

import (
	authpb "auth/api/auth/v1"
	"auth/internal/biz"
	"auth/internal/biz/po"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
	uc     *biz.AuthUsecase
	rp     *biz.RolePermissionUsecase
	logger *log.Helper
}

func NewAuthService(logger log.Logger, uc *biz.AuthUsecase, rp *biz.RolePermissionUsecase) *AuthService {
	return &AuthService{logger: log.NewHelper(logger), uc: uc, rp: rp}
}

func (s *AuthService) Register(ctx context.Context, req *authpb.RegisterReq) (*emptypb.Empty, error) {
	username, password := req.Username, req.Password
	if len(username) == 0 || len(password) == 0 {
		s.logger.WithContext(ctx).Warnf("username and password is required")
		return nil, authpb.ErrorUserOrPasswordEmpty("用户名或密码不能为空")
	}
	_, err := s.uc.Register(ctx, &po.User{Name: username, Password: password})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *AuthService) Login(ctx context.Context, req *authpb.LoginReq) (*authpb.LoginRsp, error) {
	username, password := req.Username, req.Password
	if len(username) == 0 || len(password) == 0 {
		return nil, authpb.ErrorUserOrPasswordEmpty("用户名或密码不能为空")
	}
	token, err := s.uc.Login(ctx, &po.User{Name: username, Password: password})
	if err != nil {
		return nil, err
	}
	return &authpb.LoginRsp{Token: token}, nil
}

func (s *AuthService) Logout(ctx context.Context, req *authpb.LogoutReq) (*emptypb.Empty, error) {
	err := s.uc.Logout(ctx, req.Token)
	return nil, err
}

func (s *AuthService) Auth(ctx context.Context, req *authpb.AuthReq) (*authpb.AuthRsp, error) {
	user, err := s.uc.Auth(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &authpb.AuthRsp{
		Username: user.Name,
		Nickname: user.NickName,
		UserId:   user.InstanceId,
	}, nil
}

func (s *AuthService) GetUser(ctx context.Context, req *authpb.GetUserReq) (*authpb.GetUserRsp, error) {
	if len(req.Id) == 0 {
		return nil, authpb.ErrorUserOrPasswordEmpty("id不能为空")
	}
	user, err := s.uc.GetUserById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &authpb.GetUserRsp{
		Id:        user.InstanceId,
		Username:  user.Name,
		Nickname:  user.NickName,
		CreatedAt: user.CreatedAt.GoString(),
		UpdatedAt: user.UpdatedAt.GoString(),
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    "",
	}, nil
}

func (s *AuthService) CreateRole(ctx context.Context, req *authpb.CreateRoleReq) (*emptypb.Empty, error) {
	s.rp.CreateRole(ctx, &po.Role{
		Name: req.Name,
	})
	return nil, nil
}

func (s *AuthService) UpdateRole(ctx context.Context, req *authpb.UpdateRoleReq) (*emptypb.Empty, error) {
	s.rp.UpdateRole(ctx, &po.Role{
		InstanceId: req.Id,
		Name:       req.Name,
	})
	return nil, nil
}

func (s *AuthService) DeleteRole(ctx context.Context, req *authpb.DeleteRoleReq) (*emptypb.Empty, error) {
	s.rp.DeleteRole(ctx, &po.Role{
		InstanceId: req.Id,
	})
	return nil, nil
}

func (s *AuthService) GetRole(ctx context.Context, req *authpb.GetRoleReq) (*authpb.GetRoleRsp, error) {
	role, err := s.rp.GetRole(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &authpb.GetRoleRsp{
		Id:   role.InstanceId,
		Name: role.Name,
	}, nil
}

func (s *AuthService) ListRole(ctx context.Context, req *authpb.ListRoleReq) (*authpb.ListRoleRsp, error) {
	result, err := s.rp.ListRole(ctx, &po.PageQuery[po.Role]{
		PageNum:  req.PageNumber,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	ret := &authpb.ListRoleRsp{Total: int32(result.Total)}
	for _, r := range result.Data {
		ret.Roles = append(ret.Roles, &authpb.Role{Id: r.InstanceId, Name: r.Name})
	}
	return ret, nil
}

func (s *AuthService) CreatePermission(ctx context.Context, req *authpb.CreatePermissionReq) (*emptypb.Empty, error) {
	s.rp.CreatePermission(ctx, &po.Permission{})
	return nil, nil
}

func (s *AuthService) UpdatePermission(ctx context.Context, req *authpb.UpdatePermissionReq) (*emptypb.Empty, error) {
	s.rp.UpdatePermission(ctx, &po.Permission{})
	return nil, nil
}

func (s *AuthService) DeletePermission(ctx context.Context, req *authpb.DeletePermissionReq) (*emptypb.Empty, error) {
	s.rp.DeletePermission(ctx, &po.Permission{})
	return nil, nil
}

func (s *AuthService) GetPermission(ctx context.Context, req *authpb.GetPermissionReq) (*authpb.GetPermissionRsp, error) {
	permission, err := s.rp.GetPermission(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &authpb.GetPermissionRsp{
		Id:   permission.InstanceId,
		Name: permission.Permission,
	}, nil
}

func (s *AuthService) ListPermission(ctx context.Context, req *authpb.ListPermissionReq) (*authpb.ListPermissionRsp, error) {
	result, err := s.rp.ListPermission(ctx, &po.PageQuery[po.Permission]{})
	if err != nil {
		return nil, err
	}
	ret := &authpb.ListPermissionRsp{Total: int32(result.Total)}
	for _, p := range result.Data {
		ret.Permissions = append(ret.Permissions, &authpb.Permission{Id: p.InstanceId, Permission: p.Permission})
	}
	return ret, nil
}
