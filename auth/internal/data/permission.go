package data

import (
	"auth/internal/biz/po"
	"auth/internal/biz/repo"
	"auth/internal/conf"
	"context"
	"errors"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	hashid "github.com/xinghe903/xinghe/pkg/distribute/id"
	"gorm.io/gorm"
)

var _ repo.PermissionRepo = &permissionRepo{}

type permissionRepo struct {
	db   *gorm.DB
	log  *log.Helper
	snow *hashid.Snowflake
}

func NewPermissionRepo(c *conf.Server, data *Data, logger log.Logger, id *hashid.Snowflake) repo.PermissionRepo {
	return &permissionRepo{
		db:   data.db,
		log:  log.NewHelper(logger),
		snow: id,
	}
}

func (p *permissionRepo) Create(ctx context.Context, source *po.Permission) (string, error) {
	if len(source.InstanceId) == 0 {
		source.InstanceId = source.GenerateID(p.snow.GenerateID())
	}
	if err := p.db.Create(source).Error; err != nil {
		return "", fmt.Errorf("%s create: %s", source.TableName(), err.Error())
	}
	return source.InstanceId, nil
}
func (p *permissionRepo) Update(ctx context.Context, source *po.Permission) error {
	if len(source.InstanceId) == 0 {
		return errors.New("instance id is required")
	}
	if err := p.db.Model(&po.Permission{}).Where("instanceId = ?", source.InstanceId).Updates(source).Error; err != nil {
		return fmt.Errorf("%s update: %s", source.TableName(), err.Error())
	}
	return nil
}
func (p *permissionRepo) Get(ctx context.Context, id string) (*po.Permission, error) {
	if len(id) == 0 {
		return nil, errors.New("instance id is required")
	}
	pm := &po.Permission{}
	if err := p.db.Model(pm).Where("instanceId = ?", id).First(pm).Error; err != nil {
		return nil, fmt.Errorf("%s get: %s", pm.TableName(), err.Error())
	}
	return pm, nil
}
func (p *permissionRepo) Delete(ctx context.Context, id string) error {
	if len(id) == 0 {
		return errors.New("instance id is required")
	}
	pm := &po.Permission{}
	if err := p.db.Where("instanceId =?", id).Delete(pm).Error; err != nil {
		return fmt.Errorf("%s get: %s", pm.TableName(), err.Error())
	}
	return nil
}
func (p *permissionRepo) List(ctx context.Context, cond *po.PageQuery[po.Permission]) (*po.SearchList[po.Permission], error) {
	if cond == nil {
		return nil, errors.New("condition is required")
	}
	var rsp po.SearchList[po.Permission]
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
