package role

import (
	"mes-lite-back/internal/features/permission"
)

type Role struct {
	ID   int64  `gorm:"primaryKey"`
	Name string `gorm:"unique"`

	Permissions []permission.Permission `gorm:"many2many:role_permissions;"`
}

func (Role) TableName() string {
	return "roles"
}
