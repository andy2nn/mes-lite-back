package role

type Repository interface {
	Create(r *Role, listPerIds []int64) error
	Update(r *Role) error
	Delete(r *Role) error

	GetRole(id int64) (*Role, error)
	GetByRole(name string) (*Role, error)
	List() ([]*Role, error)
	UpdatePermissions(roleID int64, permissionIDs []int64) error
}
