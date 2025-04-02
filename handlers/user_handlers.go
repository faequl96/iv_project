package handlers

import (
	"encoding/json"
	user_dto "iv_project/dto/user"
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
		userProfileResponse := ConvertToUserProfileResponse(user.UserProfile)
		userResponse.UserProfile = &userProfileResponse
	}

	if user.IVCoin != nil {
		ivCoinResponse := ConvertToIVCoinResponse(user.IVCoin)
		userResponse.IVCoin = &ivCoinResponse
	}

	return userResponse
}

func (h *userHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIdKey).(string)
	user, err := h.UserRepositories.GetUserByID(userID)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No user found with the provided ID.",
			"id": "Pengguna tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "User retrieved successfully", ConvertToUserResponse(user))
}

func (h *userHandlers) GetUserByID(w http.ResponseWriter, r *http.Request) {
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

	id := mux.Vars(r)["id"]
	user, err := h.UserRepositories.GetUserByID(id)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No user found with the provided ID.",
			"id": "Pengguna tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "User retrieved successfully", ConvertToUserResponse(user))
}

func (h *userHandlers) GetUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := h.UserRepositories.GetUsers()
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching users.",
			"id": "Terjadi kesalahan saat mengambil user.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
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
	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "You do not have permission to access this resource.",
			"id": "Anda tidak memiliki izin untuk mengakses sumber daya ini.",
		}
		ErrorResponse(w, http.StatusForbidden, messages, lang)
		return
	}

	var request user_dto.UserRequest
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

	id := mux.Vars(r)["id"]
	user, err := h.UserRepositories.GetUserByID(id)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No user found with the provided ID.",
			"id": "Pengguna tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	user.Role = models.StringToUserRoleType(request.Role)
	if user.Role.String() == "" {
		user.Role = models.UserRoleUser
	}

	if err = h.UserRepositories.UpdateUser(user); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Failed to update user.",
			"id": "Gagal mengupdate pengguna.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "User updated successfully", ConvertToUserResponse(user))
}

func (h *userHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIdKey).(string)
	if _, err := h.UserRepositories.GetUserByID(userID); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No user found with the provided ID.",
			"id": "Pengguna tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if err := h.UserRepositories.DeleteUser(userID); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while deleting the user.",
			"id": "Terjadi kesalahan saat menghapus pengguna.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "User deleted successfully", nil)
}

func (h *userHandlers) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "You do not have permission to access this resource.",
			"id": "Anda tidak memiliki izin untuk mengakses sumber daya ini.",
		}
		ErrorResponse(w, http.StatusForbidden, messages, lang)
		return
	}

	id := mux.Vars(r)["id"]
	if _, err := h.UserRepositories.GetUserByID(id); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No user found with the provided ID.",
			"id": "Pengguna tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if err := h.UserRepositories.DeleteUser(id); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while deleting the user.",
			"id": "Terjadi kesalahan saat menghapus pengguna.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "User deleted successfully", nil)
}
