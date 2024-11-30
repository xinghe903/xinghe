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

var _ repo.UserRoleRepo = &userRoleRepo{}

type userRoleRepo struct {
	db   *gorm.DB
	log  *log.Helper
	snow *hashid.Sonyflake
}

func NewUserRoleRepo(c *conf.Server, data *Data, logger log.Logger, id *hashid.Sonyflake) repo.UserRoleRepo {
	return &userRoleRepo{
		db:   data.db,
		log:  log.NewHelper(logger),
		snow: id,
	}
}

func (p *userRoleRepo) Create(ctx context.Context, source *po.UserRole) (string, error) {
	if len(source.InstanceId) == 0 {
		source.InstanceId = source.GenerateID(p.snow.GenerateID())
	}
	if err := p.db.Create(source).Error; err != nil {
		return "", fmt.Errorf("%s create: %s", source.TableName(), err.Error())
	}
	return source.InstanceId, nil
}
func (p *userRoleRepo) Update(ctx context.Context, source *po.UserRole) error {
	if len(source.InstanceId) == 0 {
		return errors.New("instance id is required")
	}
	if err := p.db.Model(&po.UserRole{}).Where("instanceId = ?", source.InstanceId).Updates(source).Error; err != nil {
		return fmt.Errorf("%s update: %s", source.TableName(), err.Error())
	}
	return nil
}
func (p *userRoleRepo) Get(ctx context.Context, id string) (*po.UserRole, error) {
	if len(id) == 0 {
		return nil, errors.New("instance id is required")
	}
	pm := &po.UserRole{}
	if err := p.db.Model(pm).Where("instanceId = ?", id).First(pm).Error; err != nil {
		return nil, fmt.Errorf("%s get: %s", pm.TableName(), err.Error())
	}
	return pm, nil
}
func (p *userRoleRepo) Delete(ctx context.Context, id string) error {
	if len(id) == 0 {
		return errors.New("instance id is required")
	}
	pm := &po.UserRole{}
	if err := p.db.Where("instanceId =?", id).Delete(pm).Error; err != nil {
		return fmt.Errorf("%s get: %s", pm.TableName(), err.Error())
	}
	return nil
}
func (p *userRoleRepo) List(ctx context.Context, cond *bo.PageQuery[po.UserRole]) (*bo.SearchList[po.UserRole], error) {
	if cond == nil {
		return nil, errors.New("condition is required")
	}
	var rsp bo.SearchList[po.UserRole]
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

func (p *userRoleRepo) CoverRelations(ctx context.Context, id string, data []string) error {
	if len(id) == 0 || len(data) == 0 {
		return nil
	}
	var datas []*po.UserRole
	source := &po.UserRole{}
	for _, d := range data {
		datas = append(datas, &po.UserRole{
			InstanceId: source.GenerateID(p.snow.GenerateID()),
			RoleId:     d,
			UserId:     id,
		})
	}
	if err := p.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("userId = ? ", id).Delete(source).Error; err != nil {
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
