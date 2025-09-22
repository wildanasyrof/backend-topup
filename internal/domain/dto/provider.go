package dto

type ProviderRequest struct {
	Name string `json:"name" validate:"required,max=255"`
	Slug string `json:"slug" validate:"required,max=255"`
}

type ProviderUpdate struct {
	Name string `json:"name,omitempty" validate:"max=255"`
	Slug string `json:"slug,omitempty" validate:"max=255"`
}
