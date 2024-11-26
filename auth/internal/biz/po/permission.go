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
	PermissionTypeUnknown   PermissionType = "unknown"
)

type PermissionType string

func NewPermissionType(s string) PermissionType {
	t := PermissionTypeUnknown
	switch s {
	case "menu":
		t = PermissionTypeMenu
	case "api":
		t = PermissionTypeAPI
	case "operation":
		t = PermissionTypeOperation
	}
	return t
}

const (
	PermissionHasChild = 1
	PermissionNoChild  = 0
)

type Permission struct {
	gorm.Model
	InstanceId  string         `json:"instanceId,omitempty" gorm:"unique;column:instanceId;type:varchar(40);not null"`
	SubjectType PermissionType `json:"subjectType,omitempty"  gorm:"column:subjectType;type:varchar(40)"`
	SubjectId   string         `json:"subjectId,omitempty"  gorm:"column:subjectId;type:varchar(40)"`
	RootId      string         `json:"rootId,omitempty"  gorm:"column:rootId;type:varchar(40)"`
	ParentId    string         `json:"parentId,omitempty"  gorm:"column:parentId;type:varchar(40)"`
	Name        string         `json:"name,omitempty"  gorm:"column:name;type:varchar(255)"`             // 权限名称
	Permission  string         `json:"permission,omitempty"  gorm:"column:permission;type:varchar(255)"` // 权限路径
	Sort        int32          `json:"sort,omitempty"  gorm:"column:sort;type:int;not null"`
	Children    int32          `json:"children,omitempty"  gorm:"column:children;type:int;not null"` // 是否有子节点 0: 没有子节点 1: 存在子节点
}

func (u *Permission) TableName() string {
	return permissionTableName
}

func (u *Permission) GenerateID(seed uint64) string {
	return permissionPrefixId + hash.Base32Encode([]int32{int32(seed >> 32), int32(seed)})
}
