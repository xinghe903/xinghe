package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	hashid "github.com/xinghe903/xinghe/pkg/distribute/id"
	"gorm.io/gorm"
)

type SearchList[T any] struct {
	Total int64
	Data  []*T
}

type PageQuery[T any] struct {
	PageNum   int32
	PageSize  int32
	Condition *T
	Sort      map[string]string
}

type BaseOrm interface {
	GetInstanceId() string
	SetInstanceId(string)
	GenerateID(id int64) string
	TableName() string
}

// 这里T为指针类型
type BaseRepo[T BaseOrm] struct {
	db        *gorm.DB
	log       *log.Helper
	builderId *hashid.Snowflake
}

var (
	ErrInstancesEmpty = errors.New("instanceId is empty")
	ErrSourceIsNil    = errors.New("source is nil")
	ErrCondIsNil      = errors.New("cond is nil")
)

func (p *BaseRepo[T]) Create(ctx context.Context, source T) (string, error) {
	if len(source.GetInstanceId()) == 0 {
		source.SetInstanceId(source.GenerateID(p.builderId.GenerateID()))
	}
	if err := p.db.Create(source).Error; err != nil {
		return "", errors.Join(fmt.Errorf("create %s ", source.TableName()), err)
	}
	return source.GetInstanceId(), nil
}
func (p *BaseRepo[T]) Update(ctx context.Context, source T) error {
	if len(source.GetInstanceId()) == 0 {
		return ErrInstancesEmpty
	}
	if err := p.db.Model(new(T)).Where("instanceId = ?", source.GetInstanceId()).Updates(source).Error; err != nil {
		return errors.Join(fmt.Errorf("update %s ", source.TableName()), err)
	}
	return nil
}
func (p *BaseRepo[T]) Get(ctx context.Context, id string) (*T, error) {
	var pm T
	if len(id) == 0 {
		return &pm, ErrInstancesEmpty
	}
	if err := p.db.Model(&pm).Where("instanceId = ?", id).First(&pm).Error; err != nil {
		return &pm, errors.Join(fmt.Errorf("get %s ", pm.TableName()), err)
	}
	return &pm, nil
}
func (p *BaseRepo[T]) Delete(ctx context.Context, id string) error {
	if len(id) == 0 {
		return ErrInstancesEmpty
	}
	var pm T
	if err := p.db.Where("instanceId =?", id).Delete(&pm).Error; err != nil {
		return errors.Join(fmt.Errorf("delete %s ", pm.TableName()), err)
	}
	return nil
}
func (p *BaseRepo[T]) List(ctx context.Context, cond *PageQuery[T]) (*SearchList[T], error) {
	if cond == nil {
		return nil, ErrCondIsNil
	}
	var rsp SearchList[T]
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
		return nil, errors.Join(ErrInstancesEmpty, err)
	}
	return &rsp, nil
}
