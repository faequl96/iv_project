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

	if transaction.PaymentMethod == models.PaymentMethodIVCoin {
		ivCoin, err := h.IVCoinRepositories.GetIVCoinByUserID(transaction.UserID)
		if err != nil {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "No iv coin found with the provided user.",
				"id": "IV coin tidak ditemukan dengan pengguna yang diberikan.",
			}
			ErrorResponse(w, http.StatusNotFound, messages, lang)
			return
		}

		if ivCoin.Balance < transaction.IVCTotalPrice {
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": fmt.Sprintf("Insufficient IVCoin balance: %d/%d IVC.", ivCoin.Balance, transaction.IVCTotalPrice),
				"id": fmt.Sprintf("IV koin tidak cukup: %d/%d IVC.", ivCoin.Balance, transaction.IVCTotalPrice),
			}
			ErrorResponse(w, http.StatusPaymentRequired, messages, lang)
			return
		}
		transaction.Status = models.TransactionStatusConfirmed

		if transaction.ProductType == models.ProductInvitation {
			invitation, err := h.InvitationRepositories.GetInvitationByID(uint(transaction.ProductID))
			if err != nil {
				lang, _ := r.Context().Value(middleware.LanguageKey).(string)
				messages := map[string]string{
					"en": "No transaction found with the provided ID.",
					"id": "Transaksi tidak ditemukan dengan ID yang diberikan.",
				}
				ErrorResponse(w, http.StatusNotFound, messages, lang)
				return
			}
			invitation.Status = models.InvitationStatusActive
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

		ivCoin.Balance = ivCoin.Balance - transaction.IVCTotalPrice

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
			lang, _ := r.Context().Value(middleware.LanguageKey).(string)
			messages := map[string]string{
				"en": "No user found with the provided ID.",
				"id": "Pengguna tidak ditemukan dengan ID yang diberikan.",
			}
			ErrorResponse(w, http.StatusNotFound, messages, lang)
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
				Email: user.UserProfile.Email,
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
