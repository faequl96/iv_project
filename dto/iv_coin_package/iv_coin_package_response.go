package iv_coin_package_dto

type IVCoinPackageResponse struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	CoinAmount       uint   `json:"coin_amount"`
	IDRPrice         uint   `json:"idr_price"`
	IDRDiscountPrice uint   `json:"idr_discount_price"`
	IVCPrice         uint   `json:"ivc_price"`
	IVCDiscountPrice uint   `json:"ivc_discount_price"`
}
