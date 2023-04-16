package base

import (
	"context"
	"email-project/model"
	"time"

	"github.com/go-kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) Authenticate(ctx context.Context, req model.AuthenticateRequest) (res model.AuthenticateResponse, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "authenticate", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Authenticate(ctx, req)
}

func (mw loggingMiddleware) CreateUserProfile(ctx context.Context, req model.CreateUserProfileRequest) (res model.UserPage, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "create user profile", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.CreateUserProfile(ctx, req)
}

func (mw loggingMiddleware) CreateUser(ctx context.Context, req model.CreateUserRequest) (res model.CreateUserResponse, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "create user", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.CreateUser(ctx, req)
}
