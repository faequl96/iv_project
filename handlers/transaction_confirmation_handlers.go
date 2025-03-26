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
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

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
				ErrorResponse(w, http.StatusNotFound, "No invitation found with the provided ID.")
				return
			}
			invitation.Status = models.InvitationStatusActive
			err = h.InvitationRepositories.UpdateInvitation(invitation)
			if err != nil {
				ErrorResponse(w, http.StatusInternalServerError, "Failed to update invitation.")
				return
			}
		}

		if transaction.ProductType == models.ProductIVCoinPackage {
			ivCoinPackage, err := h.IVCoinPackageRepositories.GetIVCoinPackageByID(uint(transaction.ProductID))
			if err != nil {
				ErrorResponse(w, http.StatusNotFound, "No iv coin package found with the provided ID.")
				return
			}

			ivCoin, err := h.IVCoinRepositories.GetIVCoinByUserID(transaction.UserID)
			if err != nil {
				ErrorResponse(w, http.StatusNotFound, "No iv coin found with the provided user.")
				return
			}

			ivCoin.Balance = ivCoin.Balance + ivCoinPackage.CoinAmount

			err = h.IVCoinRepositories.UpdateIVCoin(ivCoin)
			if err != nil {
				ErrorResponse(w, http.StatusInternalServerError, "Failed to update iv coin.")
				return
			}
		}
	}
	if transactionStatus == "cancel" || transactionStatus == "expire" || transactionStatus == "deny" || transactionStatus == "failure" {
		transaction.Status = models.TransactionStatusCanceled
	}

	err = h.TransactionRepositories.UpdateTransaction(transaction)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update transaction.")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *transactionConfirmationHandlers) ManualByAdminByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	role := r.Context().Value(middleware.RoleKey).(string)
	if role != models.UserRoleSuperAdmin.String() && role != models.UserRoleAdmin.String() {
		ErrorResponse(w, http.StatusForbidden, "You do not have permission to access this resource.")
		return
	}

	var request transaction_confirmation_dto.UpdateTransactionConfirmationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid JSON format.")
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid transaction ID format.")
		return
	}

	transaction, err := h.TransactionRepositories.GetTransactionByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No transaction found with the provided ID.")
		return
	}

	if transaction.ProductType == models.ProductInvitation {
		invitation, err := h.InvitationRepositories.GetInvitationByID(uint(transaction.ProductID))
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "No invitation found with the provided ID.")
			return
		}

		if transaction.PaymentMethod == models.PaymentMethodManualTransfer {
			transaction.Status = request.Status
			if transaction.Status == models.TransactionStatusConfirmed {
				invitation.Status = models.InvitationStatusActive
			}

			err = h.InvitationRepositories.UpdateInvitation(invitation)
			if err != nil {
				ErrorResponse(w, http.StatusInternalServerError, "Failed to update invitation.")
				return
			}
		}
	}

	if transaction.ProductType == models.ProductIVCoinPackage {
		ivCoinPackage, err := h.IVCoinPackageRepositories.GetIVCoinPackageByID(uint(transaction.ProductID))
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "No iv coin package found with the provided ID.")
			return
		}

		ivCoin, err := h.IVCoinRepositories.GetIVCoinByUserID(transaction.UserID)
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "No iv coin found with the provided user.")
			return
		}

		transaction.Status = request.Status
		if transaction.Status == models.TransactionStatusConfirmed {
			ivCoin.Balance = ivCoin.Balance + ivCoinPackage.CoinAmount
		}

		err = h.IVCoinRepositories.UpdateIVCoin(ivCoin)
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to update iv coin.")
			return
		}

	}

	err = h.TransactionRepositories.UpdateTransaction(transaction)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update transaction.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Transaction updated successfully", ConvertToTransactionResponse(transaction))
}
