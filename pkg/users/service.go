package users

import (
	"golang.org/x/crypto/bcrypt"
)

// Service implements UserService interface
type service struct {
	repo Repository
}

// NewService is a constructor for user.Service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(u *User, password string) (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Hash = string(hashed)

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return s.repo.FindByEmail(u.Email)
}

func (s *service) Login(email string, password string) (*User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
