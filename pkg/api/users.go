package api

import (
	"net/http"

	"github.com/alexedwards/scs"
	R "github.com/godwhoa/upboat/pkg/response"
	"github.com/godwhoa/upboat/pkg/users"
	"go.uber.org/zap"
)

type UsersAPI struct {
	service users.Service
	sm      *scs.Manager
	log     *zap.Logger
}

func NewUsersAPI(service users.Service, sm *scs.Manager, log *zap.Logger) *UsersAPI {
	return &UsersAPI{
		service: service,
		sm:      sm,
		log:     log,
	}
}

// Login handles login request, ok sucess sets a session cookie
func (u *UsersAPI) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &loginRequest{}
	if !DecodeValidate(w, r, req) {
		return
	}

	user, err := u.service.Login(ctx, req.Email, req.Password)
	if err == users.ErrInvalidCredentials {
		R.Respond(w, &R.Response{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if err != nil {
		R.Respond(w, R.InternalError())
		return
	}

	session := u.sm.Load(r)
	err = session.PutInt(w, "user_id", user.ID)
	if err != nil {
		R.Respond(w, R.InternalError())
		u.log.Error("Error from session.PutInt()", zap.Error(err))
		return
	}

	R.Respond(w, R.Ok("Logged in!"))
}

// Register handles user registration request
func (u *UsersAPI) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := &registerRequest{}
	if !DecodeValidate(w, r, req) {
		return
	}

	_, err := u.service.Register(ctx, &users.User{
		Email:    req.Email,
		Username: req.Username,
	}, req.Password)
	if err == users.ErrUserAlreadyExists {
		R.Respond(w, &R.Response{
			Code:    http.StatusConflict,
			Message: "User Already Exists",
			Data:    nil,
		})
		return
	}
	// log unexpected error
	if err != nil {
		R.Respond(w, R.InternalError())
		return
	}
	// Ok
	R.Respond(w, R.Ok("Registered!"))
}

// Logout clears the session
func (u *UsersAPI) Logout(w http.ResponseWriter, r *http.Request) {
	session := u.sm.Load(r)
	err := session.Remove(w, "user_id")
	if err != nil {
		R.Respond(w, R.InternalError())
		u.log.Error("Error from session.Remove()", zap.Error(err))
		return
	}
	R.Respond(w, R.Ok("Logged out!"))
}
