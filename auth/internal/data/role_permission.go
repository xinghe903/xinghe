package data

import (
	"auth/internal/biz/repo"
	"auth/internal/conf"
	"auth/internal/data/po"
	"context"
	"errors"
	"fmt"

	"github.com/xinghe903/xinghe/pkg/bo"

	"github.com/go-kratos/kratos/v2/log"
	hashid "github.com/xinghe903/xinghe/pkg/distribute/id"
	"gorm.io/gorm"
)

var _ repo.RolePermissionRepo = &rolePermissionRepo{}

type rolePermissionRepo struct {
	db   *gorm.DB
	log  *log.Helper
	snow *hashid.Sonyflake
}

func NewRolePermissionRepo(c *conf.Server, data *Data, logger log.Logger, id *hashid.Sonyflake) repo.RolePermissionRepo {
	return &rolePermissionRepo{
		db:   data.db,
		log:  log.NewHelper(logger),
		snow: id,
	}
}

func (p *rolePermissionRepo) Create(ctx context.Context, source *po.RolePermission) (string, error) {
	if len(source.InstanceId) == 0 {
		source.InstanceId = source.GenerateID(p.snow.GenerateID())
	}
	if err := p.db.Create(source).Error; err != nil {
		return "", fmt.Errorf("%s create: %s", source.TableName(), err.Error())
	}
	return source.InstanceId, nil
}
func (p *rolePermissionRepo) Update(ctx context.Context, source *po.RolePermission) error {
	if len(source.InstanceId) == 0 {
		return errors.New("instance id is required")
	}
	if err := p.db.Model(&po.RolePermission{}).Where("instanceId = ?", source.InstanceId).Updates(source).Error; err != nil {
		return fmt.Errorf("%s update: %s", source.TableName(), err.Error())
	}
	return nil
}
func (p *rolePermissionRepo) Get(ctx context.Context, id string) (*po.RolePermission, error) {
	if len(id) == 0 {
		return nil, errors.New("instance id is required")
	}
	pm := &po.RolePermission{}
	if err := p.db.Model(pm).Where("instanceId = ?", id).First(pm).Error; err != nil {
		return nil, fmt.Errorf("%s get: %s", pm.TableName(), err.Error())
	}
	return pm, nil
}
func (p *rolePermissionRepo) Delete(ctx context.Context, id string) error {
	if len(id) == 0 {
		return errors.New("instance id is required")
	}
	pm := &po.RolePermission{}
	if err := p.db.Where("instanceId =?", id).Delete(pm).Error; err != nil {
		return fmt.Errorf("%s get: %s", pm.TableName(), err.Error())
	}
	return nil
}
func (p *rolePermissionRepo) List(ctx context.Context, cond *bo.PageQuery[po.RolePermission]) (*bo.SearchList[po.RolePermission], error) {
	if cond == nil {
		return nil, errors.New("condition is required")
	}
	var rsp bo.SearchList[po.RolePermission]
	md := p.db.Model(cond.Condition)
	if cond.Condition != nil {
		md = md.Where(cond.Condition)
	}
	md = md.Count(&rsp.Total)
	if cond.PageNum != 0 && cond.PageSize != 0 {
		md = md.Offset(int(cond.PageNum-1) * int(cond.PageSize)).Limit(int(cond.PageSize))
	}
	for _, item := range cond.Sort {
		for k, v := range item {
			md = md.Order(fmt.Sprintf("%s %s", k, v))
		}
	}
	if err := md.Find(&rsp.Data).Error; err != nil {
		return nil, fmt.Errorf("%s get: %s", cond.Condition.TableName(), err.Error())
	}
	return &rsp, nil
}

func (p *rolePermissionRepo) CoverRelations(ctx context.Context, id string, data []string) error {
	if len(id) == 0 || len(data) == 0 {
		return nil
	}
	var datas []*po.RolePermission
	source := &po.RolePermission{}
	for _, d := range data {
		datas = append(datas, &po.RolePermission{
			InstanceId:   source.GenerateID(p.snow.GenerateID()),
			PermissionId: d,
			RoleId:       id,
		})
	}
	if err := p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("roleId = ? ", id).Delete(source).Error; err != nil {
			return err
		}
		if err := tx.Create(datas).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return fmt.Errorf("cover relations: %s", err.Error())
	}
	return nil
}
