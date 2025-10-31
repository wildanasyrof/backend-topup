package dto

import "github.com/wildanasyrof/backend-topup/pkg/pagination"

type CreateMenuRequest struct {
	Name string `json:"name" validate:"required,max=255"`
}

type MenuListQuery struct {
	pagination.Query // Embeds: Page, Limit, Sort, Q
}
