package dto

import "github.com/wildanasyrof/backend-topup/internal/domain/entity"

type CreatePrice struct {
	ProductID   int     `json:"product_id" validate:"required,gte=1"`
	UserLevelID int     `json:"user_level_id" validate:"required,gte=1"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

type UpdatePrice struct {
	ProductID   *int     `json:"product_id" validate:"omitempty,gte=1"`
	UserLevelID *int     `json:"user_level_id" validate:"omitempty,gte=1"`
	Price       *float64 `json:"price" validate:"omitempty,gt=0"`
}

func (req *UpdatePrice) ToEntity(price *entity.Price) {

	if req.UserLevelID != nil {
		price.UserLevelID = *req.UserLevelID
	}

	if req.ProductID != nil {
		price.ProductID = *req.ProductID
	}

	if req.Price != nil {
		price.Price = *req.Price
	}

}
