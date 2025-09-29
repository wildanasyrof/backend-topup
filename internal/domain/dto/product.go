package dto

import (
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/pkg/pagination"
)

// --- ProductCreateRequest: Used for creating a new product. ---
type ProductCreateRequest struct {
	Name        string  `form:"name" json:"name" validate:"required,min=3,max=255"`
	SkuCode     string  `form:"sku_code" json:"sku_code" validate:"required,min=1,max=255"`       // Added
	SellerName  string  `form:"seller_name" json:"seller_name" validate:"required,min=3,max=255"` // Added
	CategoryID  uint64  `form:"category_id" json:"category_id" validate:"required,gt=0"`
	ProviderID  uint64  `form:"provider_id" json:"provider_id" validate:"required,gt=0"`
	Status      string  `form:"status" json:"status" validate:"required,oneof=active inactive problem"`
	Stock       int64   `form:"stock" json:"stock" validate:"required,gte=0"`           // Added
	BasePrice   float64 `form:"base_price" json:"base_price" validate:"required,gte=0"` // Added
	Description string  `form:"description" json:"description"`
	ImgURL      string  `form:"img_url" json:"img_url" validate:"omitempty,url"`
	StartOff    string  `form:"start_off" json:"start_off" validate:"omitempty"` // Added (assuming string date/time format)
	EndOff      string  `form:"end_off" json:"end_off" validate:"omitempty"`     // Added (assuming string date/time format)
}

// --- ProductUpdateRequest: Used for partial updates. Uses pointers and omitempty. ---
type ProductUpdateRequest struct {
	Name        *string  `form:"name,omitempty" json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	SkuCode     *string  `form:"sku_code,omitempty" json:"sku_code,omitempty" validate:"omitempty,min=1,max=255"`       // Added
	SellerName  *string  `form:"seller_name,omitempty" json:"seller_name,omitempty" validate:"omitempty,min=3,max=255"` // Added
	CategoryID  *uint64  `form:"category_id,omitempty" json:"category_id,omitempty" validate:"omitempty,gt=0"`
	ProviderID  *uint64  `form:"provider_id,omitempty" json:"provider_id,omitempty" validate:"omitempty,gt=0"`
	Status      *string  `form:"status,omitempty" json:"status,omitempty" validate:"omitempty,oneof=active inactive problem"`
	Stock       *int64   `form:"stock,omitempty" json:"stock,omitempty" validate:"omitempty,gte=0"`           // Added
	BasePrice   *float64 `form:"base_price,omitempty" json:"base_price,omitempty" validate:"omitempty,gte=0"` // Added
	Description *string  `form:"description,omitempty" json:"description,omitempty"`
	ImgURL      *string  `form:"img_url,omitempty" json:"img_url,omitempty" validate:"omitempty,url"`
	StartOff    *string  `form:"start_off,omitempty" json:"start_off,omitempty"` // Added
	EndOff      *string  `form:"end_off,omitempty" json:"end_off,omitempty"`     // Added
}

// --- ProductListQuery: Used for filtering/sorting product lists. ---
type ProductListQuery struct {
	pagination.Query

	ProviderID *uint `query:"provider_id"`
	CategoryID *uint `query:"category_id"`
	LevelID    *uint `query:"level_id"` // if prices vary by user level
	Active     *bool `query:"active"`
}

// ToEntity maps the provided update fields from the DTO to the entity.
func (updateReq *ProductUpdateRequest) ToEntity(
	product *entity.Product,
) {
	if updateReq.Name != nil {
		product.Name = *updateReq.Name
	}
	if updateReq.SkuCode != nil { // Added
		product.SkuCode = *updateReq.SkuCode
	}
	if updateReq.SellerName != nil { // Added
		product.SellerName = *updateReq.SellerName
	}
	if updateReq.CategoryID != nil {
		product.CategoryID = int(*updateReq.CategoryID)
	}
	if updateReq.ProviderID != nil {
		product.ProviderID = int64(*updateReq.ProviderID)
	}
	if updateReq.Status != nil {
		product.Status = entity.CatStatus(*updateReq.Status)
	}
	if updateReq.Stock != nil { // Added
		product.Stock = *updateReq.Stock
	}
	if updateReq.BasePrice != nil { // Added
		product.BasePrice = *updateReq.BasePrice
	}
	if updateReq.Description != nil {
		product.Description = *updateReq.Description
	}
	if updateReq.ImgURL != nil {
		product.ImgUrl = *updateReq.ImgURL
	}
	if updateReq.StartOff != nil { // Added
		product.StartOff = *updateReq.StartOff
	}
	if updateReq.EndOff != nil { // Added
		product.EndOff = *updateReq.EndOff
	}
}
