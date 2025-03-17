package invitation_theme_dto

import "iv_project/models"

// InvitationThemeResponse represents the response format for an invitation theme.
type InvitationThemeResponse struct {
	ID          uint             `json:"id"`
	Title       string           `json:"title"`
	NormalPrice int              `json:"normal_price"`
	DiskonPrice int              `json:"diskon_price"`
	Category    string           `json:"category"`
	Reviews     []*models.Review `json:"reviews,omitempty"`
}
