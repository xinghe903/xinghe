package data

import (
	"auth/internal/biz/repo"
	"auth/internal/conf"
	"auth/internal/data/po"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/xinghe903/xinghe/pkg/bo"
	hashid "github.com/xinghe903/xinghe/pkg/distribute/id"
	"gorm.io/gorm"
)

var _ repo.AuthRepo = &authRepo{}

type authRepo struct {
	db   *gorm.DB
	log  *log.Helper
	snow *hashid.Sonyflake
}

func NewAuthRepo(c *conf.Server, data *Data, logger log.Logger, id *hashid.Sonyflake) repo.AuthRepo {
	return &authRepo{
		db:   data.db,
		log:  log.NewHelper(logger),
		snow: id,
	}
}

func (p *authRepo) Create(ctx context.Context, source *po.Auth) (string, error) {
	if len(source.InstanceId) == 0 {
		source.InstanceId = source.GenerateID(p.snow.GenerateID())
	}
	if err := p.db.Create(source).Error; err != nil {
		return "", fmt.Errorf("%s create: %s", source.TableName(), err.Error())
	}
	return source.InstanceId, nil
}
func (p *authRepo) Update(ctx context.Context, source *po.Auth) error {
	if len(source.InstanceId) == 0 {
		return errors.New("instance id is required")
	}
	if err := p.db.Model(&po.Auth{}).Where("instanceId = ?", source.InstanceId).Updates(source).Error; err != nil {
		return fmt.Errorf("%s update: %s", source.TableName(), err.Error())
	}
	return nil
}
func (p *authRepo) UpdateByCode(ctx context.Context, source *po.Auth) error {
	if len(source.Code) == 0 {
		return errors.New("code id is required")
	}
	if err := p.db.Model(&po.Auth{}).Where("code = ?", source.Code).Updates(source).Error; err != nil {
		return fmt.Errorf("%s update: %s", source.TableName(), err.Error())
	}
	return nil
}
func (p *authRepo) Get(ctx context.Context, id string) (*po.Auth, error) {
	if len(id) == 0 {
		return nil, errors.New("instance id is required")
	}
	pm := &po.Auth{}
	if err := p.db.Model(pm).Where("instanceId = ?", id).First(pm).Error; err != nil {
		return nil, fmt.Errorf("%s get: %s", pm.TableName(), err.Error())
	}
	return pm, nil
}
func (p *authRepo) Delete(ctx context.Context, id string) error {
	if len(id) == 0 {
		return errors.New("instance id is required")
	}
	pm := &po.Auth{}
	if err := p.db.Where("instanceId =?", id).Delete(pm).Error; err != nil {
		return fmt.Errorf("%s get: %s", pm.TableName(), err.Error())
	}
	return nil
}
func (p *authRepo) List(ctx context.Context, cond *bo.PageQuery[po.Auth]) (*bo.SearchList[po.Auth], error) {
	if cond == nil {
		return nil, errors.New("condition is required")
	}
	var rsp bo.SearchList[po.Auth]
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

func (p *authRepo) GetByToken(ctx context.Context, token string) (*po.Auth, error) {
	if len(token) == 0 {
		return nil, errors.New("token is required")
	}
	pm := &po.Auth{}
	if err := p.db.Model(pm).Where("token = ?", token).First(pm).Error; err != nil {
		return nil, fmt.Errorf("%s get: %s", pm.TableName(), err.Error())
	}
	return pm, nil
}

func (p *authRepo) UpdateByToken(ctx context.Context, source *po.Auth) error {
	if len(source.Token.String) == 0 {
		return errors.New("Token is required")
	}
	if err := p.db.Model(&po.Auth{}).Where("token = ?", source.Token).Updates(source).Error; err != nil {
		return fmt.Errorf("%s update: %s", source.TableName(), err.Error())
	}
	return nil
}

func (p *authRepo) ClearExpiredToken(ctx context.Context, date time.Time) error {
	var source po.Auth
	if err := p.db.Model(&po.Auth{}).Where("expired_at <= ?", date).Updates(map[string]interface{}{
		"token":  nil,
		"status": po.StatusUserLogout,
	}).Error; err != nil {
		return fmt.Errorf("%s update: %s", source.TableName(), err.Error())
	}
	return nil
}
