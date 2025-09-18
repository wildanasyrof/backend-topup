package dto

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}
