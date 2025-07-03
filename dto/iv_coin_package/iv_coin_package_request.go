package iv_coin_package_dto

type CreateIVCoinPackageRequest struct {
	Name                string `json:"name" validate:"required"`
	CoinAmount          uint   `json:"coin_amount" validate:"required"`
	IDRPrice            uint   `json:"idr_price" validate:"required"`
	DiscountCategoryIDs []uint `json:"discount_category_ids"`
}

type UpdateIVCoinPackageRequest struct {
	Name                string `json:"name"`
	CoinAmount          uint   `json:"coin_amount"`
	IDRPrice            uint   `json:"idr_price"`
	DiscountCategoryIDs []uint `json:"discount_category_ids"`
}
