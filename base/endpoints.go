package base

import (
	"context"
	"email-project/model"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	AuthenticateEndpoint      endpoint.Endpoint
	CreateUserProfileEndpoint endpoint.Endpoint
	CreateUserEndpoint        endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		AuthenticateEndpoint:      makeAuthenticateEndpoint(s),
		CreateUserProfileEndpoint: makeCreateUserProfileEndpoint(s),
		CreateUserEndpoint:        makeCreateUserEndpoint(s),
	}
}

func makeAuthenticateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.AuthenticateRequest)
		res, err := s.Authenticate(ctx, req)
		if err != nil {
			return model.AuthenticateResponse{"", err}, nil
		}
		return model.AuthenticateResponse{res.Token, nil}, nil
	}
}

func makeCreateUserProfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.CreateUserProfileRequest)
		res, err := s.CreateUserProfile(ctx, req)
		if err != nil {
			return model.UserPage{0, "", ""}, nil
		}
		return model.UserPage{res.ID, res.Name, res.Email}, nil
	}
}

func makeCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(model.CreateUserRequest)
		token, err := s.CreateUser(ctx, req)
		if err != nil {
			return model.CreateUserResponse{Err: err.Error()}, nil
		}
		return model.CreateUserResponse{Token: token.Token, Err: token.Err}, nil
	}
}
