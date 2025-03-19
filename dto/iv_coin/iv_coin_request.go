package iv_coin_dto

type IVCoinRequest struct {
	Balance uint `json:"balance" validate:"required"`
}
