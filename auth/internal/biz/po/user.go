package po

import (
	"github.com/xinghe903/xinghe/pkg/distribute/hash"
	"gorm.io/gorm"
)

const (
	userTableName = "x_auth_user"
	userPrefixId  = "auth-"
)

const (
	StatusUserActive   StatusUser = "active"
	StatusUserInactive StatusUser = "inactive"
	StatusUserDeleted  StatusUser = "deleted"
)

type StatusUser string

type User struct {
	gorm.Model
	InstanceId string     `json:"instanceId,omitempty" gorm:"unique;column:instanceId;type:varchar(40);not null"`
	Name       string     `json:"name,omitempty" gorm:"column:name;type:varchar(40)"`
	NickName   string     `json:"nickname,omitempty" gorm:"column:nickname;type:varchar(40)"`
	Password   string     `json:"password,omitempty"  gorm:"column:password;type:varchar(40)"`
	Email      string     `json:"email,omitempty"  gorm:"column:email;type:varchar(40)"`
	Phone      string     `json:"phone,omitempty"  gorm:"column:phone;type:varchar(40)"`
	Status     StatusUser `json:"status,omitempty"  gorm:"column:status;type:varchar(40)"`
}

func (u *User) TableName() string {
	return userTableName
}

func (u *User) GenerateID(seed int64) string {
	return hash.GetHashId(seed, userPrefixId)
}
