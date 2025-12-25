package user

import "gorm.io/gorm"

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(u *User) error {
	return r.db.Create(u).Error
}

func (r *GormRepository) GetByID(id int64) (*User, error) {
	var u User

	err := r.db.
		Preload("Role").
		Preload("Role.Permissions").
		First(&u, id).
		Error

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *GormRepository) List() ([]*User, error) {
	var users []*User
	return users, r.db.Find(&users).Error
}

func (r *GormRepository) Update(u *User) error {
	return r.db.Save(u).Error
}

func (r *GormRepository) Delete(u *User) error {
	return r.db.Delete(u).Error
}

func (r *GormRepository) GetByUsername(username string) (*User, error) {
	var u User

	err := r.db.
		Preload("Role").
		Preload("Role.Permissions").
		Where("username = ?", username).
		First(&u).
		Error

	if err != nil {
		return nil, err
	}

	return &u, nil
}
