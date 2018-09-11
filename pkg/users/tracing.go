package users

import (
	"context"

	"go.opencensus.io/trace"
)

// Tracing is middleware that provides tracing
func Tracing(service Service) Service {
	return &tracingMiddleware{service}
}

// tracingMiddleware wraps around `Service` to log useful errors
type tracingMiddleware struct {
	service Service
}

func (m *tracingMiddleware) Register(ctx context.Context, user *User, password string) (*User, error) {
	ctx, span := trace.StartSpan(ctx, "posts.Service.Register")
	defer span.End()
	return m.service.Register(ctx, user, password)
}

func (m *tracingMiddleware) Login(ctx context.Context, email string, password string) (*User, error) {
	ctx, span := trace.StartSpan(ctx, "users.Service.Login")
	defer span.End()
	return m.service.Login(ctx, email, password)
}
