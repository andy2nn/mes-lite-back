package role

// ServiceInterface определяет методы, используемые handler’ом
type ServiceInterface interface {
	CreateRole(r *Role, listPerIds []int64) error
	GetRole(id int64) (*Role, error)
	ListRoles() ([]*Role, error)
	UpdateRole(r *Role) error
	DeleteRole(id int64) error
	UpdatePermissions(roleID int64, permissionIDs []int64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateRole(r *Role, listPerIds []int64) error {
	return s.repo.Create(r, listPerIds)
}

func (s *Service) GetRole(id int64) (*Role, error) {
	return s.repo.GetRole(id)
}

func (s *Service) GetRoleByName(name string) (*Role, error) {
	return s.repo.GetByRole(name)
}

func (s *Service) ListRoles() ([]*Role, error) {
	return s.repo.List()
}

func (s *Service) UpdateRole(r *Role) error {
	return s.repo.Update(r)
}

func (s *Service) DeleteRole(id int64) error {
	r, err := s.repo.GetRole(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(r)
}

func (s *Service) UpdatePermissions(roleID int64, permissionIDs []int64) error {
	return s.repo.UpdatePermissions(roleID, permissionIDs)
}
