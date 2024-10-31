package data

import (
	"auth/internal/biz/po"
	"auth/internal/biz/repo"
	"auth/internal/conf"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	hashid "github.com/xinghe903/xinghe/pkg/distribute/id"
	"gorm.io/gorm"
)

var _ repo.RoleRepo = &roleRepo{}

type roleRepo struct {
	db        *gorm.DB
	log       *log.Helper
	builderId *hashid.Snowflake
}

func NewRoleRepo(c *conf.Server, data *Data, logger log.Logger) repo.RoleRepo {
	nodeId, err := strconv.Atoi(c.GetNodeId())
	if err != nil {
		panic(err.Error())
	}
	return &roleRepo{
		db:        data.db,
		log:       log.NewHelper(logger),
		builderId: hashid.NewSnowflake(int64(nodeId)),
	}
}

func (p *roleRepo) Create(ctx context.Context, source *po.Role) (string, error) {
	if len(source.InstanceId) == 0 {
		source.InstanceId = source.GenerateID(p.builderId.GenerateID())
	}
	if err := p.db.Create(source).Error; err != nil {
		return "", fmt.Errorf("%s create: %s", source.TableName(), err.Error())
	}
	return source.InstanceId, nil
}
func (p *roleRepo) Update(ctx context.Context, source *po.Role) error {
	if len(source.InstanceId) == 0 {
		return errors.New("instance id is required")
	}
	if err := p.db.Model(&po.Role{}).Where("instanceId = ?", source.InstanceId).Updates(source).Error; err != nil {
		return fmt.Errorf("%s update: %s", source.TableName(), err.Error())
	}
	return nil
}
func (p *roleRepo) Get(ctx context.Context, id string) (*po.Role, error) {
	if len(id) == 0 {
		return nil, errors.New("instance id is required")
	}
	pm := &po.Role{}
	if err := p.db.Model(pm).Where("instanceId = ?", id).First(pm).Error; err != nil {
		return nil, fmt.Errorf("%s get: %s", pm.TableName(), err.Error())
	}
	return pm, nil
}
func (p *roleRepo) Delete(ctx context.Context, id string) error {
	if len(id) == 0 {
		return errors.New("instance id is required")
	}
	pm := &po.Role{}
	if err := p.db.Where("instanceId =?", id).Delete(pm).Error; err != nil {
		return fmt.Errorf("%s get: %s", pm.TableName(), err.Error())
	}
	return nil
}
func (p *roleRepo) List(ctx context.Context, cond *po.PageQuery[po.Role]) (*po.SearchList[po.Role], error) {
	if cond == nil {
		return nil, errors.New("condition is required")
	}
	var rsp po.SearchList[po.Role]
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