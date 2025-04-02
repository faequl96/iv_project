package handlers

import (
	"encoding/json"
	transaction_confirmation_dto "iv_project/dto/transaction_confirmation"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type transactionConfirmationHandlers struct {
	TransactionRepositories   repositories.TransactionRepositories
	InvitationRepositories    repositories.InvitationRepositories
	IVCoinPackageRepositories repositories.IVCoinPackageRepositories
	IVCoinRepositories        repositories.IVCoinRepositories
}

func TransactionConfirmationHandler(
	TransactionRepositories repositories.TransactionRepositories,
	InvitationRepositories repositories.InvitationRepositories,
	IVCoinPackageRepositories repositories.IVCoinPackageRepositories,
	IVCoinRepositories repositories.IVCoinRepositories,
) *transactionConfirmationHandlers {
	return &transactionConfirmationHandlers{
		TransactionRepositories,
		InvitationRepositories,
		IVCoinPackageRepositories,
		IVCoinRepositories,
	}
}

func (h *transactionConfirmationHandlers) AutoByMidtrans(w http.ResponseWriter, r *http.Request) {
	var notification map[string]any
	json.NewDecoder(r.Body).Decode(&notification)

	transactionStatus := notification["transaction_status"].(string)
	referenceNumber := notification["order_id"].(string)

	transaction, _ := h.TransactionRepositories.GetTransactionByReferenceNumber(referenceNumber)

	if transactionStatus == "pending" {
		transaction.Status = models.TransactionStatusPending
	}
	if transactionStatus == "settlement" {
		transaction.Status = models.TransactionStatusConfirmed

		if transaction.ProductType == models.ProductInvitation {
			invitation, err := h.InvitationRepositories.GetInvitationByID(uint(transaction.ProductID))
			if err != nil {
				return
			}
			invitation.Status = models.InvitationStatusActive
			err = h.InvitationRepositories.UpdateInvitation(invitation)
			if err != nil {
				return
			}
		}

		if transaction.ProductType == models.ProductIVCoinPackage {
			ivCoinPackage, err := h.IVCoinPackageRepositories.GetIVCoinPackageByID(uint(transaction.ProductID))
			if err != nil {
				return
			}

			ivCoin, err := h.IVCoinRepositories.GetIVCoinByUserID(transaction.UserID)
			if err != nil {
				return
			}

			ivCoin.Balance = ivCoin.Balance + ivCoinPackage.CoinAmount

			err = h.IVCoinRepositories.UpdateIVCoin(ivCoin)
			if err != nil {
				return
			}
		}
	}
	if transactionStatus == "cancel" || transactionStatus == "expire" || transactionStatus == "deny" || transactionStatus == "failure" {
		transaction.Status = models.TransactionStatusCanceled
	}

	if err := h.TransactionRepositories.UpdateTransaction(transaction); err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *transactionConfirmationHandlers) ManualByAdminByID(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "You do not have permission to access this resource.",
			"id": "Anda tidak memiliki izin untuk mengakses sumber daya ini.",
		}
		ErrorResponse(w, http.StatusForbidden, messages, lang)
		return
	}

	var request transaction_confirmation_dto.TransactionConfirmationRequest
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
			"en": "Invalid transaction ID format. Please provide a numeric ID.",
			"id": "Format ID transaksi tidak valid. Harap berikan ID dalam format angka.",
		}
		ErrorResponse(w, http.StatusBadRequest, messages, lang)
		return
	}

	transaction, err := h.TransactionRepositories.GetTransactionByID(uint(id))
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "No transaction found with the provided ID.",
			"id": "Transaksi tidak ditemukan dengan ID yang diberikan.",
		}
		ErrorResponse(w, http.StatusNotFound, messages, lang)
		return
	}

	if transaction.ProductType == models.ProductInvitation {
		invitation, err := h.InvitationRepositories.GetInvitationByID(uint(transaction.ProductID))
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "No invitation found with the provided ID.",
				"id": "Undangan tidak ditemukan dengan ID yang diberikan.",
			}
			ErrorResponse(w, http.StatusNotFound, messages, lang)
			return
		}

		if transaction.PaymentMethod == models.PaymentMethodManualTransfer {
			transaction.Status = request.Status
			if transaction.Status == models.TransactionStatusConfirmed {
				invitation.Status = models.InvitationStatusActive
			}

			err = h.InvitationRepositories.UpdateInvitation(invitation)
			if err != nil {
				lang, _ := r.Context().Value(middleware.LanguageKey).(string)
				messages := map[string]string{
					"en": "Failed to update invitation",
					"id": "Gagal mengupdate undangan",
				}
				ErrorResponse(w, http.StatusInternalServerError, messages, lang)
				return
			}
		}
	}

	if transaction.ProductType == models.ProductIVCoinPackage {
		ivCoinPackage, err := h.IVCoinPackageRepositories.GetIVCoinPackageByID(uint(transaction.ProductID))
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "No iv coin package found with the provided ID.",
				"id": "Paket iv coin tidak ditemukan dengan ID yang diberikan.",
			}
			ErrorResponse(w, http.StatusNotFound, messages, lang)
			return
		}

		ivCoin, err := h.IVCoinRepositories.GetIVCoinByUserID(transaction.UserID)
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "No iv coin found with the provided user ID.",
				"id": "IV coin tidak ditemukan dengan ID user yang diberikan.",
			}
			ErrorResponse(w, http.StatusNotFound, messages, lang)
			return
		}

		transaction.Status = request.Status
		if transaction.Status == models.TransactionStatusConfirmed {
			ivCoin.Balance = ivCoin.Balance + ivCoinPackage.CoinAmount
		}

		err = h.IVCoinRepositories.UpdateIVCoin(ivCoin)
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "Failed to update iv coin.",
				"id": "Gagal mengupdate iv coin.",
			}
			ErrorResponse(w, http.StatusInternalServerError, messages, lang)
			return
		}

	}

	err = h.TransactionRepositories.UpdateTransaction(transaction)
	if err != nil {
		lang, _ := r.Context().Value(middleware.LanguageKey).(string)
		messages := map[string]string{
			"en": "Failed to update transaction.",
			"id": "Gagal mengupdate transaksi.",
		}
		ErrorResponse(w, http.StatusInternalServerError, messages, lang)
		return
	}

	SuccessResponse(w, http.StatusOK, "Transaction updated successfully", ConvertToTransactionResponse(transaction))
}
