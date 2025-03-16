package handlers

import (
	"encoding/json"
	"iv_project/dto"
	invitation_theme_dto "iv_project/dto/invitation_theme"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"
	"strconv"

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

	request := new(invitation_theme_dto.InvitationThemeRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	invitationTheme := models.InvitationTheme{
		Title:       request.Title,
		NormalPrice: request.NormalPrice,
		DiskonPrice: request.DiskonPrice,
		Category:    request.Category,
	}

	err := h.InvitationThemeRepositories.CreateInvitationTheme(invitationTheme)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Create Success")
}

func (h *invitationThemeHandlers) GetInvitationThemeByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	invitationTheme, err := h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitationTheme}
	json.NewEncoder(w).Encode(response)
}

func (h *invitationThemeHandlers) GetInvitationThemes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	invitationThemes, err := h.InvitationThemeRepositories.GetInvitationThemes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
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
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitationThemes}
	json.NewEncoder(w).Encode(response)
}

func (h *invitationThemeHandlers) UpdateInvitationTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(invitation_theme_dto.InvitationThemeRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	invitationTheme, err := h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
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

	err = h.InvitationThemeRepositories.UpdateInvitationTheme(invitationTheme)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Update Success")
}

func (h *invitationThemeHandlers) DeleteInvitationTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	invitationTheme, err := h.InvitationThemeRepositories.GetInvitationThemeByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = h.InvitationThemeRepositories.DeleteInvitationTheme(invitationTheme)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Delete Success")
}
