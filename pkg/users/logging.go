package users

import (
	"context"

	"go.uber.org/zap"
)

// Middleware is anything that wraps around a `Service`
type Middleware func(Service) Service

// Logging is a middleware that provides logging to Service
func Logging(log *zap.Logger) Middleware {
	return func(service Service) Service {
		return &loggingMiddleware{service, log}
	}
}

// Chain lets you chain multiple middleware
func Chain(service Service, middlewares ...Middleware) Service {
	if len(middlewares) == 0 {
		return service
	}

	// Wrap the first middleware with the service
	s := middlewares[len(middlewares)-1](service)
	// Wrap that with the rest of the middleware chain
	for i := len(middlewares) - 2; i >= 0; i-- {
		s = middlewares[i](s)
	}

	return s
}

// loggingMiddleware wraps around `Service` to log useful errors
type loggingMiddleware struct {
	service Service
	log     *zap.Logger
}

func (m *loggingMiddleware) Register(ctx context.Context, user *User, password string) (u *User, err error) {
	u, err = m.service.Register(ctx, user, password)
	if err != ErrUserAlreadyExists && err != nil {
		m.log.Error("Error from users.Service.Register()", zap.Error(err))
	}
	return
}

func (m *loggingMiddleware) Login(ctx context.Context, email string, password string) (u *User, err error) {
	u, err = m.service.Login(ctx, email, password)
	if err != ErrInvalidCredentials && err != nil {
		m.log.Error("Error from users.Service.Login()", zap.Error(err))
	}
	return
}
