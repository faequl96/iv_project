package iv_coin_dto

// IVCoinRequest represents the data structure for updating or setting the user's IVCoin balance
type IVCoinRequest struct {
	Balance uint `json:"balance" binding:"required"` // User's IVCoin balance, required field and must be a non-negative value
}
