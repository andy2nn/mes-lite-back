package user

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCreds = errors.New("invalid username or password")

// ServiceInterface определяет методы, используемые handler’ом
type ServiceInterface interface {
	CreateUser(u *User, rawPassword string) error
	GetUser(id int64) (*User, error)
	ListUsers() ([]*User, error)
	UpdateUser(u *User, newPassword string) error
	DeleteUser(id int64) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(u *User, rawPassword string) error {
	if rawPassword == "" {
		return fmt.Errorf("password required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	return s.repo.Create(u)
}

func (s *Service) GetUser(id int64) (*User, error) {
	return s.repo.GetByID(id)
}

//TODO: обработка по имени (хэндлер тоже)

func (s *Service) ListUsers() ([]*User, error) {
	return s.repo.List()
}

// TODO: проверить и поправить
func (s *Service) UpdateUser(u *User, newPassword string) error {
	if newPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hash)
	} else {
		existing, err := s.repo.GetByID(u.ID)
		if err != nil {
			return err
		}
		u.Password = existing.Password
	}

	return s.repo.Update(u)
}

func (s *Service) DeleteUser(id int64) error {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(u)
}
