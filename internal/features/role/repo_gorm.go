package role

import (
	"mes-lite-back/internal/features/permission"

	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(role *Role, permissionIDs []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(role).Error; err != nil {
			return err
		}

		if len(permissionIDs) > 0 {
			for _, permID := range permissionIDs {
				rp := &permission.RolePermission{
					RoleID:       role.ID,
					PermissionID: permID,
				}
				if err := tx.Create(rp).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *GormRepository) GetRole(id int64) (*Role, error) {
	var role Role

	err := r.db.
		Preload("Role.Permissions").
		First(&role, id).
		Error

	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *GormRepository) List() ([]*Role, error) {
	var roles []*Role
	return roles, r.db.Find(&roles).Error
}

func (r *GormRepository) Update(role *Role) error {
	return r.db.Create(role).Error
}

func (r *GormRepository) Delete(role *Role) error {
	return r.db.Delete(role).Error
}

func (r *GormRepository) GetByRole(name string) (*Role, error) {
	var role Role
	err := r.db.
		Preload("Role.Permissions").
		Where("name = ?", name).
		First(&role).
		Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *GormRepository) UpdatePermissions(roleID int64, permissionIDs []int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		var role Role
		if err := tx.First(&role, roleID).Error; err != nil {
			return err
		}

		if err := tx.Where("role_id = ?", roleID).
			Delete(&permission.RolePermission{}).Error; err != nil {
			return err
		}

		for _, permID := range permissionIDs {

			var perm permission.Permission
			if err := tx.First(&perm, permID).Error; err != nil {
				return err
			}

			rp := &permission.RolePermission{
				RoleID:       roleID,
				PermissionID: permID,
			}
			if err := tx.Create(rp).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
