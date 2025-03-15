package handlers

import (
	"encoding/json"
	"iv_project/dto"
	review_dto "iv_project/dto/review"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type reviewHandlers struct {
	ReviewRepositories repositories.ReviewRepositories
}

func ReviewHandlers(ReviewRepositories repositories.ReviewRepositories) *reviewHandlers {
	return &reviewHandlers{ReviewRepositories}
}

func (h *reviewHandlers) CreateReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(review_dto.ReviewRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	review := models.Review{
		UserID:            request.UserID,
		InvitationThemeID: uint(request.InvitationThemeID),
		Star:              request.Star,
		Comment:           request.Comment,
	}

	err := h.ReviewRepositories.CreateReview(review)
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

func (h *reviewHandlers) GetReviewsByInvitationThemeID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["invitation_theme_id"])

	reviews, err := h.ReviewRepositories.GetReviewsByInvitationThemeID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: reviews}
	json.NewEncoder(w).Encode(response)
}

func (h *reviewHandlers) DeleteReview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	review, err := h.ReviewRepositories.GetReviewByID(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	err = h.ReviewRepositories.DeleteReview(review)
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
