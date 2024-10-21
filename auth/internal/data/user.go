package data

import (
	"auth/internal/biz/po"
	"auth/internal/biz/repo"
	"auth/internal/conf"
	"context"
	"errors"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	hashid "github.com/xinghe903/xinghe/pkg/distribute/id"
	"gorm.io/gorm"
)

var (
	ErrCreateUser   = errors.New("create user failed")
	ErrUpdateUser   = errors.New("update user failed")
	ErrQueryUser    = errors.New("query user failed")
	ErrParamIsEmpty = errors.New("param is empty")
)

var _ repo.UserRepo = &userRepo{}

type userRepo struct {
	db        *gorm.DB
	log       *log.Helper
	builderId *hashid.Snowflake
}

func NewUserRepo(c *conf.Server, data *Data, logger log.Logger) repo.UserRepo {
	nodeId, err := strconv.Atoi(c.GetNodeId())
	if err != nil {
		panic(err.Error())
	}
	return &userRepo{
		db:        data.db,
		log:       log.NewHelper(logger),
		builderId: hashid.NewSnowflake(int64(nodeId)),
	}
}

func (p *userRepo) CreateUser(ctx context.Context, source *po.User) (string, error) {
	if len(source.InstanceId) == 0 {
		source.InstanceId = source.GenerateID(p.builderId.GenerateID())
	}
	if err := p.db.Create(source).Error; err != nil {
		return "", errors.Join(ErrCreateUser, err)
	}
	return source.InstanceId, nil
}
func (p *userRepo) UpdateUser(ctx context.Context, source *po.User) error {
	if len(source.InstanceId) == 0 {
		return ErrParamIsEmpty
	}
	if err := p.db.Model(&po.User{}).Where("instanceId = ?", source.InstanceId).Updates(source).Error; err != nil {
		return errors.Join(ErrUpdateUser, err)
	}
	return nil
}
func (p *userRepo) GetUser(ctx context.Context, id string) (*po.User, error) {
	if len(id) == 0 {
		return nil, ErrParamIsEmpty
	}
	pm := &po.User{}
	if err := p.db.Model(pm).Where("instanceId = ?", id).First(pm).Error; err != nil {
		return nil, errors.Join(ErrQueryUser, err)
	}
	return pm, nil
}
func (p *userRepo) DeleteUser(ctx context.Context, id string) error {
	if len(id) == 0 {
		return ErrParamIsEmpty
	}
	if err := p.db.Where("instanceId =?", id).Delete(&po.User{}).Error; err != nil {
		return errors.Join(ErrQueryUser, err)
	}
	return nil
}
func (p *userRepo) ListUser(ctx context.Context, cond *po.PageQuery[po.User]) (*po.SearchList[po.User], error) {
	var rsp po.SearchList[po.User]
	pageSize, pageNum := cond.PageSize, cond.PageNum
	if pageSize == 0 {
		pageSize = 10
	}
	if pageNum == 0 {
		pageNum = 1
	}

	if err := p.db.Model(cond.Condition).
		Where(cond.Condition).
		Count(&rsp.Total).
		Offset(int(pageNum-1) * int(pageSize)).
		Limit(int(pageSize)).
		Find(&rsp.Data).Error; err != nil {
		return nil, errors.Join(ErrQueryUser, err)
	}
	return &rsp, nil
}
