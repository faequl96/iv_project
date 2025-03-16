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

	// Decode JSON request
	var request user_dto.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	// Validasi request
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	// Buat user dengan ID dari Firebase
	user := models.User{
		ID:       request.ID, // Gunakan Firebase UID dari request
		UserName: request.UserName,
		Email:    request.Email,
		FullName: request.FullName,
		IVCoin: &models.IVCoin{
			Balance: 0,
		},
	}

	err = h.UserRepositories.CreateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.SuccessResult{Code: http.StatusOK, Data: "User registered successfully"})
}

func (h *userHandlers) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	user, err := h.UserRepositories.GetUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{Code: http.StatusNotFound, Message: "User not found"}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: user}
	json.NewEncoder(w).Encode(response)
}

func (h *userHandlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.UserRepositories.GetUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: users}
	json.NewEncoder(w).Encode(response)
}

func (h *userHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	request := new(user_dto.UserRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := h.UserRepositories.GetUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{Code: http.StatusNotFound, Message: "User not found"}
		json.NewEncoder(w).Encode(response)
		return
	}
	if request.UserName != "" {
		user.UserName = request.UserName
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.FullName != "" {
		user.FullName = request.FullName
	}

	err = h.UserRepositories.UpdateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: "User updated successfully"}
	json.NewEncoder(w).Encode(response)
}

func (h *userHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	user, err := h.UserRepositories.GetUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := dto.ErrorResult{Code: http.StatusNotFound, Message: "User not found"}
		json.NewEncoder(w).Encode(response)
		return
	}

	err = h.UserRepositories.DeleteUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: "User deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
