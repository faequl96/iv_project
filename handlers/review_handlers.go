package handlers

import (
	"encoding/json"
	review_dto "iv_project/dto/review"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type reviewHandlers struct {
	ReviewRepositories repositories.ReviewRepositories
	UserRepositories   repositories.UserRepositories
}

func ReviewHandlers(
	ReviewRepositories repositories.ReviewRepositories,
	UserRepositories repositories.UserRepositories,
) *reviewHandlers {
	return &reviewHandlers{ReviewRepositories, UserRepositories}
}

func ConvertToReviewResponse(review *models.Review) review_dto.ReviewResponse {
	reviewResponse := review_dto.ReviewResponse{
		ID:        review.ID,
		Star:      review.Star,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt.Format(time.RFC3339),
		UpdatedAt: review.UpdatedAt.Format(time.RFC3339),
	}

	if review.User != nil {
		userResponse := ConvertToUserResponse(review.User)
		reviewResponse.User = &userResponse
	}

	return reviewResponse
}

func (h *reviewHandlers) CreateReview(w http.ResponseWriter, r *http.Request) {
	var request review_dto.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format. Please check your input.")
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	userID := r.Context().Value(middleware.UserIdKey).(string)
	review := &models.Review{
		Star:              request.Star,
		Comment:           request.Comment,
		UserID:            userID,
		InvitationThemeID: request.InvitationThemeID,
	}

	if err := h.ReviewRepositories.CreateReview(review); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while create the review. Please try again.")
		return
	}

	user, err := h.UserRepositories.GetUserByID(userID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "User with ID "+userID+" not found")
		return
	}

	review.User = user

	SuccessResponse(w, http.StatusCreated, "Review successfully created", ConvertToReviewResponse(review))
}

func (h *reviewHandlers) GetReviewByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format. Please use a number.")
		return
	}

	review, err := h.ReviewRepositories.GetReviewByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Review not found with the given ID.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Review found successfully", ConvertToReviewResponse(review))
}

func (h *reviewHandlers) GetReviews(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	reviews, err := h.ReviewRepositories.GetReviews()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching iv coin packages.")
		return
	}

	if len(reviews) == 0 {
		SuccessResponse(w, http.StatusOK, "No reviews available at this moment", []review_dto.ReviewResponse{})
		return
	}

	var responses []review_dto.ReviewResponse
	for _, review := range reviews {
		responses = append(responses, ConvertToReviewResponse(&review))
	}

	SuccessResponse(w, http.StatusOK, "Reviews retrieved successfully", responses)
}

func (h *reviewHandlers) GetReviewsByInvitationThemeID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["invitationThemeId"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format. Please use a number.")
		return
	}

	reviews, err := h.ReviewRepositories.GetReviewsByInvitationThemeID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while retrieving reviews. Please try again.")
		return
	}

	if len(reviews) == 0 {
		SuccessResponse(w, http.StatusOK, "No reviews available for this invitation theme", []review_dto.ReviewResponse{})
		return
	}

	var response []review_dto.ReviewResponse
	for _, review := range reviews {
		response = append(response, ConvertToReviewResponse(&review))
	}

	SuccessResponse(w, http.StatusOK, "Reviews retrieved successfully", response)
}

func (h *reviewHandlers) UpdateReviewByID(w http.ResponseWriter, r *http.Request) {
	var request review_dto.UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format. Please check your input.")
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format. Please use a number.")
		return
	}

	review, err := h.ReviewRepositories.GetReviewByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "Review not found with the given ID.")
		return
	}

	if request.Star != 0 {
		review.Star = request.Star
	}
	if request.Comment != "" {
		review.Comment = request.Comment
	}

	if err := h.ReviewRepositories.UpdateReview(review); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while updating the review. Please try again.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Review successfully updated", ConvertToReviewResponse(review))
}

func (h *reviewHandlers) DeleteReviewByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid ID format. Please use a number.")
		return
	}

	if _, err = h.ReviewRepositories.GetReviewByID(uint(id)); err != nil {
		ErrorResponse(w, http.StatusNotFound, "No review found with the provided ID.")
		return
	}

	if err := h.ReviewRepositories.DeleteReview(uint(id)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while deleting the review. Please try again.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Review successfully deleted", nil)
}
