package handlers

import (
	"encoding/json"
	iv_coin_dto "iv_project/dto/iv_coin"
	user_dto "iv_project/dto/user"
	user_profile_dto "iv_project/dto/user_profile"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"

	"github.com/gorilla/mux"
)

type userHandlers struct {
	UserRepositories repositories.UserRepositories
}

func UserHandlers(UserRepositories repositories.UserRepositories) *userHandlers {
	return &userHandlers{UserRepositories}
}

func convertToUserResponse(user *models.User) user_dto.UserResponse {
	return user_dto.UserResponse{
		ID: user.ID,
		UserProfile: &user_profile_dto.UserProfileResponse{
			ID:        user.UserProfile.ID,
			FirstName: user.UserProfile.FirstName,
			LastName:  user.UserProfile.LastName,
		},
		IVCoin: &iv_coin_dto.IVCoinResponse{
			ID:      user.IVCoin.ID,
			Balance: user.IVCoin.Balance,
		},
	}
}

func (h *userHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(user_dto.CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	user := &models.User{
		ID: request.ID,
		UserProfile: &models.UserProfile{
			FirstName: request.UserProfile.FirstName,
			LastName:  request.UserProfile.LastName,
		},
		IVCoin: &models.IVCoin{
			Balance: 0,
		},
	}

	if err := h.UserRepositories.CreateUser(user); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to register user: "+err.Error())
		return
	}

	SuccessResponse(w, http.StatusCreated, "User registered successfully", convertToUserResponse(user))
}

func (h *userHandlers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	user, err := h.UserRepositories.GetUserByID(id)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "User with ID "+id+" not found")
		return
	}

	SuccessResponse(w, http.StatusOK, "User retrieved successfully", convertToUserResponse(user))
}

func (h *userHandlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.UserRepositories.GetUsers()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve users: "+err.Error())
		return
	}

	var userResponses []user_dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, convertToUserResponse(&user))
	}

	SuccessResponse(w, http.StatusOK, "Users retrieved successfully", userResponses)
}

func (h *userHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
