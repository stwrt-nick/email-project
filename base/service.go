package base

import (
	"context"
	"email-project/api"
	"email-project/model"
	"errors"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type baseService struct {
	db     mongo.Database
	dbInfo model.DBInfo
	logger *logrus.Logger
	users  map[int]model.UserPage
	nextID int
}

func NewService(db mongo.Database, dbInfo model.DBInfo) Service {
	return &baseService{
		db:     db,
		dbInfo: dbInfo,
	}
}

type Service interface {
	Authenticate(ctx context.Context, req model.AuthenticateRequest) (res model.AuthenticateResponse, err error)
	CreateUserProfile(ctx context.Context, req model.CreateUserProfileRequest) (model.UserPage, error)
	CreateUser(ctx context.Context, req model.CreateUserRequest) (res model.CreateUserResponse, err error)
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
func (s *baseService) CreateUserProfile(ctx context.Context, req model.CreateUserProfileRequest) (model.UserPage, error) {
	// Generate new user ID
	s.nextID++

	// Create new user profile
	user := model.UserPage{
		ID:    s.nextID,
		Name:  req.Name,
		Email: req.Email,
	}

	// Add new user to map
	s.users[user.ID] = user

	return user, nil
}

func (s *baseService) CreateUser(ctx context.Context, req model.CreateUserRequest) (res model.CreateUserResponse, err error) {
	// check if user exists

	count, err := s.db.Collection(s.dbInfo.DBCollectionName).CountDocuments(ctx, bson.M{"username": req.Username})
	if err != nil {
		s.logger.Error(err)
		return res, err
	}
	if count > 0 {
		return res, ErrUserExists
	}

	hashPass, err := api.HashPassword(req.Password)
	if err != nil {
		return res, err
	}
	// create user
	user := model.User{
		Username: req.Username,
		Password: hashPass,
	}

	_, err = s.db.Collection(s.dbInfo.DBCollectionName).InsertOne(ctx, user)
	if err != nil {
		return res, err
	}

	return res, err
}
