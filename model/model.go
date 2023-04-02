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

// User represents a user profile
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
