package iv_coin_package_dto

type IVCoinPackageResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	CoinAmount    uint   `json:"coin_amount"`
	Price         uint   `json:"price"`
	DiscountPrice uint   `json:"discount_price"`
}
