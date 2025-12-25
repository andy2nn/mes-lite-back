package user

type Repository interface {
	Create(u *User) error
	Update(u *User) error
	Delete(u *User) error

	GetByID(id int64) (*User, error)
	GetByUsername(username string) (*User, error)
	List() ([]*User, error)
}
