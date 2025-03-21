package iv_coin_package_dto

import discount_category_dto "iv_project/dto/discount_category"

type IVCoinPackageResponse struct {
	ID                 uint                                             `json:"id"`
	Name               string                                           `json:"name"`
	CoinAmount         uint                                             `json:"coin_amount"`
	IDRPrice           uint                                             `json:"idr_price"`
	IDRDiscountPrice   uint                                             `json:"idr_discount_price"`
	DiscountCategories []discount_category_dto.DiscountCategoryResponse `json:"discount_categories"`
}
