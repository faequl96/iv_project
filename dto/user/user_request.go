package user_dto

type CreateUserRequest struct {
	ID       string `json:"id" binding:"required"`
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
}
