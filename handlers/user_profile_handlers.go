package handlers

import (
	"encoding/json"
	user_profile_dto "iv_project/dto/user_profile"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type userProfileHandlers struct {
	UserProfileRepositories repositories.UserProfileRepositories
}

func UserProfileHandlers(UserProfileRepositories repositories.UserProfileRepositories) *userProfileHandlers {
	return &userProfileHandlers{UserProfileRepositories}
}

func ConvertToUserProfileResponse(userProfile *models.UserProfile) user_profile_dto.UserProfileResponse {
	return user_profile_dto.UserProfileResponse{
		ID:        userProfile.ID,
		FirstName: userProfile.FirstName,
		LastName:  userProfile.LastName,
	}
}

func (h *userProfileHandlers) GetUserProfileByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	userProfile, err := h.UserProfileRepositories.GetUserProfileByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "UserProfile not found")
		return
	}

	SuccessResponse(w, http.StatusOK, "UserProfile retrieved successfully", ConvertToUserProfileResponse(userProfile))
}

func (h *userProfileHandlers) GetUserProfileByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["userId"]
	userProfile, err := h.UserProfileRepositories.GetUserProfileByUserID(id)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No user profile found with the provided user.")
		return
	}

	SuccessResponse(w, http.StatusOK, "IV coin retrieved successfully", ConvertToUserProfileResponse(userProfile))
}

func (h *userProfileHandlers) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	userProfile, err := h.UserProfileRepositories.GetUserProfileByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "UserProfile not found")
		return
	}

	request := new(user_profile_dto.UpdateUserProfileRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	if request.FirstName != "" {
		userProfile.FirstName = request.FirstName
	}
	if request.LastName != "" {
		userProfile.LastName = request.LastName
	}

	if err = h.UserProfileRepositories.UpdateUserProfile(userProfile); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update UserProfile: "+err.Error())
		return
	}

	SuccessResponse(w, http.StatusOK, "UserProfile updated successfully", ConvertToUserProfileResponse(userProfile))
}
