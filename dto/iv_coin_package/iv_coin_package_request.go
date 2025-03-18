package iv_coin_package_dto

type CreateIVCoinPackageRequest struct {
	Name             string `json:"name" binding:"required"`
	CoinAmount       uint   `json:"coin_amount" binding:"required,gte=0"`
	IDRPrice         uint   `json:"idr_price" binding:"required,gte=0"`
	IDRDiscountPrice uint   `json:"idr_discount_price" binding:"required,gte=0"`
	IVCPrice         uint   `json:"ivc_price" binding:"required,gte=0"`
	IVCDiscountPrice uint   `json:"ivc_discount_price" binding:"required,gte=0"`
}

type UpdateIVCoinPackageRequest struct {
	Name             string `json:"name"`
	CoinAmount       uint   `json:"coin_amount"`
	IDRPrice         uint   `json:"idr_price"`
	IDRDiscountPrice uint   `json:"idr_discount_price"`
	IVCPrice         uint   `json:"ivc_price"`
	IVCDiscountPrice uint   `json:"ivc_discount_price"`
}
