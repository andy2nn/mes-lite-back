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

// UserPermissionCheck - структура для проверки прав пользователя
type UserPermissionCheck struct {
	UserID         int64  `json:"user_id"`
	PermissionCode string `json:"permission_code"`
}

// PermissionAssignRequest - запрос на назначение разрешений
type PermissionAssignRequest struct {
	RoleID        int64   `json:"role_id" validate:"required"`
	PermissionIDs []int64 `json:"permission_ids"`
}

// PermissionWithAssignment - разрешение с информацией о назначении роли
type PermissionWithAssignment struct {
	Permission
	Assigned bool `json:"assigned"`
}
