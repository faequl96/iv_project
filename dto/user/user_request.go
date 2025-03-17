package user_dto

// CreateUserRequest represents the required data for creating a new user
type CreateUserRequest struct {
	ID       string `json:"id" binding:"required"` // Firebase UID
	Email    string `json:"email" binding:"required,email"`
	UserName string `json:"user_name" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
}

// UpdateUserRequest represents the data structure for updating an existing user
// Fields are optional; only provided fields will be updated
type UpdateUserRequest struct {
	FullName string `json:"full_name"` // Updated full name (optional)
}
