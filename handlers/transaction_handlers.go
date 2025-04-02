package handlers

import (
	"encoding/json"
	transaction_dto "iv_project/dto/transaction"
	"iv_project/models"
	"iv_project/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type transactionHandlers struct {
	TransactionRepositories   repositories.TransactionRepositories
	InvitationRepositories    repositories.InvitationRepositories
	IVCoinPackageRepositories repositories.IVCoinPackageRepositories
	IVCoinRepositories        repositories.IVCoinRepositories
	UserRepositories          repositories.UserRepositories
	VoucherCodeRepositories   repositories.VoucherCodeRepositories
}

func TransactionHandler(
	TransactionRepositories repositories.TransactionRepositories,
	InvitationRepositories repositories.InvitationRepositories,
	IVCoinPackageRepositories repositories.IVCoinPackageRepositories,
	IVCoinRepositories repositories.IVCoinRepositories,
	UserRepositories repositories.UserRepositories,
	VoucherCodeRepositories repositories.VoucherCodeRepositories,
) *transactionHandlers {
	return &transactionHandlers{
		TransactionRepositories,
		InvitationRepositories,
		IVCoinPackageRepositories,
		IVCoinRepositories,
		UserRepositories,
		VoucherCodeRepositories,
	}
}

func ConvertToTransactionResponse(transaction *models.Transaction) transaction_dto.TransactionResponse {
	transactionResponse := transaction_dto.TransactionResponse{
		ID:                     transaction.ID,
		ProductType:            transaction.ProductType,
		ProductName:            transaction.ProductName,
		Status:                 transaction.Status,
		PaymentMethod:          transaction.PaymentMethod,
		ReferenceNumber:        transaction.ReferenceNumber,
		IDRPrice:               transaction.IDRPrice,
		IDRDiscount:            transaction.IDRDiscount,
		IDRTotalPrice:          transaction.IDRTotalPrice,
		IVCPrice:               transaction.IVCPrice,
		IVCDiscount:            transaction.IVCDiscount,
		IVCTotalPrice:          transaction.IVCTotalPrice,
		VoucherCodeID:          transaction.VoucherCodeID,
		VoucherCodeName:        transaction.VoucherCodeName,
		IDRVoucherCodeDiscount: transaction.IDRVoucherCodeDiscount,
		IVCVoucherCodeDiscount: transaction.IVCVoucherCodeDiscount,
		PaymentProofImageUrl:   transaction.PaymentProofImageUrl,
		CreatedAt:              transaction.CreatedAt.Format(time.RFC3339),
	}

	return transactionResponse
}

func (h *transactionHandlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {
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
		ProductType: request.ProductType,
		Status:      models.TransactionStatusPending,
		ProductID:   request.ProductID,
		UserID:      request.UserID,
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

		ivCoin, _ := h.IVCoinRepositories.GetIVCoinByUserID(transaction.UserID)
		if ivCoin != nil {
			if ivCoin.Balance > transaction.IVCTotalPrice {
				transaction.PaymentMethod = models.PaymentMethodIVCoin
				transaction.ReferenceNumber = GenerateReferenceNumber(transaction.PaymentMethod.String())
			} else {
				transaction.PaymentMethod = models.PaymentMethodGopay
				transaction.ReferenceNumber = GenerateReferenceNumber(transaction.PaymentMethod.String())
			}
		} else {
			transaction.PaymentMethod = models.PaymentMethodGopay
			transaction.ReferenceNumber = GenerateReferenceNumber(transaction.PaymentMethod.String())
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

		transaction.PaymentMethod = models.PaymentMethodGopay
		transaction.ReferenceNumber = GenerateReferenceNumber(transaction.PaymentMethod.String())
	}

	err := h.TransactionRepositories.CreateTransaction(transaction)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to create transaction.")
		return
	}

	SuccessResponse(w, http.StatusCreated, "Transaction created successfully", ConvertToTransactionResponse(transaction))
}

func (h *transactionHandlers) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
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

func (h *transactionHandlers) UpdateTransactionByID(w http.ResponseWriter, r *http.Request) {
	var request transaction_dto.UpdateTransactionRequest
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
		ErrorResponse(w, http.StatusBadRequest, "Invalid transaction ID format.")
		return
	}

	transaction, err := h.TransactionRepositories.GetTransactionByID(uint(id))
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "No transaction found with the provided ID.")
		return
	}

	if transaction.Status == models.TransactionStatusConfirmed || transaction.Status == models.TransactionStatusCanceled {
		ErrorResponse(w, http.StatusNotFound, "The transaction cannot be updated, because the transaction has been completed.")
		return
	}

	if request.PaymentMethod.String() != "" {
		transaction.PaymentMethod = request.PaymentMethod
		transaction.ReferenceNumber = GenerateReferenceNumber(request.PaymentMethod.String())
	}

	if request.VoucherCodeID != 0 {
		voucherCode, err := h.VoucherCodeRepositories.GetVoucherCodeByID(uint(request.VoucherCodeID))
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "No voucher code found with the provided ID.")
			return
		}
		transaction.VoucherCodeID = voucherCode.ID
		transaction.VoucherCodeName = voucherCode.Name

		if transaction.ProductType == models.ProductInvitation {
			totalIDRVoucherCodeDiscount := CalculateDiscountedPrice(transaction.IDRTotalPrice, voucherCode.DiscountPercentage)
			transaction.IDRVoucherCodeDiscount = transaction.IDRTotalPrice - totalIDRVoucherCodeDiscount
			transaction.IDRTotalPrice = totalIDRVoucherCodeDiscount
			totalIVCVoucherCodeDiscount := CalculateDiscountedPrice(transaction.IVCTotalPrice, voucherCode.DiscountPercentage)
			transaction.IVCVoucherCodeDiscount = transaction.IVCTotalPrice - totalIVCVoucherCodeDiscount
			transaction.IVCTotalPrice = totalIVCVoucherCodeDiscount
		}

		if transaction.ProductType == models.ProductIVCoinPackage {
			totalIDRVoucherCodeDiscount := CalculateDiscountedPrice(transaction.IDRTotalPrice, voucherCode.DiscountPercentage)
			transaction.IDRVoucherCodeDiscount = transaction.IDRTotalPrice - totalIDRVoucherCodeDiscount
			transaction.IDRTotalPrice = totalIDRVoucherCodeDiscount
		}
	}

	if request.VoucherCodeID == 0 {
		transaction.VoucherCodeID = 0
		transaction.VoucherCodeName = ""

		if transaction.ProductType == models.ProductInvitation {
			transaction.IDRVoucherCodeDiscount = 0
			transaction.IDRTotalPrice = transaction.IDRPrice - transaction.IDRDiscount
			transaction.IVCVoucherCodeDiscount = 0
			transaction.IVCTotalPrice = transaction.IVCPrice - transaction.IVCDiscount
		}

		if transaction.ProductType == models.ProductIVCoinPackage {
			transaction.IDRVoucherCodeDiscount = 0
			transaction.IDRTotalPrice = transaction.IDRPrice - transaction.IDRDiscount
		}
	}

	err = h.TransactionRepositories.UpdateTransaction(transaction)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Failed to update transaction.")
		return
	}

	SuccessResponse(w, http.StatusCreated, "Transaction created successfully", ConvertToTransactionResponse(transaction))
}

func (h *transactionHandlers) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
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
