package api

import (
	v "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l loginRequest) Validate() error {
	return v.ValidateStruct(&l,
		v.Field(&l.Email, v.Required, is.Email),
		v.Field(&l.Password, v.Required),
	)
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r registerRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Username, v.Required),
		v.Field(&r.Email, v.Required, is.Email),
		v.Field(&r.Password, v.Required),
	)
}

type createRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (r createRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Body, v.Required),
		v.Field(&r.Title, v.Required, v.Length(1, 200)),
	)
}

type updateRequest struct {
	createRequest
}

type voteRequest struct {
	Delta int `json:"delta"`
}

func (r voteRequest) Validate() error {
	return v.ValidateStruct(&r,
		v.Field(&r.Delta, v.Required, v.In(-1, +1)),
	)
}
