package user

import (
	"mes-lite-back/internal/features/role"
	"time"
)

type User struct {
	ID        int64  `gorm:"primaryKey"`
	Username  string `gorm:"unique"`
	Password  string
	FullName  string
	RoleID    int64
	CreatedAt time.Time
	Email     string    `json:"email"`
	Role      role.Role `gorm:"foreignKey:RoleID"`
}

type RefreshToken struct {
	ID        int64 `gorm:"primaryKey"`
	UserID    int64
	Token     string `gorm:"unique"`
	ExpiresAt time.Time
	CreatedAt time.Time
}
