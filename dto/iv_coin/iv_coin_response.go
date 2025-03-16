package iv_coin_dto

// IVCoinResponse represents the user's IVCoin balance details
type IVCoinResponse struct {
	Balance uint `json:"balance"` // User's IVCoin balance, must be a non-negative value
}
