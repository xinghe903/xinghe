package data

import (
	"auth/internal/biz/repo"
	"auth/internal/conf"
	"auth/internal/data/po"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/xinghe903/xinghe/pkg/bo"
	hashid "github.com/xinghe903/xinghe/pkg/distribute/id"
	"gorm.io/gorm"
)

var _ repo.UserRepo = &userRepo{}

type userRepo struct {
	db   *gorm.DB
	log  *log.Helper
	snow *hashid.Sonyflake
}

func NewUserRepo(c *conf.Server, data *Data, logger log.Logger, id *hashid.Sonyflake) repo.UserRepo {
	return &userRepo{
		db:   data.db,
		log:  log.NewHelper(logger),
		snow: id,
	}
}

func (p *userRepo) Create(ctx context.Context, source *po.User) (string, error) {
	if len(source.InstanceId) == 0 {
		source.InstanceId = strconv.FormatUint(p.snow.GenerateID(), 10)
	}
	if err := p.db.Create(source).Error; err != nil {
		return "", fmt.Errorf("%s create: %s", source.TableName(), err.Error())
	}
	return source.InstanceId, nil
}
func (p *userRepo) Update(ctx context.Context, source *po.User) error {
	if len(source.InstanceId) == 0 {
		return errors.New("instance id is required")
	}
	if err := p.db.Model(&po.User{}).Where("instanceId = ?", source.InstanceId).Updates(source).Error; err != nil {
		return fmt.Errorf("%s update: %s", source.TableName(), err.Error())
	}
	return nil
}
func (p *userRepo) Get(ctx context.Context, id string) (*po.User, error) {
	if len(id) == 0 {
		return nil, errors.New("instance id is required")
	}
	pm := &po.User{}
	if err := p.db.Model(pm).Where("instanceId = ?", id).First(pm).Error; err != nil {
		return nil, fmt.Errorf("%s get: %s", pm.TableName(), err.Error())
	}
	return pm, nil
}
func (p *userRepo) Delete(ctx context.Context, id string) error {
	if len(id) == 0 {
		return errors.New("instance id is required")
	}
	pm := &po.User{}
	if err := p.db.Where("instanceId =?", id).Delete(pm).Error; err != nil {
		return fmt.Errorf("%s get: %s", pm.TableName(), err.Error())
	}
	return nil
}
func (p *userRepo) List(ctx context.Context, cond *bo.PageQuery[po.User], username string) (*bo.SearchList[po.User], error) {
	if cond == nil {
		return nil, errors.New("condition is required")
	}
	var rsp bo.SearchList[po.User]
	md := p.db.Model(cond.Condition)
	if cond.Condition != nil {
		md = md.Where(cond.Condition)
	}
	if len(username) != 0 {
		md = md.Where("username like ?", "%"+username+"%")
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
