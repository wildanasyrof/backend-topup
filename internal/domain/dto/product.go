package dto

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
)

type ProductCreateRequest struct {
	Name        string `form:"name" validate:"required,min=3,max=255"`
	CategoryID  uint64 `form:"category_id" validate:"required,gt=0"`
	ProviderID  uint64 `form:"provider_id" validate:"required,gt=0"`
	Status      string `form:"status" validate:"required,oneof=active inactive problem"`
	Description string `form:"description"`
	ImgURL      string `form:"img_url"`
}

type ProductUpdateRequest struct {
	Name        *string `form:"name,omitempty" validate:"omitempty,min=3,max=255"`
	CategoryID  *uint64 `form:"category_id,omitempty" validate:"omitempty,gt=0"`
	ProviderID  *uint64 `form:"provider_id,omitempty" validate:"omitempty,gt=0"`
	Status      *string `form:"status,omitempty" validate:"omitempty,oneof=active inactive problem"`
	Description *string `form:"description,omitempty"`
	ImgURL      *string `form:"img_url,omitempty" validate:"omitempty,url"`
}

func (updateReq *ProductUpdateRequest) ToEntity(
	product *entity.Product,
) {
	// Check if the Name field is provided in the request.
	if updateReq.Name != nil {
		product.Name = *updateReq.Name
	}

	// Check if the CategoryID field is provided.
	if updateReq.CategoryID != nil {
		product.CategoryID = int(*updateReq.CategoryID)
	}

	// Check if the ProviderID field is provided.
	if updateReq.ProviderID != nil {
		product.ProviderID = int64(*updateReq.ProviderID)
	}

	// Check if the Status field is provided.
	if updateReq.Status != nil {
		product.Status = entity.CatStatus(*updateReq.Status)
	}

	// Check if the Description field is provided.
	if updateReq.Description != nil {
		product.Description = *updateReq.Description
	}

	// Check if the ImgURL field is provided.
	if updateReq.ImgURL != nil {
		product.ImgUrl = *updateReq.ImgURL
	}
}
