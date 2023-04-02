package base

import (
	"context"
	"email-project/model"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type baseService struct {
	db     mongo.Client
	users  map[int]model.User
	nextID int
}

func NewService(db mongo.Client) Service {
	return &baseService{}
}

type Service interface {
	Authenticate(ctx context.Context, req model.AuthenticateRequest) (res model.AuthenticateResponse, err error)
	CreateUserProfile(ctx context.Context, req model.CreateUserProfileRequest) (model.User, error)
}

func (s *baseService) Authenticate(ctx context.Context, req model.AuthenticateRequest) (res model.AuthenticateResponse, err error) {
	// TODO: Implement authentication logic here
	// For example, check if the username and password match a record in your database
	// Hash the password before comparing it to the stored hash
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

// CreateUserProfile creates a new user profile
func (s *baseService) CreateUserProfile(ctx context.Context, req model.CreateUserProfileRequest) (model.User, error) {
	// Generate new user ID
	s.nextID++

	// Create new user profile
	user := model.User{
		ID:    s.nextID,
		Name:  req.Name,
		Email: req.Email,
	}

	// Add new user to map
	s.users[user.ID] = user

	return user, nil
}
