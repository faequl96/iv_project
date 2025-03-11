package handlers

import (
	"encoding/json"
	"iv_project/dto"
	invitation_dto "iv_project/dto/invitation_theme"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type invitationThemeHandlers struct {
	InvitationThemeRepositories repositories.InvitationThemeRepositories
}

func InvitationThemeHandler(InvitationThemeRepositories repositories.InvitationThemeRepositories) *invitationThemeHandlers {
	return &invitationThemeHandlers{InvitationThemeRepositories}
}

func (h *invitationThemeHandlers) CreateInvitationTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(invitation_dto.InvitationThemeRequest)
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

	invitationTheme := models.InvitationTheme{
		Title:       request.Title,
		NormalPrice: request.NormalPrice,
		DiskonPrice: request.DiskonPrice,
		Category:    request.Category,
		Rating:      request.Rating,
	}

	err = h.InvitationThemeRepositories.CreateInvitationTheme(invitationTheme)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: "Create Success"}
	json.NewEncoder(w).Encode(response)
}

func (h *invitationThemeHandlers) GetInvitationThemes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	invitationThemes, err := h.InvitationThemeRepositories.GetInvitationThemes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitationThemes}
	json.NewEncoder(w).Encode(response)
}

func (h *invitationThemeHandlers) GetInvitationThemesByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	category := mux.Vars(r)["category"]

	invitationThemes, err := h.InvitationThemeRepositories.GetInvitationThemesByCategory(category)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitationThemes}
	json.NewEncoder(w).Encode(response)
}

func (h *invitationThemeHandlers) GetInvitationThemeByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	invitationTheme, err := h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitationTheme}
	json.NewEncoder(w).Encode(response)
}

func (h *invitationThemeHandlers) UpdateInvitationTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(invitation_dto.InvitationThemeRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	invitationTheme, err := h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.Title != "" {
		invitationTheme.Title = request.Title
	}
	invitationTheme.NormalPrice = request.NormalPrice
	invitationTheme.DiskonPrice = request.DiskonPrice
	if request.Category != "" {
		invitationTheme.Category = request.Category
	}
	if request.Rating > 1 && request.Rating <= 5 {
		invitationTheme.Rating = request.Rating
	}

	err = h.InvitationThemeRepositories.UpdateInvitationTheme(invitationTheme)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: "Update Success"}
	json.NewEncoder(w).Encode(response)
}

func (h *invitationThemeHandlers) DeleteInvitationTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	invitationTheme, err := h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	err = h.InvitationThemeRepositories.DeleteInvitationTheme(invitationTheme)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: "Delete Success"}
	json.NewEncoder(w).Encode(response)
}
