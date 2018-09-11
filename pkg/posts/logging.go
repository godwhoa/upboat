package posts

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

type loggingMiddleware struct {
	service Service
	log     *zap.Logger
}

func isImportant(err error) bool {
	switch err {
	case ErrPostNotFound, nil, ErrUnauthorized:
		return false
	default:
		return true
	}
}

func (m *loggingMiddleware) Create(ctx context.Context, post *Post) (id int, err error) {
	id, err = m.service.Create(ctx, post)
	if isImportant(err) {
		m.log.Error("Error from posts.Service.Create()", zap.Error(err))
	}
	return
}

func (m *loggingMiddleware) Get(ctx context.Context, postID int) (post *Post, err error) {
	post, err = m.service.Get(ctx, postID)
	if isImportant(err) {
		m.log.Error("Error from posts.Service.Get()", zap.Error(err))
	}
	return
}

func (m *loggingMiddleware) Edit(ctx context.Context, post *Post) (err error) {
	err = m.service.Edit(ctx, post)
	if isImportant(err) {
		m.log.Error("Error from posts.Service.Edit()", zap.Error(err))
	}
	return
}

func (m *loggingMiddleware) Delete(ctx context.Context, postID, authorID int) (err error) {
	err = m.service.Delete(ctx, postID, authorID)
	if isImportant(err) {
		m.log.Error("Error from posts.Service.Delete()", zap.Error(err))
	}
	return
}

func (m *loggingMiddleware) Vote(ctx context.Context, postID, voterID, delta int) (err error) {
	err = m.service.Vote(ctx, postID, voterID, delta)
	if isImportant(err) {
		m.log.Error("Error from posts.Service.Vote()", zap.Error(err))
	}
	return
}

func (m *loggingMiddleware) Unvote(ctx context.Context, postID, voterID int) (err error) {
	err = m.service.Unvote(ctx, postID, voterID)
	if isImportant(err) {
		m.log.Error("Error from posts.Service.Unvote()", zap.Error(err))
	}
	return
}

func (m *loggingMiddleware) Votes(ctx context.Context, postID int) (votes int, err error) {
	votes, err = m.service.Votes(ctx, postID)
	if isImportant(err) {
		m.log.Error("Error from posts.Service.Votes()", zap.Error(err))
	}
	return
}
