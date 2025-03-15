package handlers

import (
	"encoding/json"
	"fmt"
	"iv_project/dto"
	user_dto "iv_project/dto/user"
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

func (h *userHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("tes1")

	request := new(user_dto.CreateUserRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Println("tes2")

	// validation := validator.New()
	// err := validation.Struct(request)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
	// 	json.NewEncoder(w).Encode(response)
	// 	return
	// }

	fmt.Println("tes3")

	user := models.User{
		ID:       request.ID,
		UserName: request.UserName,
		Email:    request.Email,
		FullName: request.FullName,
		IVCoin: models.IVCoin{
			Balance: 0,
		},
	}

	fmt.Println(user)

	err := h.UserRepositories.CreateUser(user)
	fmt.Println("tes5")
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
