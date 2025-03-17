package iv_coin_package_dto

type CreateIVCoinPackageRequest struct {
	Name          string `json:"name" binding:"required"`
	CoinAmount    uint   `json:"coin_amount" binding:"required,gte=0"`
	Price         uint   `json:"price" binding:"required,gte=0"`
	DiscountPrice uint   `json:"discount_price" binding:"required,gte=0"`
}

type UpdateIVCoinPackageRequest struct {
	Name          string `json:"name"`
	CoinAmount    uint   `json:"coin_amount"`
	Price         uint   `json:"price"`
	DiscountPrice uint   `json:"discount_price"`
}
