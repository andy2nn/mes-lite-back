package permission

type Repository interface {
	Create(p *Permission) error
	Update(p *Permission) error
	Delete(p *Permission) error

	GetPermissionById(id int64) (*Permission, error)
	GetPermissionByName(name string) (*Permission, error)
	List() ([]*Permission, error)
}
