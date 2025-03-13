package user_dto

type UserResponse struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}
