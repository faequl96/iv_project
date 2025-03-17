package invitation_theme_dto

import review_dto "iv_project/dto/review"

type InvitationThemeResponse struct {
	ID          uint                         `json:"id"`
	Title       string                       `json:"title"`
	NormalPrice int                          `json:"normal_price"`
	DiskonPrice int                          `json:"diskon_price"`
	Category    string                       `json:"category"`
	Reviews     []*review_dto.ReviewResponse `json:"reviews,omitempty"`
}
