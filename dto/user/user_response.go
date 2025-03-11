package user_dto

type UserResponse struct {
	ID       uint   `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}
