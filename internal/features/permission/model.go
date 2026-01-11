package permission

import (
	"time"
)

type Permission struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"unique;not null;size:100" json:"code"`
	Name        string    `gorm:"not null;size:255" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	Category    string    `gorm:"size:100" json:"category,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Permission) TableName() string {
	return "permissions"
}

type RolePermission struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID       int64     `gorm:"not null" json:"role_id"`
	PermissionID int64     `gorm:"not null" json:"permission_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

