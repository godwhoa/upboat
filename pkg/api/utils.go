package api

import (
	"encoding/json"
	"net/http"

	v "github.com/go-ozzo/ozzo-validation"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Respond responds with JSON encoded `response`
func Respond(w http.ResponseWriter, r *response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)
	if json.NewEncoder(w).Encode(r) != nil {
		http.Error(w, "Error occured while encoding JSON", http.StatusInternalServerError)
	}
}

// DecodeValidate decodes JSON, validates it and responds accordingly
func DecodeValidate(w http.ResponseWriter, r *http.Request, obj v.Validatable) (ok bool) {
	if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
		Respond(w, JSONError())
		return false
	}
	if err := obj.Validate(); err != nil {
		Respond(w, ValidationError(err))
		return false
	}
	return true
}

// Helper functions for syntatic sugar
func Ok(msg string) *response {
	return &response{
		Code:    http.StatusOK,
		Message: msg,
		Data:    nil,
	}
}

func OkData(msg string, data interface{}) *response {
	return &response{
		Code:    http.StatusOK,
		Message: msg,
		Data:    data,
	}
}

func Created(msg string, data interface{}) *response {
	return &response{
		Code:    http.StatusCreated,
		Message: msg,
		Data:    data,
	}
}

func NotFound(msg string) *response {
	return &response{
		Code:    http.StatusNotFound,
		Message: msg,
		Data:    nil,
	}
}

func Unauthorized(msg string) *response {
	return &response{
		Code:    http.StatusUnauthorized,
		Message: msg,
		Data:    nil,
	}
}

func InternalError() *response {
	return &response{
		Code:    http.StatusInternalServerError,
		Message: "Internal Error",
		Data:    nil,
	}
}

func JSONError() *response {
	return &response{
		Code:    http.StatusBadRequest,
		Message: "Invalid JSON",
		Data:    nil,
	}
}

func ValidationError(err error) *response {
	return &response{
		Code:    http.StatusBadRequest,
		Message: "Validation Error",
		Data:    err,
	}
}
