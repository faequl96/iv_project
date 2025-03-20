package handlers

import (
	"encoding/json"
	iv_coin_dto "iv_project/dto/iv_coin"
	user_dto "iv_project/dto/user"
	user_profile_dto "iv_project/dto/user_profile"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type userHandlers struct {
	UserRepositories repositories.UserRepositories
}

func UserHandlers(UserRepositories repositories.UserRepositories) *userHandlers {
	return &userHandlers{UserRepositories}
}

func ConvertToUserResponse(user *models.User) user_dto.UserResponse {
	userResponse := user_dto.UserResponse{ID: user.ID, Role: user.Role}
	if user.UserProfile != nil {
		userResponse.UserProfile = &user_profile_dto.UserProfileResponse{
			ID:        user.UserProfile.ID,
			FirstName: user.UserProfile.FirstName,
			LastName:  user.UserProfile.LastName,
		}
	}
	if user.IVCoin != nil {
		userResponse.IVCoin = &iv_coin_dto.IVCoinResponse{
			ID:      user.IVCoin.ID,
			Balance: user.IVCoin.Balance,
		}
	}

	return userResponse
}

func (h *userHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(middleware.UserIdKey).(string)
	user, err := h.UserRepositories.GetUserByID(userID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "User with ID "+userID+" not found")
		return
	}

	SuccessResponse(w, http.StatusOK, "User retrieved successfully", ConvertToUserResponse(user))
}

func (h *userHandlers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	id := mux.Vars(r)["id"]
	user, err := h.UserRepositories.GetUserByID(id)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "User with ID "+id+" not found")
		return
	}

	SuccessResponse(w, http.StatusOK, "User retrieved successfully", ConvertToUserResponse(user))
}

func (h *userHandlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	users, err := h.UserRepositories.GetUsers()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve users: "+err.Error())
		return
	}

	if len(users) == 0 {
		SuccessResponse(w, http.StatusOK, "No users available at this moment", []user_dto.UserResponse{})
		return
	}

	var userResponses []user_dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, ConvertToUserResponse(&user))
	}

	SuccessResponse(w, http.StatusOK, "Users retrieved successfully", userResponses)
}

func (h *userHandlers) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	var request user_dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	id := mux.Vars(r)["id"]
	user, err := h.UserRepositories.GetUserByID(id)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	user.Role = models.StringToUserRole(request.Role)
	if user.Role.String() == "" {
		user.Role = models.UserRoleUser
	}

	if err = h.UserRepositories.UpdateUser(user); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update User: "+err.Error())
		return
	}

	SuccessResponse(w, http.StatusOK, "User updated successfully", ConvertToUserResponse(user))
}

func (h *userHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := r.Context().Value(middleware.UserIdKey).(string)
	if _, err := h.UserRepositories.GetUserByID(userID); err != nil {
		ErrorResponse(w, http.StatusNotFound, "User with ID "+userID+" not found")
		return
	}

	if err := h.UserRepositories.DeleteUser(userID); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to delete user: "+err.Error())
		return
	}

	SuccessResponse(w, http.StatusOK, "User deleted successfully", nil)
}

func (h *userHandlers) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	id := mux.Vars(r)["id"]
	if _, err := h.UserRepositories.GetUserByID(id); err != nil {
		ErrorResponse(w, http.StatusNotFound, "User with ID "+id+" not found")
		return
	}

	if err := h.UserRepositories.DeleteUser(id); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to delete user: "+err.Error())
		return
	}

	SuccessResponse(w, http.StatusOK, "User deleted successfully", nil)
}
