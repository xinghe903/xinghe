package service

import (
	authpb "auth/api/auth/v1"
	"auth/internal/biz"
	"auth/internal/biz/po"
	"context"
	"database/sql"
	"time"

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

func (s *AuthService) CreateUser(ctx context.Context, req *authpb.CreateUserReq) (*authpb.CreateUserRsp, error) {
	rsp, err := s.uc.CreateUser(ctx, &po.User{
		Name:     req.Username,
		NickName: req.Nickname,
		Email:    sql.NullString{String: req.Email, Valid: true},
		Phone:    sql.NullString{String: req.Phone, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return &authpb.CreateUserRsp{
		Id:        rsp.InstanceId,
		Username:  req.Username,
		Nickname:  req.Nickname,
		Password:  rsp.Password,
		CreatedAt: time.Now().Format(time.DateTime),
		UpdatedAt: time.Now().Format(time.DateTime),
		Email:     req.Email,
		Phone:     req.Phone,
	}, nil
}

func (s *AuthService) UpdateUser(ctx context.Context, req *authpb.UpdateUserReq) (*emptypb.Empty, error) {
	err := s.uc.UpdateUser(ctx, &po.User{
		InstanceId: req.Id,
		NickName:   req.Nickname,
		Email:      sql.NullString{String: req.Email, Valid: true},
		Phone:      sql.NullString{String: req.Phone, Valid: true},
	})
	return nil, err
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
		CreatedAt: user.CreatedAt.Format(time.DateTime),
		UpdatedAt: user.UpdatedAt.Format(time.DateTime),
		Email:     user.Email.String,
		Phone:     user.Phone.String,
		Avatar:    "",
	}, nil
}

func (s *AuthService) ListUser(ctx context.Context, req *authpb.ListUserReq) (*authpb.ListUserRsp, error) {
	result, err := s.uc.ListUser(ctx, &po.PageQuery[po.User]{
		PageNum:  req.PageNumber,
		PageSize: req.PageSize,
	}, req.Username)
	if err != nil {
		return nil, err
	}
	ret := &authpb.ListUserRsp{Total: int32(result.Total)}
	for _, user := range result.Data {
		ret.Users = append(ret.Users, &authpb.User{
			Id:        user.InstanceId,
			Username:  user.Name,
			Nickname:  user.NickName,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
			UpdatedAt: user.UpdatedAt.Format(time.DateTime),
			Email:     user.Email.String,
			Phone:     user.Phone.String,
			Avatar:    "",
		})
	}
	return ret, nil
}

func (s *AuthService) CreateRole(ctx context.Context, req *authpb.CreateRoleReq) (*authpb.CreateRoleRsp, error) {
	id, err := s.rp.CreateRole(ctx, &po.Role{
		Name: req.Name,
	})
	return &authpb.CreateRoleRsp{Id: id, Name: req.Name}, err
}

func (s *AuthService) UpdateRole(ctx context.Context, req *authpb.UpdateRoleReq) (*emptypb.Empty, error) {
	err := s.rp.UpdateRole(ctx, &po.Role{
		InstanceId: req.Id,
		Name:       req.Name,
	})
	return nil, err
}

func (s *AuthService) DeleteRole(ctx context.Context, req *authpb.DeleteRoleReq) (*emptypb.Empty, error) {
	err := s.rp.DeleteRole(ctx, &po.Role{
		InstanceId: req.Id,
	})
	return nil, err
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
		ret.Roles = append(ret.Roles, &authpb.Role{
			Id:        r.InstanceId,
			Name:      r.Name,
			UpdatedAt: r.UpdatedAt.Format(time.DateTime),
			CreatedAt: r.CreatedAt.Format(time.DateTime),
		})
	}
	return ret, nil
}

func (s *AuthService) CreatePermission(ctx context.Context, req *authpb.CreatePermissionReq) (*emptypb.Empty, error) {
	_, err := s.rp.CreatePermission(ctx, &po.Permission{
		Permission:  req.Permission,
		SubjectType: po.NewPermissionType(req.SubjectType),
		SubjectId:   req.SubjectId,
		RootId:      req.RootId,
		ParentId:    req.ParentId,
		Name:        req.Name,
		Sort:        req.Sort,
	})
	return nil, err
}

func (s *AuthService) UpdatePermission(ctx context.Context, req *authpb.UpdatePermissionReq) (*emptypb.Empty, error) {
	err := s.rp.UpdatePermission(ctx, &po.Permission{
		InstanceId:  req.Id,
		Permission:  req.Permission,
		SubjectType: po.NewPermissionType(req.SubjectType),
		SubjectId:   req.SubjectId,
		RootId:      req.RootId,
		ParentId:    req.ParentId,
		Name:        req.Name,
		Sort:        req.Sort,
	})
	return nil, err
}

func (s *AuthService) DeletePermission(ctx context.Context, req *authpb.DeletePermissionReq) (*emptypb.Empty, error) {
	err := s.rp.DeletePermission(ctx, &po.Permission{InstanceId: req.Id})
	return nil, err
}

func (s *AuthService) GetPermission(ctx context.Context, req *authpb.GetPermissionReq) (*authpb.GetPermissionRsp, error) {
	permission, err := s.rp.GetPermission(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &authpb.GetPermissionRsp{
		Id:          permission.InstanceId,
		Name:        permission.Name,
		SubjectType: string(permission.SubjectType),
		SubjectId:   permission.SubjectId,
		RootId:      permission.RootId,
		ParentId:    permission.ParentId,
		Sort:        permission.Sort,
		Permission:  permission.Permission,
		CreatedAt:   permission.CreatedAt.Format(time.DateTime),
		UpdatedAt:   permission.UpdatedAt.Format(time.DateTime),
	}, nil
}

func (s *AuthService) ListPermission(ctx context.Context, req *authpb.ListPermissionReq) (*authpb.ListPermissionRsp, error) {
	result, err := s.rp.ListPermission(ctx, &po.PageQuery[po.Permission]{})
	if err != nil {
		return nil, err
	}
	ret := &authpb.ListPermissionRsp{Total: int32(result.Total)}
	for _, permission := range result.Data {
		ret.Permissions = append(ret.Permissions, &authpb.Permission{
			Id:          permission.InstanceId,
			Name:        permission.Name,
			SubjectType: string(permission.SubjectType),
			SubjectId:   permission.SubjectId,
			RootId:      permission.RootId,
			ParentId:    permission.ParentId,
			Sort:        permission.Sort,
			Permission:  permission.Permission,
			CreatedAt:   permission.CreatedAt.Format(time.DateTime),
			UpdatedAt:   permission.UpdatedAt.Format(time.DateTime),
			Children:    permission.Children,
		})
	}
	return ret, nil
}
