package users

import (
	"context"

	"github.com/godwhoa/upboat/pkg/errors"
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

func (s *service) Register(ctx context.Context, u *User, password string) (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.E(errors.Internal, errors.Op("bcrypt.GenerateFromPassword"), err)
	}
	u.Hash = string(hashed)

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return s.repo.FindByEmail(ctx, u.Email)
}

func (s *service) Login(ctx context.Context, email string, password string) (*User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
