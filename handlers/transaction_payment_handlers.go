package handlers

import (
	"fmt"
	"iv_project/models"
	"iv_project/pkg/middleware"
	"iv_project/repositories"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type transactionPaymentHandlers struct {
	TransactionRepositories   repositories.TransactionRepositories
	InvitationRepositories    repositories.InvitationRepositories
	IVCoinPackageRepositories repositories.IVCoinPackageRepositories
	IVCoinRepositories        repositories.IVCoinRepositories
	UserRepositories          repositories.UserRepositories
}

func TransactionPaymentHandler(
	TransactionRepositories repositories.TransactionRepositories,
	InvitationRepositories repositories.InvitationRepositories,
	IVCoinPackageRepositories repositories.IVCoinPackageRepositories,
	IVCoinRepositories repositories.IVCoinRepositories,
	UserRepositories repositories.UserRepositories,
) *transactionPaymentHandlers {
	return &transactionPaymentHandlers{
		TransactionRepositories,
		InvitationRepositories,
		IVCoinPackageRepositories,
		IVCoinRepositories,
		UserRepositories,
	}
}

func (h *transactionPaymentHandlers) IssueByID(w http.ResponseWriter, r *http.Request) {
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

	if transaction.PaymentMethod == models.PaymentMethodIVCoin {
		ivCoin, err := h.IVCoinRepositories.GetIVCoinByUserID(transaction.UserID)
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "No iv coin found with the provided user.")
			return
		}

		if ivCoin.Balance < transaction.IVCTotalPrice {
			ErrorResponse(w, http.StatusPaymentRequired, fmt.Sprintf("Insufficient IVCoin balance: %d/%d IVC.", ivCoin.Balance, transaction.IVCTotalPrice))
			return
		}
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

		ivCoin.Balance = ivCoin.Balance - transaction.IVCTotalPrice

		err = h.IVCoinRepositories.UpdateIVCoin(ivCoin)
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to update iv coin.")
			return
		}

		err = h.TransactionRepositories.UpdateTransaction(transaction)
		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Failed to update transaction.")
			return
		}

		SuccessResponse(w, http.StatusOK, "Transaction issue successfully", ConvertToTransactionResponse(transaction))
	}

	if transaction.PaymentMethod == models.PaymentMethodManualTransfer {
		uploadedFile, ok := r.Context().Value(middleware.UploadsKey).(string)
		if ok {
			transaction.PaymentProofImageUrl = uploadedFile
		}
		SuccessResponse(w, http.StatusOK, "Upload payment proof successfully", ConvertToTransactionResponse(transaction))
	}

	if transaction.PaymentMethod == models.PaymentMethodAutoTransfer || transaction.PaymentMethod == models.PaymentMethodGopay {
		user, err := h.UserRepositories.GetUserByID(transaction.UserID)
		if err != nil {
			ErrorResponse(w, http.StatusNotFound, "User with ID "+transaction.UserID+" not found")
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

		if transaction.PaymentMethod == models.PaymentMethodAutoTransfer {
			snapReq.EnabledPayments = append(snapReq.EnabledPayments, snap.PaymentTypeBankTransfer)
		}
		if transaction.PaymentMethod == models.PaymentMethodGopay {
			snapReq.EnabledPayments = append(snapReq.EnabledPayments, snap.PaymentTypeGopay)
			snapReq.Gopay = &snap.GopayDetails{
				EnableCallback: true,
				CallbackUrl:    "yourapp://payment-callback",
			}
		}

		snapResp, _ := s.CreateTransaction(snapReq)

		SuccessResponse(w, http.StatusCreated, "Connection to payment gateway successfully", map[string]any{
			"midtrans_redirect_url": snapResp.RedirectURL,
			"transaction":           ConvertToTransactionResponse(transaction),
		})
	}

}
