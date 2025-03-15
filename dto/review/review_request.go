package review_dto

type ReviewRequest struct {
	UserID            string `json:"user_id" binding:"required"`
	InvitationThemeID int    `json:"invitation_theme_id" binding:"required"`
	Star              int    `json:"star" binding:"required"`
	Comment           string `json:"comment" binding:"required"`
}
