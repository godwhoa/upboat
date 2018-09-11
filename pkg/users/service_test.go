package users

import (
	"testing"

	qt "github.com/frankban/quicktest"
	"golang.org/x/crypto/bcrypt"
)

type mockRepo struct {
	createerr     bool
	finderr       bool
	createcalled  bool
	findcalled    bool
	byemailcalled bool
	u             *User
}

func (r *mockRepo) Create(user *User) error {
	if r.createerr {
		return ErrUserAlreadyExists
	}
	r.createcalled = true
	r.u = user
	return nil

}
func (r *mockRepo) Find(id int) (*User, error) {
	r.findcalled = true
	if r.finderr {
		return nil, ErrUserNotFound
	}
	return r.u, nil
}
func (r *mockRepo) FindByEmail(email string) (*User, error) {
	r.byemailcalled = true
	if r.finderr {
		return nil, ErrUserNotFound
	}
	return r.u, nil
}

func TestNewService(t *testing.T) {
	c := qt.New(t)
	service := NewService(&mockRepo{})
	c.Assert(service, qt.Not(qt.IsNil))
}

// Ensure Register hashes the password
func TestService_Register_EnsureHashing(t *testing.T) {
	c := qt.New(t)
	service := NewService(&mockRepo{})

	user, err := service.Register(&User{Username: "blah", Email: "blah@blah.com"}, "password")
	c.Assert(user, qt.Not(qt.IsNil))
	c.Assert(err, qt.IsNil)

	c.Assert(user.Hash, qt.Not(qt.Equals), "password")

	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte("password"))
	c.Assert(err, qt.IsNil)
}

// Ensure Register calls UserRepository.Create()
func TestService_Register_EnsureCreate(t *testing.T) {
	c := qt.New(t)
	repo := &mockRepo{}
	service := NewService(repo)
	service.Register(&User{Username: "blah", Email: "blah@blah.com"}, "password")
	c.Assert(repo.createcalled, qt.Equals, true)
}

// Ensure it propogades ErrUserAlreadyExists
func TestService_Register_AlreadyExists(t *testing.T) {
	c := qt.New(t)
	repo := &mockRepo{createerr: true}
	service := NewService(repo)
	user, err := service.Register(&User{Username: "blah", Email: "blah@blah.com"}, "password")
	c.Assert(err, qt.Equals, ErrUserAlreadyExists)
	c.Assert(user, qt.IsNil)
}

// "OK" case for register
func TestService_Register_OK(t *testing.T) {
	c := qt.New(t)
	repo := &mockRepo{}
	service := NewService(repo)
	user, err := service.Register(&User{Username: "blah", Email: "blah@blah.com"}, "password")
	c.Assert(err, qt.IsNil)
	c.Assert(user, qt.Not(qt.IsNil))
}

func TestService_Login_OK(t *testing.T) {
	c := qt.New(t)
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	u := &User{
		ID:       0,
		Username: "blah",
		Email:    "blah@blah.com",
		Hash:     string(hash),
	}
	service := NewService(&mockRepo{u: u})
	user, err := service.Login("blah@blah.com", "password")
	c.Assert(err, qt.IsNil)
	c.Assert(user, qt.Not(qt.IsNil))
}

func TestService_Login_NotFound(t *testing.T) {
	c := qt.New(t)
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	u := &User{
		ID:       0,
		Username: "blah",
		Email:    "blah@blah.com",
		Hash:     string(hash),
	}
	service := NewService(&mockRepo{u: u, finderr: true})
	user, err := service.Login("apple@kak.com", "password")
	c.Assert(user, qt.IsNil)
	c.Assert(err, qt.Equals, ErrUserNotFound)
}

func TestService_LoginInvalid(t *testing.T) {
	c := qt.New(t)
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	u := &User{
		ID:       0,
		Username: "blah",
		Email:    "blah@blah.com",
		Hash:     string(hash),
	}
	service := NewService(&mockRepo{u: u})
	user, err := service.Login("blah@blah.com", "passwordddd")
	c.Assert(user, qt.IsNil)
	c.Assert(err, qt.Equals, ErrInvalidCredentials)
}
