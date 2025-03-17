package handlers

import (
	"encoding/json"
	iv_coin_dto "iv_project/dto/iv_coin"
	user_dto "iv_project/dto/user"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type userHandlers struct {
	UserRepositories repositories.UserRepositories
}

// UserHandlers initializes the handler with the given repository.
func UserHandlers(UserRepositories repositories.UserRepositories) *userHandlers {
	return &userHandlers{UserRepositories}
}

// Convert model User ke DTO UserResponse
func convertToUserResponse(user *models.User) user_dto.UserResponse {
	return user_dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		UserName: user.UserName,
		FullName: user.FullName,
		IVCoin: &iv_coin_dto.IVCoinResponse{
			Balance: user.IVCoin.Balance,
		},
	}
}

// CreateUser handles user registration
func (h *userHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode request body
	request := new(user_dto.CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	// Validate request fields
	validation := validator.New()
	if err := validation.Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	// Create a new user with default IVCoin balance
	user := &models.User{
		ID:       request.ID,
		Email:    request.Email,
		UserName: request.UserName,
		FullName: request.FullName,
		IVCoin: &models.IVCoin{
			Balance: 0, // Set default balance
		},
	}

	// Store user in the database
	if err := h.UserRepositories.CreateUser(user); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to register user: "+err.Error())
		return
	}

	// Convert model ke DTO dan kirim response
	SuccessResponse(w, http.StatusCreated, "User registered successfully", convertToUserResponse(user))
}

// GetUserByID retrieves a user by their ID
func (h *userHandlers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from request URL
	id := mux.Vars(r)["id"]

	// Fetch the user from the database
	user, err := h.UserRepositories.GetUserByID(id)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "User with ID "+id+" not found")
		return
	}

	// Convert model ke DTO dan kirim response
	SuccessResponse(w, http.StatusOK, "User retrieved successfully", convertToUserResponse(user))
}

// GetUsers retrieves all users
func (h *userHandlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Fetch all users from the database
	users, err := h.UserRepositories.GetUsers()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve users: "+err.Error())
		return
	}

	// Convert slice of users ke slice UserResponse DTO
	var userResponses []user_dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, convertToUserResponse(&user))
	}

	SuccessResponse(w, http.StatusOK, "Users retrieved successfully", userResponses)
}

// UpdateUser modifies an existing user's details
func (h *userHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from request URL
	id := mux.Vars(r)["id"]

	// Fetch the user from the database
	user, err := h.UserRepositories.GetUserByID(id)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "User with ID "+id+" not found")
		return
	}

	// Decode the incoming JSON request body
	request := new(user_dto.UpdateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	// Update only the provided fields
	if request.FullName != "" {
		user.FullName = request.FullName
	}

	// Save the updated user data
	if err = h.UserRepositories.UpdateUser(user); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update user: "+err.Error())
		return
	}

	// Convert model ke DTO dan kirim response
	SuccessResponse(w, http.StatusOK, "User updated successfully", convertToUserResponse(user))
}

// DeleteUser removes a user from the database
func (h *userHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from request URL
	id := mux.Vars(r)["id"]

	// Check if the user exists before attempting to delete
	if _, err := h.UserRepositories.GetUserByID(id); err != nil {
		ErrorResponse(w, http.StatusNotFound, "User with ID "+id+" not found")
		return
	}

	// Delete the user from the database
	if err := h.UserRepositories.DeleteUser(id); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to delete user: "+err.Error())
		return
	}

	SuccessResponse(w, http.StatusOK, "User deleted successfully", nil)
}
