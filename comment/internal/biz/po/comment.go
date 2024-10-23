package po

import (
	"github.com/xinghe903/xinghe/pkg/distribute/hash"
	"gorm.io/gorm"
)

const (
	commentTableName = "x_comment_base"
	commentPrefixId  = "auth-"
)

const (
	CommentSubjectTypeVideo    CommentSubjectType = "video"    // 视频
	CommentSubjectTypeArticle  CommentSubjectType = "article"  // 文章
	CommentSubjectTypeMusic    CommentSubjectType = "music"    // 音乐
	CommentSubjectTypeBook     CommentSubjectType = "book"     //  书
	CommentSubjectTypeOffering CommentSubjectType = "offering" // 商品
)

type CommentSubjectType string

type Comment struct {
	gorm.Model
	InstanceId  string             `json:"instanceId,omitempty" gorm:"unique;column:instanceId;type:varchar(40);not null"`
	UserId      string             `json:"nickname,omitempty" gorm:"column:nickname;type:varchar(40)"`       //  评论用户ID
	ReplyUserId string             `json:"replyUserId,omitempty" gorm:"column:replyUserId;type:varchar(40)"` // 回复用户ID
	SubjectId   string             `json:"subjectId,omitempty" gorm:"column:subjectId;type:varchar(40)"`     // 评论对象ID
	SubjectType CommentSubjectType `json:"subjectType,omitempty" gorm:"column:subjectType;type:varchar(40)"` // 评论对象类型
	RootId      string             `json:"rootId,omitempty" gorm:"column:rootId;type:varchar(40)"`           //  顶级评论ID
	ParentId    string             `json:"parentId,omitempty" gorm:"column:parentId;type:varchar(40)"`       //  父级评论ID
	ReplyCount  int32              `json:"replyCount,omitempty" gorm:"column:replyCount;type:int(10)"`       //  回复数
	LikeCount   int32              `json:"likeCount,omitempty" gorm:"column:likeCount;type:int(10)"`         // 点赞数
	Content     string             `json:"content,omitempty" gorm:"column:content;type:TEXT"`                //  评论内容
}

func (c *Comment) TableName() string {
	return commentTableName
}

func (c *Comment) GenerateID(seed int64) string {
	return hash.GetHashId(seed, commentTableName)
}
