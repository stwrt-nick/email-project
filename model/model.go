package model

// Authenticate Request
type AuthenticateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Authenticate Response
type AuthenticateResponse struct {
	Token string `json:"token"`
	Err   error  `json:"error,omitempty"`
}

// User Profile Request
type CreateUserProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// User represents a user page
type UserPage struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// User represents an account
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Create User Request
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Create User Response
type CreateUserResponse struct {
	Token string `json:"token"`
	Err   string `json:"err,omitempty"`
}

// DB Info
type DBInfo struct {
	DBName           string `json:"dbName"`
	DBCollectionName string `json:"dbCollectionName"`
}
