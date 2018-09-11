package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexedwards/scs"
	qt "github.com/frankban/quicktest"
	"github.com/godwhoa/upboat/pkg/upboat"
	"go.uber.org/zap"
)

type mockService struct {
	loginerr    bool
	regerr      bool
	logincalled bool
}

func (s *mockService) Register(u *upboat.User, password string) (*upboat.User, error) {
	if s.regerr {
		return nil, upboat.ErrUserAlreadyExists
	}
	return u, nil
}

func (s *mockService) Login(email string, password string) (*upboat.User, error) {
	s.logincalled = true
	if s.loginerr {
		return nil, upboat.ErrInvalidCredentials
	}
	return &upboat.User{
		ID:       0,
		Username: "blah",
		Email:    "blah@blah.com",
		Hash:     "kjadkjglkjlkj",
	}, nil
}

func deps() (*zap.Logger, *scs.Manager, *mockService) {
	log, _ := zap.NewProduction()
	sm := scs.NewCookieManager("ksajkjgfkjkjkjkjkjijijkjdkljfkl")
	return log, sm, &mockService{}
}

func post(endpoint, payload string) (*http.Request, error) {
	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer([]byte(payload)))
	return req, err
}

func TestLogin_Invalid_Input(t *testing.T) {
	c := qt.New(t)
	req, err := post("/api/login", `{"email":"blah@blah.com"}`)
	c.Assert(err, qt.IsNil)
	c.Assert(req, qt.Not(qt.IsNil))

	rr := httptest.NewRecorder()

	log, sm, service := deps()

	http.HandlerFunc(Login(service, sm, log)).
		ServeHTTP(rr, req)
	c.Assert(service.logincalled, qt.Equals, false)
	c.Assert(rr.Code, qt.Equals, http.StatusBadRequest)
}

func TestLoginOK(t *testing.T) {
	c := qt.New(t)
	req, err := post("/api/login", `{"email":"blah@blah.com", "password":"password"}`)
	c.Assert(err, qt.IsNil)
	c.Assert(req, qt.Not(qt.IsNil))

	rr := httptest.NewRecorder()

	log, sm, service := deps()

	http.HandlerFunc(Login(service, sm, log)).
		ServeHTTP(rr, req)
	c.Assert(rr.Code, qt.Equals, http.StatusOK)
}

func TestLogout(t *testing.T) {
	c := qt.New(t)
	req, err := post("/api/login", `{"email":"blah@blah.com", "password":"password"}`)
	c.Assert(err, qt.IsNil)
	c.Assert(req, qt.Not(qt.IsNil))

	rr := httptest.NewRecorder()

	log, sm, service := deps()

	http.HandlerFunc(Login(service, sm, log)).
		ServeHTTP(rr, req)
	c.Assert(rr.Code, qt.Equals, http.StatusOK)

	req, err = http.NewRequest("GET", "/api/logout", nil)
	c.Assert(err, qt.IsNil)
	c.Assert(req, qt.Not(qt.IsNil))
	rr = httptest.NewRecorder()

	http.HandlerFunc(Logout(sm, log)).
		ServeHTTP(rr, req)
	c.Assert(rr.Code, qt.Equals, http.StatusOK)
}

func TestRegisterAlreadyExists(t *testing.T) {
	c := qt.New(t)
	payload := []byte(`{"email":"blah@blah.com", "username": "blah", "password":"password"}`)
	req, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(payload))
	c.Assert(err, qt.IsNil)

	rr := httptest.NewRecorder()

	log, _ := zap.NewProduction()
	service := &mockService{regerr: true}

	http.HandlerFunc(Register(service, log)).
		ServeHTTP(rr, req)
	c.Assert(rr.Code, qt.Equals, http.StatusBadRequest)
}

func TestRegisterOK(t *testing.T) {
	c := qt.New(t)
	req, err := post("/api/register", `{"email":"blah@blah.com", "username": "blah", "password":"password"}`)
	c.Assert(err, qt.IsNil)

	rr := httptest.NewRecorder()

	log, _ := zap.NewProduction()
	service := &mockService{}

	http.HandlerFunc(Register(service, log)).
		ServeHTTP(rr, req)
	c.Assert(rr.Code, qt.Equals, http.StatusOK)
}
