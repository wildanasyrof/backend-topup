package dto

type UpdateUserRequest struct {
	Name     *string `json:"name" validate:"omitempty,min=2,max=100"`
	Whatsapp *string `json:"whatsapp" validate:"omitempty,min=10,max=15"`
}
