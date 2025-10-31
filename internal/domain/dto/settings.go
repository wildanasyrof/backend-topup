package dto

import "github.com/wildanasyrof/backend-topup/pkg/pagination"

type CreateSettingsRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Value string `json:"value" validate:"required"`
}

type UpdateSettingsRequest struct {
	Name  string `json:"name" validate:"omitempty,min=3,max=100"`
	Value string `json:"value" validate:"omitempty"`
}

type SettingsListQuery struct {
	pagination.Query // Embeds: Page, Limit, Sort, Q

	// Filter spesifik
	Name *string `query:"name"`
}
