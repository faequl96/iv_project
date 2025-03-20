package ad_mob_dto

type AdMobRequest struct {
	Amount uint `json:"amount" validate:"required"`
}
