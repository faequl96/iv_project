package handlers

import (
	"encoding/json"
	"iv_project/dto"
	invitation_dto "iv_project/dto/invitation"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type invitationDataHandlers struct {
	InvitationDataRepositories repositories.InvitationDataRepositories
}

func InvitationDataHandler(InvitationDataRepositories repositories.InvitationDataRepositories) *invitationDataHandlers {
	return &invitationDataHandlers{InvitationDataRepositories}
}

func (h *invitationDataHandlers) CreateInvitationData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode JSON request
	var request invitation_dto.InvitationDataRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// Validasi request
	validation := validator.New()
	if err := validation.Struct(request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// Parse event date
	eventDate, err := time.Parse(time.RFC3339, request.EventDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	uploadedFiles, ok := r.Context().Value(middleware.UploadedFilesKey).(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Tidak ada file yang diunggah"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Buat objek InvitationData
	invitationData := models.InvitationData{
		EventName: request.EventName,
		EventDate: eventDate,
		Location:  request.Location,
		Gallery: &models.Gallery{
			ImageURL1:  uploadedFiles["image_url_1"],
			ImageURL2:  uploadedFiles["image_url_2"],
			ImageURL3:  uploadedFiles["image_url_3"],
			ImageURL4:  uploadedFiles["image_url_4"],
			ImageURL5:  uploadedFiles["image_url_5"],
			ImageURL6:  uploadedFiles["image_url_6"],
			ImageURL7:  uploadedFiles["image_url_7"],
			ImageURL8:  uploadedFiles["image_url_8"],
			ImageURL9:  uploadedFiles["image_url_9"],
			ImageURL10: uploadedFiles["image_url_10"],
			ImageURL11: uploadedFiles["image_url_11"],
			ImageURL12: uploadedFiles["image_url_12"],
		},
	}

	// Simpan ke database
	err = h.InvitationDataRepositories.CreateInvitationData(invitationData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// Beri respons sukses
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitationData}
	json.NewEncoder(w).Encode(response)
}

func (h *invitationDataHandlers) GetInvitationDataByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	invitationData, err := h.InvitationDataRepositories.GetInvitationDataByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitationData}
	json.NewEncoder(w).Encode(response)
}

// Update InvitationData
func (h *invitationDataHandlers) UpdateInvitationData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request invitation_dto.InvitationDataRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	invitationData, err := h.InvitationDataRepositories.GetInvitationDataByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// Parse event date
	eventDate, err := time.Parse(time.RFC3339, request.EventDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	uploadedFiles, ok := r.Context().Value(middleware.UploadedFilesKey).(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Tidak ada file yang diunggah"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.EventName != "" {
		invitationData.EventName = request.EventName
	}
	if request.EventDate != "" {
		invitationData.EventDate = eventDate
	}
	if request.Location != "" {
		invitationData.Location = request.Location
	}
	invitationData.Gallery.ImageURL1 = uploadedFiles["image_url_1"]
	invitationData.Gallery.ImageURL2 = uploadedFiles["image_url_2"]
	invitationData.Gallery.ImageURL3 = uploadedFiles["image_url_3"]
	invitationData.Gallery.ImageURL4 = uploadedFiles["image_url_4"]
	invitationData.Gallery.ImageURL5 = uploadedFiles["image_url_5"]
	invitationData.Gallery.ImageURL6 = uploadedFiles["image_url_6"]
	invitationData.Gallery.ImageURL7 = uploadedFiles["image_url_7"]
	invitationData.Gallery.ImageURL8 = uploadedFiles["image_url_8"]
	invitationData.Gallery.ImageURL9 = uploadedFiles["image_url_9"]
	invitationData.Gallery.ImageURL10 = uploadedFiles["image_url_10"]
	invitationData.Gallery.ImageURL11 = uploadedFiles["image_url_11"]
	invitationData.Gallery.ImageURL12 = uploadedFiles["image_url_12"]

	err = h.InvitationDataRepositories.UpdateInvitationData(invitationData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitationData}
	json.NewEncoder(w).Encode(response)
}

// Delete InvitationData
func (h *invitationDataHandlers) DeleteInvitationData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	invitationData, err := h.InvitationDataRepositories.GetInvitationDataByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = h.InvitationDataRepositories.DeleteInvitationData(invitationData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: "InvitationData deleted successfully"}
	json.NewEncoder(w).Encode(response)
}
