package permission

import "gorm.io/gorm"

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(p *Permission) error {
	return r.db.Create(p).Error
}

func (r *GormRepository) Update(p *Permission) error {
	return r.db.Save(p).Error
}

func (r *GormRepository) Delete(p *Permission) error {
	return r.db.Delete(p).Error
}

func (r *GormRepository) GetPermissionById(id int64) (*Permission, error) {
	var perm Permission
	err := r.db.First(&perm, id).Error
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *GormRepository) GetPermissionByName(name string) (*Permission, error) {
	var perm Permission
	err := r.db.Where("name = ?", name).First(&perm).Error
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *GormRepository) List() ([]*Permission, error) {
	var perms []*Permission
	return perms, r.db.Find(&perms).Error
}
