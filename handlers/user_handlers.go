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

func UserHandlers(UserRepositories repositories.UserRepositories) *userHandlers {
	return &userHandlers{UserRepositories}
}

func (h *userHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(user_dto.CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	user := models.User{
		ID:       request.ID,
		UserName: request.UserName,
		Email:    request.Email,
		FullName: request.FullName,
	}

	err = h.UserRepositories.CreateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: "Success"}
	json.NewEncoder(w).Encode(response)
}

func (h *userHandlers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	user, err := h.UserRepositories.GetUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	UserResponse := user_dto.UserResponse{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		FullName: user.FullName,
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: UserResponse}
	json.NewEncoder(w).Encode(response)
}
