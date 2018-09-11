package posts

import (
	"context"
	"time"

	"go.opencensus.io/trace"
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
	ctx, span := trace.StartSpan(ctx, "Create")
	defer span.End()
	id, err = m.service.Create(ctx, post)
	if err != nil {
		m.log.Error("Error from posts.Service.Create()", zap.Error(err))
	}
	return
}

func (m *loggingMiddleware) Get(postID int) (post *Post, err error) {
	t0 := time.Now()
	defer func() {
		m.log.Info("Latency of posts.Service.Get()",
			zap.String("latency", time.Since(t0).String()),
		)
	}()
	post, err = m.service.Get(postID)
	if err != ErrPostNotFound && err != nil {
		m.log.Error("Error from posts.Service.Get()", zap.Error(err))
	}
	return
}

func (m *loggingMiddleware) Edit(post *Post) (err error) {
	t0 := time.Now()
	defer func() {
		m.log.Info("Latency of posts.Service.Edit()",
			zap.String("latency", time.Since(t0).String()),
		)
	}()
	err = m.service.Edit(post)
	if isImportant(err) {
		m.log.Error("Error from posts.Service.Edit()", zap.Error(err))
	}
	return
}

func (m *loggingMiddleware) Delete(postID, authorID int) (err error) {
	t0 := time.Now()
	defer func() {
		m.log.Info("Latency of posts.Service.Delete()",
			zap.String("latency", time.Since(t0).String()),
		)
	}()
	err = m.service.Delete(postID, authorID)
	if isImportant(err) {
		m.log.Error("Error from posts.Service.Delete()", zap.Error(err))
	}
	return
}

func (m *loggingMiddleware) Vote(postID, voterID, delta int) (err error) {
	t0 := time.Now()
	defer func() {
		m.log.Info("Latency of posts.Service.Vote()",
			zap.String("latency", time.Since(t0).String()),
		)
	}()
	err = m.service.Vote(postID, voterID, delta)
	if isImportant(err) {
		m.log.Error("Error from posts.Service.Vote()", zap.Error(err))
	}
	return
}

func (m *loggingMiddleware) Unvote(postID, voterID int) (err error) {
	t0 := time.Now()
	defer func() {
		m.log.Info("Latency of posts.Service.Unvote()",
			zap.String("latency", time.Since(t0).String()),
		)
	}()
	err = m.service.Unvote(postID, voterID)
	if isImportant(err) {
		m.log.Error("Error from posts.Service.Unvote()", zap.Error(err))
	}
	return
}

func (m *loggingMiddleware) Votes(postID int) (votes int, err error) {
	t0 := time.Now()
	defer func() {
		m.log.Info("Latency of posts.Service.Votes()",
			zap.String("latency", time.Since(t0).String()),
		)
	}()
	votes, err = m.service.Votes(postID)
	if isImportant(err) {
		m.log.Error("Error from posts.Service.Votes()", zap.Error(err))
	}
	return
}
