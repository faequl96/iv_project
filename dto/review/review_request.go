package review_dto

type CreateReviewRequest struct {
	UserID            string `json:"user_id" binding:"required"`
	InvitationThemeID uint   `json:"invitation_theme_id" binding:"required"`
	Star              int    `json:"star" binding:"required"`
	Comment           string `json:"comment" binding:"required"`
}

type UpdateReviewRequest struct {
	Star    int    `json:"star"`
	Comment string `json:"comment"`
}
