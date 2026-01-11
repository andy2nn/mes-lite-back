package permission

// ServiceInterface определяет методы, используемые handler’ом
type ServiceInterface interface {
	CreatePermission(p *Permission) error
	GetPermissionById(id int64) (*Permission, error)
	GetPermissionByName(name string) (*Permission, error)
	ListPermissions() ([]*Permission, error)
	UpdatePermission(p *Permission) error
	DeletePermission(id int64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePermission(p *Permission) error {
	return s.repo.Create(p)
}

func (s *Service) GetPermissionById(id int64) (*Permission, error) {
	return s.repo.GetPermissionById(id)
}

func (s *Service) GetPermissionByName(name string) (*Permission, error) {
	return s.repo.GetPermissionByName(name)
}

func (s *Service) ListPermissions() ([]*Permission, error) {
	return s.repo.List()
}

func (s *Service) UpdatePermission(p *Permission) error {
	return s.repo.Update(p)
}

func (s *Service) DeletePermission(id int64) error {
	p, err := s.repo.GetPermissionById(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(p)
}
