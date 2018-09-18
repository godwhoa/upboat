package api

import (
	"encoding/json"
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
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		R.Respond(w, R.Err(err))
		return
	}
	if err := req.Validate(); err != nil {
		R.Respond(w, R.Err(err))
		return
	}

	user, err := u.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		R.Respond(w, R.Err(err))
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
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		R.Respond(w, R.Err(err))
		return
	}
	if err := req.Validate(); err != nil {
		R.Respond(w, R.Err(err))
		return
	}

	user := &users.User{
		Email:    req.Email,
		Username: req.Username,
	}
	_, err := u.service.Register(ctx, user, req.Password)
	if err != nil {
		R.Respond(w, R.Err(err))
		return
	}
	R.Respond(w, R.Ok("Registered!"))
}

// Logout clears the session
func (u *UsersAPI) Logout(w http.ResponseWriter, r *http.Request) {
	session := u.sm.Load(r)
	err := session.Remove(w, "user_id")
	if err != nil {
		R.Respond(w, R.Err(err))
		u.log.Error("Error from session.Remove()", zap.Error(err))
		return
	}
	R.Respond(w, R.Ok("Logged out!"))
}
