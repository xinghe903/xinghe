package po

import (
	"github.com/xinghe903/xinghe/pkg/distribute/hash"
	"gorm.io/gorm"
)

const (
	permissionTableName = "x_auth_permission"
	permissionPrefixId  = "perm-"
)

const (
	PermissionTypeMenu      PermissionType = "menu"
	PermissionTypeAPI       PermissionType = "api"
	PermissionTypeOperation PermissionType = "operation"
)

type PermissionType string

type Permission struct {
	gorm.Model
	InstanceId  string         `json:"instanceId,omitempty" gorm:"unique;column:instanceId;type:varchar(40);not null"`
	SubjectType PermissionType `json:"subjectType,omitempty"  gorm:"column:subjectType;type:varchar(40)"`
	SubjectId   string         `json:"subjectId,omitempty"  gorm:"column:subjectId;type:varchar(40)"`
	Permission  string         `json:"permission,omitempty"  gorm:"column:permission;type:varchar(255)"`
}

func (u *Permission) TableName() string {
	return permissionTableName
}

func (u *Permission) GenerateID(seed uint64) string {
	return permissionPrefixId + hash.Base32Encode([]int32{int32(seed >> 32), int32(seed)})
}
