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

	"github.com/gorilla/mux"
)

type invitationHandlers struct {
	InvitationRepositories repositories.InvitationRepositories
}

func InvitationHandler(InvitationRepositories repositories.InvitationRepositories) *invitationHandlers {
	return &invitationHandlers{InvitationRepositories}
}

func (h *invitationHandlers) CreateInvitation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(invitation_dto.InvitationRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	eventDate, err := time.Parse(time.RFC3339, request.InvitationData.EventDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Invalid event_date format. Use RFC3339 (e.g. 2025-06-15T14:00:00Z)"}
		json.NewEncoder(w).Encode(response)
		return
	}

	uploadedFiles, ok := r.Context().Value(middleware.UploadedFilesKey).(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: "Tidak ada file yang diunggah"}
		json.NewEncoder(w).Encode(response)
		return
	}

	invitation := models.Invitation{
		UserID: request.UserID,
		Status: request.Status,
		InvitationData: models.InvitationData{
			EventName:         request.InvitationData.EventName,
			EventDate:         eventDate,
			Location:          request.InvitationData.Location,
			GalleryImageURL1:  uploadedFiles["gallery_image_url_1"],
			GalleryImageURL2:  uploadedFiles["gallery_image_url_2"],
			GalleryImageURL3:  uploadedFiles["gallery_image_url_3"],
			GalleryImageURL4:  uploadedFiles["gallery_image_url_4"],
			GalleryImageURL5:  uploadedFiles["gallery_image_url_5"],
			GalleryImageURL6:  uploadedFiles["gallery_image_url_6"],
			GalleryImageURL7:  uploadedFiles["gallery_image_url_7"],
			GalleryImageURL8:  uploadedFiles["gallery_image_url_8"],
			GalleryImageURL9:  uploadedFiles["gallery_image_url_9"],
			GalleryImageURL10: uploadedFiles["gallery_image_url_10"],
			GalleryImageURL11: uploadedFiles["gallery_image_url_11"],
			GalleryImageURL12: uploadedFiles["gallery_image_url_12"],
		},
	}

	err = h.InvitationRepositories.CreateInvitation(invitation)
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

func (h *invitationHandlers) GetInvitationByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	invitation, err := h.InvitationRepositories.GetInvitationByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitation}
	json.NewEncoder(w).Encode(response)
}

func (h *invitationHandlers) GetInvitationsByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	invitations, err := h.InvitationRepositories.GetInvitationsByUserID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: invitations}
	json.NewEncoder(w).Encode(response)
}
