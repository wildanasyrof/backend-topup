package dto

type CreateBannerRequest struct {
	ImgUrl string `form:"img_url" validate:"required"`
}
