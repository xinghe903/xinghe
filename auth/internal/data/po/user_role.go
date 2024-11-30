package po

import (
	"github.com/xinghe903/xinghe/pkg/distribute/hash"
	"gorm.io/gorm"
)

const (
	userRoleTableName = "x_auth_user_role"
	userRolePrefixId  = "userRole-"
)

type UserRole struct {
	gorm.Model
	InstanceId string `json:"instanceId,omitempty" gorm:"unique;column:instanceId;type:varchar(40);not null"`
	UserId     string `json:"userId,omitempty"  gorm:"column:userId;type:varchar(40)"`
	RoleId     string `json:"roleId,omitempty" gorm:"column:roleId;type:varchar(40)"`
}

func (u *UserRole) TableName() string {
	return userRoleTableName
}

func (u *UserRole) GenerateID(seed uint64) string {
	return userRolePrefixId + hash.Base32Encode([]int32{int32(seed >> 32), int32(seed)})
}
