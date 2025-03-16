package user_dto

type UserRequest struct {
	ID       string        `json:"id" binding:"required"` // Firebase UID
	Email    string        `json:"email" binding:"required,email"`
	UserName string        `json:"user_name" binding:"required"`
	FullName string        `json:"full_name" binding:"required"`
	IVCoin   IVCoinRequest `json:"iv_coin"`
}

type IVCoinRequest struct {
	Balance uint `json:"balance" binding:"required"`
}
