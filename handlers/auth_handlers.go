package handlers

import (
	"encoding/json"
	auth_dto "iv_project/dto/auth"
	"iv_project/models"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/pkg/middleware"
	"iv_project/pkg/utils"
	"iv_project/repositories"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type authHandlers struct {
	JwtServices      jwtToken.JWTServices
	UserRepositories repositories.UserRepositories
}

func AuthHandlers(JwtServices jwtToken.JWTServices, UserRepositories repositories.UserRepositories) *authHandlers {
	return &authHandlers{JwtServices, UserRepositories}
}

func ConvertToAuthResponse(token string, user *models.User) auth_dto.AuthResponse {
	authResponse := auth_dto.AuthResponse{Token: token}

	if user != nil {
		userResponse := ConvertToUserResponse(user)
		authResponse.User = userResponse
	}

	return authResponse
}

func (h *authHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var request auth_dto.AuthRequest
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

	var user *models.User

	userFinded, err := h.UserRepositories.GetUserByID(request.ID)
	if err != nil {
		firstName, lastName := utils.SplitName(request.DisplayName)
		user = &models.User{
			ID:          request.ID,
			UnixID:      utils.GenerateUnixID(),
			UserProfile: &models.UserProfile{UserID: request.ID, Email: request.Email, FirstName: firstName, LastName: lastName},
			IVCoin:      &models.IVCoin{Balance: 0, UserID: request.ID, AdMobMarker: 0, AdMobLastUpdateAt: time.Now()},
		}
	}
	if userFinded != nil {
		user = userFinded
	}

	if user.Role == "" {
		user.Role = models.UserRoleUser
		if user.UserProfile.Email == "faequl96@gmail.com" {
			user.Role = models.UserRoleSuperAdmin
		}
	}

	if userFinded == nil {
		if err := h.UserRepositories.CreateUser(user); err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "Failed to create user",
				"id": "Gagal membuat pengguna",
			}
			ErrorResponse(w, http.StatusInternalServerError, messages, lang)
			return
		}
	}

	token, _ := h.JwtServices.GenerateToken(user.ID, user.Role)

	SuccessResponse(w, http.StatusOK, "User login successfully", ConvertToAuthResponse(token, user))
}
