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
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid request format",
			"id": "Format request tidak valid",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	if err := validator.New().Struct(request); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Validation failed. Please complete the request field",
			"id": "Validasi gagal. Silahkan lengkapi field request",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
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
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An Error occurred while creating the review. Please try again later.",
			"id": "Terjadi kesalahan saat membuat review. Coba lagi nanti.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	user, err := h.UserRepositories.GetUserByID(userID)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No user found with the provided ID.",
			"id": "Pengguna tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	review.User = user

	SuccessResponse(w, http.StatusCreated, "Review successfully created", ConvertToReviewResponse(review))
}

func (h *reviewHandlers) GetReviewByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid review ID format. Please provide a numeric ID.",
			"id": "Format ID review tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	review, err := h.ReviewRepositories.GetReviewByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No review found with the provided ID.",
			"id": "Review tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Review found successfully", ConvertToReviewResponse(review))
}

func (h *reviewHandlers) GetReviews(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "You do not have permission to access this resource.",
			"id": "Anda tidak memiliki izin untuk mengakses sumber daya ini.",
		}
		ErrorResponse(w, http.StatusForbidden, messages, lang)
		return
	}

	reviews, err := h.ReviewRepositories.GetReviews()
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while fetching reviews.",
			"id": "Terjadi kesalahan saat mengambil review.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
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
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid invitation theme ID format. Please provide a numeric ID.",
			"id": "Format ID tema undangan tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	reviews, err := h.ReviewRepositories.GetReviewsByInvitationThemeID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while retrieving reviews. Please try again.",
			"id": "Terjadi kesalahan saat mengambil review. Silahkan coba lagi nanti",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
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
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid request format",
			"id": "Format request tidak valid",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	if err := validator.New().Struct(request); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Validation failed. Please complete the request field",
			"id": "Validasi gagal. Silahkan lengkapi field request",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid review ID format. Please provide a numeric ID.",
			"id": "Format ID review tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	review, err := h.ReviewRepositories.GetReviewByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No review found with the provided ID.",
			"id": "Review tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if request.Star != 0 {
		review.Star = request.Star
	}
	if request.Comment != "" {
		review.Comment = request.Comment
	}

	if err := h.ReviewRepositories.UpdateReview(review); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while updating the review.",
			"id": "Terjadi kesalahan saat mengupdate review.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Review successfully updated", ConvertToReviewResponse(review))
}

func (h *reviewHandlers) DeleteReviewByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Invalid review ID format. Please provide a numeric ID.",
			"id": "Format ID review tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	if _, err = h.ReviewRepositories.GetReviewByID(uint(id)); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No review found with the provided ID.",
			"id": "Review tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if err := h.ReviewRepositories.DeleteReview(uint(id)); err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "An error occurred while deleting the review.",
			"id": "Terjadi kesalahan saat menghapus review.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Review successfully deleted", nil)
}
