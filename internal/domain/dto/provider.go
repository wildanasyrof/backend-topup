package dto

type ProviderRequest struct {
	Name string `json:"name" validate:"required,max=255"`
	Ref  string `json:"ref" validate:"required,max=255"`
}

type ProviderUpdate struct {
	Name string `json:"name,omitempty" validate:"max=255"`
	Ref  string `json:"ref,omitempty" validate:"max=255"`
}
