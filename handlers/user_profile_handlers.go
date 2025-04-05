package handlers

import (
	"encoding/json"
	user_profile_dto "iv_project/dto/user_profile"
	"iv_project/models"
	"iv_project/pkg/middleware"
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
	userProfileResponse := user_profile_dto.UserProfileResponse{
		ID:        userProfile.ID,
		Email:     userProfile.Email,
		FirstName: userProfile.FirstName,
		LastName:  userProfile.LastName,
	}

	return userProfileResponse
}

func (h *userProfileHandlers) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIdKey).(string)
	userProfile, err := h.UserProfileRepositories.GetUserProfileByUserID(userID)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No user profile found with the provided user ID.",
			"id": "Profile pengguna tidak ditemukan dengan ID pengguna yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "User profile retrieved successfully", ConvertToUserProfileResponse(userProfile))
}

func (h *userProfileHandlers) GetUserProfileByID(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "You do not have permission to access this resource.",
			"id": "Anda tidak memiliki izin untuk mengakses sumber daya ini.",
		}
		ErrorResponse(w, http.StatusForbidden, messages, lang)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid user profile ID format. Please provide a numeric ID.",
			"id": "Format ID profil pengguna tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	userProfile, err := h.UserProfileRepositories.GetUserProfileByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No user profile found with the provided ID.",
			"id": "Profile pengguna tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "User profile retrieved successfully", ConvertToUserProfileResponse(userProfile))
}

func (h *userProfileHandlers) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var request user_profile_dto.UserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid request format",
			"id": "Format request tidak valid",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	if err := validator.New().Struct(request); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Validation failed. Please complete the request field",
			"id": "Validasi gagal. Silahkan lengkapi field request",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	userID := r.Context().Value(middleware.UserIdKey).(string)
	userProfile, err := h.UserProfileRepositories.GetUserProfileByUserID(userID)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No user profile found with the provided user ID.",
			"id": "Profile pengguna tidak ditemukan dengan ID pengguna yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if request.FirstName != "" {
		userProfile.FirstName = request.FirstName
	}
	if request.LastName != "" {
		userProfile.LastName = request.LastName
	}

	if err = h.UserProfileRepositories.UpdateUserProfile(userProfile); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Failed to update user profile.",
			"id": "Gagal mengupdate profil pengguna.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "User profile updated successfully", ConvertToUserProfileResponse(userProfile))
}

func (h *userProfileHandlers) UpdateUserProfileByID(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "You do not have permission to access this resource.",
			"id": "Anda tidak memiliki izin untuk mengakses sumber daya ini.",
		}
		ErrorResponse(w, http.StatusForbidden, messages, lang)
		return
	}

	var request user_profile_dto.UserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid request format",
			"id": "Format request tidak valid",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	if err := validator.New().Struct(request); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Validation failed. Please complete the request field",
			"id": "Validasi gagal. Silahkan lengkapi field request",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid user profile ID format. Please provide a numeric ID.",
			"id": "Format ID profil pengguna tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	userProfile, err := h.UserProfileRepositories.GetUserProfileByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No user profile found with the provided ID.",
			"id": "Profile pengguna tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if request.FirstName != "" {
		userProfile.FirstName = request.FirstName
	}
	if request.LastName != "" {
		userProfile.LastName = request.LastName
	}

	if err = h.UserProfileRepositories.UpdateUserProfile(userProfile); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Failed to update user profile.",
			"id": "Gagal mengupdate profil pengguna.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "User profile updated successfully", ConvertToUserProfileResponse(userProfile))
}
