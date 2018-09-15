package response

import (
	"encoding/json"
	"net/http"
)

// Response is a standardized JSON response container
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Respond responds with JSON encoded `response`
func Respond(w http.ResponseWriter, r *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)
	if json.NewEncoder(w).Encode(r) != nil {
		http.Error(w, "Error occured while encoding JSON", http.StatusInternalServerError)
	}
}

// Helper functions for syntatic sugar
func Ok(msg string) *Response {
	return &Response{
		Code:    http.StatusOK,
		Message: msg,
		Data:    nil,
	}
}

func OkData(msg string, data interface{}) *Response {
	return &Response{
		Code:    http.StatusOK,
		Message: msg,
		Data:    data,
	}
}

func Created(msg string, data interface{}) *Response {
	return &Response{
		Code:    http.StatusCreated,
		Message: msg,
		Data:    data,
	}
}

func NotFound(msg string) *Response {
	return &Response{
		Code:    http.StatusNotFound,
		Message: msg,
		Data:    nil,
	}
}

func Unauthorized(msg string) *Response {
	return &Response{
		Code:    http.StatusUnauthorized,
		Message: msg,
		Data:    nil,
	}
}

func InternalError() *Response {
	return &Response{
		Code:    http.StatusInternalServerError,
		Message: "Internal Error",
		Data:    nil,
	}
}

func JSONError() *Response {
	return &Response{
		Code:    http.StatusBadRequest,
		Message: "Invalid JSON",
		Data:    nil,
	}
}

func ValidationError(err error) *Response {
	return &Response{
		Code:    http.StatusBadRequest,
		Message: "Validation Error",
		Data:    err,
	}
}
