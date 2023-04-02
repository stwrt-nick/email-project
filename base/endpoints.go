package base

import (
	"context"
	"email-project/model"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	AuthenticateEndpoint      endpoint.Endpoint
	CreateUserProfileEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		AuthenticateEndpoint:      makeAuthenticateEndpoint(s),
		CreateUserProfileEndpoint: makeCreateUserProfileEndpoint(s),
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
			return model.User{0, "", ""}, nil
		}
		return model.User{res.ID, res.Name, res.Email}, nil
	}
}
