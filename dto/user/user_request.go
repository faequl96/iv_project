package user_dto

// CreateUserRequest represents the required data for creating a new user
type CreateUserRequest struct {
	ID       string `json:"id" binding:"required"`          // Firebase UID, uniquely identifies the user
	Email    string `json:"email" binding:"required,email"` // User's email address, must be a valid email format
	UserName string `json:"user_name" binding:"required"`   // Username, required for user identification
	FullName string `json:"full_name" binding:"required"`   // Full name of the user, required field
}

// UpdateUserRequest represents the data structure for updating an existing user
// Fields are optional; only provided fields will be updated
type UpdateUserRequest struct {
	FullName string `json:"full_name"` // Updated full name (optional)
}
