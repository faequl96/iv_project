package handlers

import (
	"encoding/json"
	"fmt"
	transaction_dto "iv_project/dto/transaction"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type transactionHandlers struct {
	TransactionRepositories   repositories.TransactionRepositories
	InvitationRepositories    repositories.InvitationRepositories
	IVCoinPackageRepositories repositories.IVCoinPackageRepositories
	IVCoinRepositories        repositories.IVCoinRepositories
	UserRepositories          repositories.UserRepositories
}

func TransactionHandler(
	TransactionRepositories repositories.TransactionRepositories,
	InvitationRepositories repositories.InvitationRepositories,
	IVCoinPackageRepositories repositories.IVCoinPackageRepositories,
	IVCoinRepositories repositories.IVCoinRepositories,
	UserRepositories repositories.UserRepositories,
) *transactionHandlers {
	return &transactionHandlers{
		TransactionRepositories,
		InvitationRepositories,
		IVCoinPackageRepositories,
		IVCoinRepositories,
		UserRepositories,
	}
}

func ConvertToTransactionResponse(transaction *models.Transaction) transaction_dto.TransactionResponse {
	transactionResponse := transaction_dto.TransactionResponse{
		ID:              transaction.ID,
		ProductType:     transaction.ProductType,
		ProductName:     transaction.ProductName,
		Status:          transaction.Status,
		PaymentMethod:   transaction.PaymentMethod,
		ReferenceNumber: transaction.ReferenceNumber,
		IDRPrice:        transaction.IDRPrice,
		IDRDiscount:     transaction.IDRDiscount,
		IDRTotalPrice:   transaction.IDRTotalPrice,
		IVCPrice:        transaction.IVCPrice,
		IVCDiscount:     transaction.IVCDiscount,
		IVCTotalPrice:   transaction.IVCTotalPrice,
		CreatedAt:       transaction.CreatedAt.Format(time.RFC3339),
	}

	return transactionResponse
}

func (h *transactionHandlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request transaction_dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request format. Please check your input.")
		return
	}

	if err := validator.New().Struct(request); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	transaction := &models.Transaction{
		ProductType:     request.ProductType,
		Status:          models.TransactionStatusPending,
		PaymentMethod:   request.PaymentMethod,
		ReferenceNumber: GenerateReferenceNumber(request.PaymentMethod.String()),
		ProductID:       request.ProductID,
		UserID:          request.UserID,
	}

	if request.ProductType == models.ProductInvitation {
		invitation, err := h.InvitationRepositories.GetInvitationByID(uint(request.ProductID))
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "No invitation found with the provided ID.")
			return
		}

		transaction.IDRPrice = invitation.InvitationTheme.IDRPrice
		transaction.IDRDiscount = invitation.InvitationTheme.IDRPrice - invitation.InvitationTheme.IDRDiscountPrice
		transaction.IDRTotalPrice = invitation.InvitationTheme.IDRDiscountPrice
		transaction.IVCPrice = invitation.InvitationTheme.IVCPrice
		transaction.IVCDiscount = invitation.InvitationTheme.IVCPrice - invitation.InvitationTheme.IVCDiscountPrice
		transaction.IVCTotalPrice = invitation.InvitationTheme.IVCDiscountPrice

		transaction.ProductName = invitation.InvitationTheme.Name

		if request.PaymentMethod == models.PaymentMethodIVCoin {
			ivCoin, err := h.IVCoinRepositories.GetIVCoinByUserID(request.UserID)
			if err != nil {
				ErrorResponse(w, http.StatusNotFound, "No iv coin found with the provided user.")
				return
			}

			if ivCoin.Balance < transaction.IVCTotalPrice {
				ErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Insufficient IVCoin balance: %d/%d IVC.", ivCoin.Balance, transaction.IVCTotalPrice))
				return
			}

			transaction.Status = models.TransactionStatusConfirmed
			invitation.Status = models.InvitationStatusActive

			err = h.InvitationRepositories.UpdateInvitation(invitation)
			if err != nil {
				ErrorResponse(w, http.StatusInternalServerError, "Failed to update invitation.")
				return
			}

			ivCoin.Balance = ivCoin.Balance - transaction.IVCTotalPrice

			err = h.IVCoinRepositories.UpdateIVCoin(ivCoin)
			if err != nil {
				ErrorResponse(w, http.StatusInternalServerError, "Failed to update iv coin.")
				return
			}
		}

	}

	if request.ProductType == models.ProductIVCoinPackage {
		ivCoinPackage, err := h.IVCoinPackageRepositories.GetIVCoinPackageByID(uint(request.ProductID))
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "No iv coin package found with the provided ID.")
			return
		}

		transaction.IDRPrice = ivCoinPackage.IDRPrice
		transaction.IDRDiscount = ivCoinPackage.IDRPrice - ivCoinPackage.IDRDiscountPrice
		transaction.IDRTotalPrice = ivCoinPackage.IDRDiscountPrice

		transaction.ProductName = ivCoinPackage.Name
	}

	err := h.TransactionRepositories.CreateTransaction(transaction)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to create transaction.")
		return
	}

	if request.PaymentMethod == models.PaymentMethodAutoTransfer || request.PaymentMethod == models.PaymentMethodGopay {
		user, err := h.UserRepositories.GetUserByID(request.UserID)
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "User with ID "+request.UserID+" not found")
			return
		}

		var s = snap.Client{}
		s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

		snapReq := &snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  transaction.ReferenceNumber,
				GrossAmt: int64(transaction.IDRTotalPrice),
			},
			CustomerDetail: &midtrans.CustomerDetails{
				FName: user.UserProfile.FirstName + " " + user.UserProfile.LastName,
				Email: user.Email,
			},
			EnabledPayments: []snap.SnapPaymentType{},
		}

		if request.PaymentMethod == models.PaymentMethodAutoTransfer {
			snapReq.EnabledPayments = append(snapReq.EnabledPayments, snap.PaymentTypeBankTransfer)
		}
		if request.PaymentMethod == models.PaymentMethodGopay {
			snapReq.EnabledPayments = append(snapReq.EnabledPayments, snap.PaymentTypeGopay)
		}

		snapResp, _ := s.CreateTransaction(snapReq)

		SuccessResponse(w, http.StatusCreated, "Transaction created successfully", map[string]any{
			"transaction":           ConvertToTransactionResponse(transaction),
			"midtrans_token":        snapResp.Token,
			"midtrans_redirect_url": snapResp.RedirectURL,
		})
		return
	}

	SuccessResponse(w, http.StatusCreated, "Transaction created successfully", ConvertToTransactionResponse(transaction))
}

func (h *transactionHandlers) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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

	SuccessResponse(w, http.StatusOK, "Transaction retrieved successfully", ConvertToTransactionResponse(transaction))
}

func (h *transactionHandlers) GetTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	transactions, err := h.TransactionRepositories.GetTransactions()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching transactions.")
		return
	}

	var transactionResponses []transaction_dto.TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, ConvertToTransactionResponse(&transaction))
	}

	SuccessResponse(w, http.StatusOK, "Transactions retrieved successfully", transactionResponses)
}

func (h *transactionHandlers) GetTransactionsByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := mux.Vars(r)["userId"]

	transactions, err := h.TransactionRepositories.GetTransactionsByUserID(userID)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while fetching transactions.")
		return
	}

	var transactionResponses []transaction_dto.TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, ConvertToTransactionResponse(&transaction))
	}

	SuccessResponse(w, http.StatusOK, "Transactions retrieved successfully", transactionResponses)
}

func (h *transactionHandlers) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid transaction ID format.")
		return
	}

	if _, err = h.TransactionRepositories.GetTransactionByID(uint(id)); err != nil {
		ErrorResponse(w, http.StatusNotFound, "No transaction found with the provided ID.")
		return
	}

	if err := h.TransactionRepositories.DeleteTransaction(uint(id)); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "An error occurred while deleting the transaction.")
		return
	}

	SuccessResponse(w, http.StatusOK, "Transaction deleted successfully", nil)
}
