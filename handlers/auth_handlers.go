package handlers

import (
	"encoding/json"
	auth_dto "iv_project/dto/auth"
	iv_coin_dto "iv_project/dto/iv_coin"
	user_dto "iv_project/dto/user"
	user_profile_dto "iv_project/dto/user_profile"
	"iv_project/models"
	jwtToken "iv_project/pkg/jwt"
	"iv_project/repositories"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
)

type authHandlers struct {
	UserRepositories repositories.UserRepositories
}

func AuthHandlers(UserRepositories repositories.UserRepositories) *authHandlers {
	return &authHandlers{UserRepositories}
}

func ConvertToAuthResponse(token string, user *models.User) auth_dto.AuthResponse {
	userResponse := user_dto.UserResponse{
		ID:   user.ID,
		Role: user.Role,
	}
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

	return auth_dto.AuthResponse{
		Token: token,
		User:  userResponse,
	}
}

func (h *userHandlers) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request user_dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format: "+err.Error())
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	var user *models.User
	var message string

	userFinded, err := h.UserRepositories.GetUserByID(request.ID)
	if err != nil {
		user = &models.User{
			ID: request.ID,
			IVCoin: &models.IVCoin{
				Balance: 0,
				UserID:  request.ID,
			},
		}
		message = "User created successfully"
	}
	if userFinded != nil {
		user = userFinded
		message = "User already created"
	}

	if user.Email == "faequl96@gmail.com" {
		user.Role = models.UserRoleSuperAdmin
	} else {
		user.Role = models.UserRoleUser
	}

	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix() // 2 day expired

	token, _ := jwtToken.GenerateToken(&claims)

	if userFinded != nil {
		if err := h.UserRepositories.CreateUser(user); err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
			return
		}
	}

	SuccessResponse(w, http.StatusCreated, message, ConvertToAuthResponse(token, user))
}
