package po

import (
	"github.com/xinghe903/xinghe/pkg/distribute/hash"
	"gorm.io/gorm"
)

const (
	rolePerTableName = "x_auth_role_permission"
	rolePerPrefixId  = "roleper-"
)

type RolePermission struct {
	gorm.Model
	InstanceId   string `json:"instanceId,omitempty" gorm:"unique;column:instanceId;type:varchar(40);not null"`
	RoleId       string `json:"roleId,omitempty" gorm:"column:roleId;type:varchar(40)"`
	PermissionId string `json:"permissionId,omitempty"  gorm:"column:permissionId;type:varchar(40)"`
}

func (u *RolePermission) TableName() string {
	return rolePerTableName
}

func (u *RolePermission) GenerateID(seed int64) string {
	return hash.GetHashId(seed, rolePerPrefixId)
}
