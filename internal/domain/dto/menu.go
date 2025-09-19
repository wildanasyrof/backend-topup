package dto

type CreateMenuRequest struct {
	Name string `json:"name" validate:"required,max=255"`
}
