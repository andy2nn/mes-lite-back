package user

import (
	"errors"

	"gorm.io/gorm"
)

type refreshTokenRepo struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) Get(token string) (*RefreshToken, error) {
	var rt RefreshToken

	err := r.db.
		Where("token = ?", token).
		First(&rt).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // <- НЕ найдено = не ошибка
		}
		return nil, err
	}

	return &rt, nil
}

func (r *refreshTokenRepo) Save(token *RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepo) Delete(token string) error {
	return r.db.
		Where("token = ?", token).
		Delete(&RefreshToken{}).
		Error
}
