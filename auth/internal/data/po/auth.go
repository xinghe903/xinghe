package po

import (
	"database/sql"
	"time"

	"github.com/xinghe903/xinghe/pkg/distribute/hash"
	"gorm.io/gorm"
)

const (
	authTableName = "x_auth_auth"
	authPrefixId  = "auth-"
)

const (
	StatusUserLogin  StatusAuth = "login"
	StatusUserLogout StatusAuth = "logout"
)

type StatusAuth string

type Auth struct {
	gorm.Model
	InstanceId string         `json:"instanceId,omitempty" gorm:"unique;column:instanceId;type:varchar(40);not null"`
	Name       string         `json:"name,omitempty" gorm:"unique;column:name;type:varchar(40)"`
	NickName   string         `json:"nickname,omitempty" gorm:"column:nickname;type:varchar(40)"`
	Code       string         `json:"code,omitempty" gorm:"unique;column:code;type:varchar(40)"`
	Token      sql.NullString `json:"token,omitempty" gorm:"unique;column:token;type:varchar(40)"`
	ExpiredAt  time.Time      `json:"expired_at,omitempty" gorm:"column:expired_at;type:datetime"` // 到期时间
	Status     StatusAuth     `json:"status,omitempty"  gorm:"column:status;type:varchar(40)"`
}

func (u *Auth) TableName() string {
	return authTableName
}

func (u *Auth) GenerateID(seed uint64) string {
	return authPrefixId + hash.Base32Encode([]int32{int32(seed >> 32), int32(seed)})
}
