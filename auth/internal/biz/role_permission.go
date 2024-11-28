package biz

import (
	authpb "auth/api/auth/v1"
	"auth/internal/biz/repo"
	"auth/internal/conf"
	"auth/internal/data/po"
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/xinghe903/xinghe/pkg/bo"
)

type RolePermissionUsecase struct {
	log    *log.Helper
	rRepo  repo.RoleRepo
	pRepo  repo.PermissionRepo
	rpRepo repo.RolePermissionRepo
}

func NewRolePermissionUsecase(c *conf.Config, logger log.Logger,
	r repo.RoleRepo, p repo.PermissionRepo, rp repo.RolePermissionRepo) *RolePermissionUsecase {
	return &RolePermissionUsecase{
		log:    log.NewHelper(logger),
		rRepo:  r,
		pRepo:  p,
		rpRepo: rp,
	}
}

func (a *RolePermissionUsecase) CreateRole(ctx context.Context, role *po.Role) (string, error) {
	id, err := a.rRepo.Create(ctx, role)
	if err != nil {
		a.log.WithContext(ctx).Errorf("create role: %v", err.Error())
		return "", authpb.ErrorCreateRole("创建角色失败")
	}
	return id, nil
}

func (a *RolePermissionUsecase) UpdateRole(ctx context.Context, role *po.Role) error {
	err := a.rRepo.Update(ctx, &po.Role{InstanceId: role.InstanceId, Name: role.Name})
	if err != nil {
		a.log.WithContext(ctx).Errorf("update role: %v", err.Error())
		return authpb.ErrorUpdateRole("更新角色失败")
	}
	return nil
}

func (a *RolePermissionUsecase) DeleteRole(ctx context.Context, role *po.Role) error {
	err := a.rRepo.Delete(ctx, role.InstanceId)
	if err != nil {
		a.log.WithContext(ctx).Errorf("delete role: %v", err.Error())
		return authpb.ErrorDeleteRole("删除角色失败")
	}
	return nil
}

func (a *RolePermissionUsecase) GetRole(ctx context.Context, id string) (*po.Role, error) {
	role, err := a.rRepo.Get(ctx, id)
	if err != nil {
		a.log.WithContext(ctx).Errorf("get role: %v", err.Error())
		return nil, authpb.ErrorGetRole("获取角色信息失败")
	}
	return role, nil
}

func (a *RolePermissionUsecase) ListRole(ctx context.Context, cond *bo.PageQuery[po.Role]) (*bo.SearchList[po.Role], error) {
	cond.Sort = []map[string]string{{"updated_at": "desc"}}
	list, err := a.rRepo.List(ctx, cond)
	if err != nil {
		a.log.WithContext(ctx).Errorf("list role: %v", err.Error())
		return nil, authpb.ErrorGetRole("获取角色列表失败")
	}
	return list, nil
}

func (a *RolePermissionUsecase) CreatePermission(ctx context.Context, req *po.Permission) (string, error) {
	id, err := a.pRepo.Create(ctx, req)
	if err != nil {
		a.log.WithContext(ctx).Errorf("create permission: %v", err.Error())
		return "", authpb.ErrorCreatePermission("创建权限失败")
	}
	if len(req.ParentId) != 0 {
		a.pRepo.Update(ctx, &po.Permission{InstanceId: req.ParentId, Children: po.PermissionHasChild})
	} else {
		a.pRepo.Update(ctx, &po.Permission{InstanceId: id, RootId: id, ParentId: "root"})
	}
	return id, nil
}

func (a *RolePermissionUsecase) UpdatePermission(ctx context.Context, req *po.Permission) error {
	err := a.pRepo.Update(ctx, &po.Permission{InstanceId: req.InstanceId, Permission: req.Permission})
	if err != nil {
		a.log.WithContext(ctx).Errorf("update permission: %v", err.Error())
		return authpb.ErrorUpdatePermission("更新权限失败")
	}
	if len(req.ParentId) != 0 {
		a.pRepo.Update(ctx, &po.Permission{InstanceId: req.ParentId, Children: po.PermissionHasChild})
	}
	// TODO remove no child
	return nil
}

func (a *RolePermissionUsecase) DeletePermission(ctx context.Context, req *po.Permission) error {
	err := a.pRepo.Delete(ctx, req.InstanceId)
	if err != nil {
		a.log.WithContext(ctx).Errorf("delete permission: %v", err.Error())
		return authpb.ErrorDeletePermission("删除权限失败")
	}
	// TODO remove no child
	return nil
}

func (a *RolePermissionUsecase) GetPermission(ctx context.Context, id string) (*po.Permission, error) {
	permission, err := a.pRepo.Get(ctx, id)
	if err != nil {
		a.log.WithContext(ctx).Errorf("get permission: %v", err.Error())
		return nil, authpb.ErrorGetPermission("获取权限信息失败")
	}
	return permission, nil
}

func (a *RolePermissionUsecase) ListPermission(ctx context.Context, cond *bo.PageQuery[po.Permission]) (*bo.SearchList[po.Permission], error) {
	cond.Sort = []map[string]string{{"sort": "asc"}, {"updated_at": "desc"}}
	list, err := a.pRepo.List(ctx, cond)
	if err != nil {
		a.log.WithContext(ctx).Errorf("list permission: %v", err.Error())
		return nil, authpb.ErrorGetPermission("获取权限列表失败")
	}
	return list, nil
}

func (a *RolePermissionUsecase) ListRolePermissions(ctx context.Context, roleId string) ([]string, error) {
	list, err := a.rpRepo.List(ctx, &bo.PageQuery[po.RolePermission]{
		Condition: &po.RolePermission{RoleId: roleId},
	})
	if err != nil {
		a.log.WithContext(ctx).Errorf("list role permissions: %v", err.Error())
		return nil, authpb.ErrorListRolepermissions("获取角色权限列表失败")
	}
	var permissionIds []string
	for _, rp := range list.Data {
		permissionIds = append(permissionIds, rp.PermissionId)
	}
	return permissionIds, nil
}

func (a *RolePermissionUsecase) UpdateRolePermissions(ctx context.Context, roleId string, permissions []string) error {
	err := a.rpRepo.CoverRelations(ctx, roleId, permissions)
	if err != nil {
		a.log.WithContext(ctx).Errorf("update role permissions: %v", err.Error())
		return authpb.ErrorUpdateRolepermission("更新角色权限列表失败")
	}
	return nil
}
