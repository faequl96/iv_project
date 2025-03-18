package review_dto

import user_dto "iv_project/dto/user"

type ReviewResponse struct {
	ID      uint                   `json:"id"`
	Star    int                    `json:"star"`
	Comment string                 `json:"comment"`
	User    *user_dto.UserResponse `json:"user"`
}
