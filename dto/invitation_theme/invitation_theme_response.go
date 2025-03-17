package invitation_theme_dto

import "iv_project/models"

type InvitationThemeResponse struct {
	ID          uint             `json:"id"`
	Title       string           `json:"title"`
	NormalPrice int              `json:"normal_price"`
	DiskonPrice int              `json:"diskon_price"`
	Category    string           `json:"category"`
	Reviews     []*models.Review `json:"reviews,omitempty"`
}
