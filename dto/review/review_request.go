package review_dto

type CreateReviewRequest struct {
	Star    int    `json:"star" validate:"required"`
	Comment string `json:"comment" validate:"required"`

	InvitationThemeID uint   `json:"invitation_theme_id" validate:"required"`
	UserID            string `json:"user_id" validate:"required"`
}

type UpdateReviewRequest struct {
	Star    int    `json:"star"`
	Comment string `json:"comment"`
}
