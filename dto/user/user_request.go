package user_dto

type CreateUserRequest struct {
	AuthID   string `json:"auth_id" binding:"required"`
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required"`
}
