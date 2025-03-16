package handlers

import (
	"encoding/json"
	"iv_project/dto"
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

// Constructor function to create a new instance of userHandlers
func UserHandlers(UserRepositories repositories.UserRepositories) *userHandlers {
	return &userHandlers{UserRepositories}
}

// successResponse sends a standardized success response with a status code, message, and data
func successResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.SuccessResult{Code: statusCode, Message: message, Data: data})
}

// errorResponse sends a standardized error response with a status code and message
func errorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.ErrorResult{Code: statusCode, Message: message})
}

// CreateUser handles user registration
func (h *userHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode request body
	request := new(user_dto.UserRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	// Validate request fields
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return
	}

	// Create a new user with default IVCoin balance
	user := models.User{
		ID:       request.ID,
		Email:    request.Email,
		UserName: request.UserName,
		FullName: request.FullName,
		IVCoin: &models.IVCoin{
			Balance: 0, // Set default balance
		},
	}

	// Store user in the database
	err = h.UserRepositories.CreateUser(user)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to register user: "+err.Error())
		return
	}

	successResponse(w, http.StatusCreated, "User registered successfully", user)
}

// GetUserByID retrieves a user by their ID
func (h *userHandlers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from request URL
	id := mux.Vars(r)["id"]

	// Fetch the user from the database
	user, err := h.UserRepositories.GetUserByID(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "User with ID "+id+" not found")
		return
	}

	successResponse(w, http.StatusOK, "User retrieved successfully", user)
}

// GetUsers retrieves all users
func (h *userHandlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Fetch all users from the database
	users, err := h.UserRepositories.GetUsers()
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to retrieve users: "+err.Error())
		return
	}

	successResponse(w, http.StatusOK, "Users retrieved successfully", users)
}

// UpdateUser modifies an existing user's details
func (h *userHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode request body
	request := new(user_dto.UserRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	// Get user ID from request URL
	id := mux.Vars(r)["id"]

	// Fetch the user from the database
	user, err := h.UserRepositories.GetUserByID(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "User with ID "+id+" not found")
		return
	}

	// Update only the provided fields
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.UserName != "" {
		user.UserName = request.UserName
	}
	if request.FullName != "" {
		user.FullName = request.FullName
	}

	// Save the updated user data
	err = h.UserRepositories.UpdateUser(user)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to update user: "+err.Error())
		return
	}

	successResponse(w, http.StatusOK, "User updated successfully", user)
}

// DeleteUser removes a user from the database
func (h *userHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from request URL
	id := mux.Vars(r)["id"]

	// Fetch the user to ensure they exist
	user, err := h.UserRepositories.GetUserByID(id)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "User with ID "+id+" not found")
		return
	}

	// Delete the user from the database
	err = h.UserRepositories.DeleteUser(user)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to delete user: "+err.Error())
		return
	}

	successResponse(w, http.StatusOK, "User deleted successfully", nil)
}
