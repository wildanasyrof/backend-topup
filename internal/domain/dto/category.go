package dto

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
)

// CreateCategoryRequest: untuk multipart/form-data
type CreateCategoryRequest struct {
	Name        string `form:"name"        validate:"required"`
	Type        string `form:"type"        validate:"required,oneof=prabayar pascabayar"`
	MenuID      int64  `form:"menu_id"     validate:"required"`
	ProviderID  int64  `form:"provider_id" validate:"required"`
	Slug        string `form:"slug"        validate:"required"`
	Status      string `form:"status"      validate:"required,oneof=inactive active problem"`
	Description string `form:"description"`
	InputType   string `form:"input_type"  validate:"required"`
	IsLogin     bool   `form:"is_login"`
	ImgUrl      string `form:"img_url"`
}

// UpdateCategoryRequest: semua opsional; image juga opsional
type UpdateCategoryRequest struct {
	Name        *string `form:"name"`
	Type        *string `form:"type"        validate:"omitempty,oneof=prabayar pascabayar"`
	MenuID      *int64  `form:"menu_id"`
	ProviderID  *int64  `form:"provider_id"`
	Slug        *string `form:"slug"`
	Status      *string `form:"status"      validate:"omitempty,oneof=inactive active problem"`
	Description *string `form:"description"`
	InputType   *string `form:"input_type"`
	IsLogin     *bool   `form:"is_login"`
	ImgUrl      string  `form:"img_url"`
}

func (req *UpdateCategoryRequest) UpdateEntity(category *entity.Category) {
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Type != nil {
		category.Type = entity.Type(*req.Type)
	}
	if req.MenuID != nil {
		category.MenuID = *req.MenuID
	}
	if req.ProviderID != nil {
		category.ProviderID = *req.ProviderID
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}
	if req.Status != nil {
		category.Status = entity.CatStatus(*req.Status)
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.InputType != nil {
		category.InputType = *req.InputType
	}
	if req.IsLogin != nil {
		category.IsLogin = *req.IsLogin
	}
	// Tambahkan kondisi untuk ImgUrl
	if req.ImgUrl != "" {
		category.ImgUrl = req.ImgUrl
	}
}
