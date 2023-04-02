package base

import (
	"context"
	"email-project/model"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type baseService struct {
	db mongo.Client
}

func NewService(db mongo.Client) Service {
	return &baseService{}
}

type Service interface {
	Authenticate(ctx context.Context, req model.AuthenticateRequest) (res model.AuthenticateResponse, err error)
}

func (s *baseService) Authenticate(ctx context.Context, req model.AuthenticateRequest) (res model.AuthenticateResponse, err error) {
	// TODO: Implement authentication logic here
	// For example, check if the username and password match a record in your database
	// Hash the password before comparing it to the stored hash
	fmt.Println("hello")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return res, err
	}
	if req.Username == "testuser" && bcrypt.CompareHashAndPassword(hashedPassword, []byte("testpassword")) == nil {
		res.Token = "testtoken"
		return res, nil
	}
	return res, errors.New("invalid username or password")
}
