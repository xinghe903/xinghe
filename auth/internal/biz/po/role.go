package po

import (
	"github.com/xinghe903/xinghe/pkg/distribute/hash"
	"gorm.io/gorm"
)

const (
	roleTableName = "x_auth_role"
	rolePrefixId  = "role-"
)

const (
	StatusRoleActive   StatusRole = "active"
	StatusRoleInactive StatusRole = "inactive"
)

type StatusRole string

type Role struct {
	gorm.Model
	InstanceId string     `json:"instanceId,omitempty" gorm:"unique;column:instanceId;type:varchar(40);not null"`
	Name       string     `json:"name,omitempty" gorm:"unique;column:name;type:varchar(40)"`
	Status     StatusRole `json:"status,omitempty"  gorm:"column:status;type:varchar(40)"`
}

func (u *Role) TableName() string {
	return roleTableName
}

func (u *Role) GenerateID(seed int64) string {
	return hash.GetHashId(seed, rolePrefixId)
}
