package po

import (
	"database/sql"

	"github.com/xinghe903/xinghe/pkg/distribute/hash"
	"gorm.io/gorm"
)

const (
	userTableName = "x_auth_user"
	userPrefixId  = "user-"
)

const (
	StatusUserActive   StatusUser = "active"
	StatusUserInactive StatusUser = "inactive"
	StatusUserDeleted  StatusUser = "deleted"
)

type StatusUser string

type User struct {
	gorm.Model
	InstanceId string         `json:"instanceId,omitempty" gorm:"unique;column:instanceId;type:varchar(40);not null"`
	Name       string         `json:"name,omitempty" gorm:"unique;column:name;type:varchar(40)"`
	NickName   string         `json:"nickname,omitempty" gorm:"column:nickname;type:varchar(40)"`
	Password   string         `json:"password,omitempty"  gorm:"column:password;type:varchar(255)"`
	Email      sql.NullString `json:"email,omitempty"  gorm:"unique;column:email;type:varchar(40)"`
	Phone      sql.NullString `json:"phone,omitempty"  gorm:"unique;column:phone;type:varchar(40)"`
	Status     StatusUser     `json:"status,omitempty"  gorm:"column:status;type:varchar(40)"`
}

func (u *User) TableName() string {
	return userTableName
}

func (u *User) GenerateID(seed uint64) string {
	return userPrefixId + hash.Base32Encode([]int32{int32(seed >> 32), int32(seed)})
}
