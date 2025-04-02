package handlers

import (
	"encoding/json"
	auth_dto "iv_project/dto/auth"
	"iv_project/models"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/repositories"
	"net/http"

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
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	var user *models.User

	userFinded, err := h.UserRepositories.GetUserByID(request.ID)
	if err != nil {
		user = &models.User{
			ID:          request.ID,
			Email:       request.Email,
			UserProfile: &models.UserProfile{UserID: request.ID},
			IVCoin:      &models.IVCoin{Balance: 0, UserID: request.ID},
		}
	}
	if userFinded != nil {
		user = userFinded
	}

	if user.Role == "" {
		user.Role = models.UserRoleUser
		if user.Email == "faequl96@gmail.com" {
			user.Role = models.UserRoleSuperAdmin
		}
	}

	if userFinded == nil {
		if err := h.UserRepositories.CreateUser(user); err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
			return
		}
	}

	token, _ := h.JwtServices.GenerateToken(user.ID, user.Role)

	SuccessResponse(w, http.StatusOK, "User login successfully", ConvertToAuthResponse(token, user))
}
