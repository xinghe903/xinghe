package data

import (
	"comment/internal/biz/po"
	"comment/internal/biz/repo"
	"comment/internal/conf"
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	hashid "github.com/xinghe903/xinghe/pkg/distribute/id"
	"gorm.io/gorm"
)

var (
	ErrCreateComment = errors.New("create Comment failed")
	ErrUpdateComment = errors.New("update Comment failed")
	ErrQueryComment  = errors.New("query Comment failed")
	ErrParamIsEmpty  = errors.New("param is empty")
)

var _ repo.CommentRepo = &commentRepo{}

type commentRepo struct {
	db        *gorm.DB
	log       *log.Helper
	builderId *hashid.Snowflake
}

func NewCommentRepo(c *conf.Server, data *Data, logger log.Logger, id *hashid.Snowflake) repo.CommentRepo {
	return &commentRepo{
		db:        data.db,
		log:       log.NewHelper(logger),
		builderId: id,
	}
}

func (p *commentRepo) CreateComment(ctx context.Context, source *po.Comment) (string, error) {
	if len(source.InstanceId) == 0 {
		source.InstanceId = source.GenerateID(p.builderId.GenerateID())
	}
	if err := p.db.Create(source).Error; err != nil {
		return "", errors.Join(ErrCreateComment, err)
	}
	return source.InstanceId, nil
}
func (p *commentRepo) UpdateComment(ctx context.Context, source *po.Comment) error {
	if len(source.InstanceId) == 0 {
		return ErrParamIsEmpty
	}
	if err := p.db.Model(&po.Comment{}).Where("instanceId = ?", source.InstanceId).Updates(source).Error; err != nil {
		return errors.Join(ErrUpdateComment, err)
	}
	return nil
}
func (p *commentRepo) GetComment(ctx context.Context, id string) (*po.Comment, error) {
	if len(id) == 0 {
		return nil, ErrParamIsEmpty
	}
	pm := &po.Comment{}
	if err := p.db.Model(pm).Where("instanceId = ?", id).First(pm).Error; err != nil {
		return nil, errors.Join(ErrQueryComment, err)
	}
	return pm, nil
}
func (p *commentRepo) DeleteComment(ctx context.Context, id string) error {
	if len(id) == 0 {
		return ErrParamIsEmpty
	}
	if err := p.db.Where("instanceId =?", id).Delete(&po.Comment{}).Error; err != nil {
		return errors.Join(ErrQueryComment, err)
	}
	return nil
}
func (p *commentRepo) ListComment(ctx context.Context, cond *po.PageQuery[po.Comment]) (*po.SearchList[po.Comment], error) {
	var rsp po.SearchList[po.Comment]
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
		return nil, errors.Join(ErrQueryComment, err)
	}
	return &rsp, nil
}
