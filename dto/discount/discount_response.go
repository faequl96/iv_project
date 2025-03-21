package discount_dto

import (
	invitation_theme_dto "iv_project/dto/invitation_theme"
	iv_coin_package_dto "iv_project/dto/iv_coin_package"
)

type DiscountResponse struct {
	InvitationThemes []invitation_theme_dto.InvitationThemeResponse `json:"invitation_themes"`
	IVCoinPackages   []iv_coin_package_dto.IVCoinPackageResponse    `json:"iv_coin_packages"`
}
