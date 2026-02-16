package payment

import (
	"POS-kasir/config"

	"POS-kasir/pkg/logger"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type IMidtrans interface {
	CreateQRISCharge(orderID string, amount int64) (*coreapi.ChargeResponse, error)
	GetQRISCharge(orderID string) (*coreapi.TransactionStatusResponse, error)
	VerifyNotificationSignature(payload MidtransNotificationPayload) error
	CancelTransaction(orderID string) (*coreapi.CancelResponse, error)
}

type MidtransNotificationPayload struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}

type MidtransService struct {
	config *config.AppConfig
	client coreapi.Client
	log    logger.ILogger
}

func NewMidtransService(cfg *config.AppConfig, log logger.ILogger) IMidtrans {
	var client coreapi.Client

	midtransEnv := midtrans.Sandbox
	if cfg.Midtrans.IsProd {
		midtransEnv = midtrans.Production
	} else {
		midtransEnv = midtrans.Sandbox
	}

	client.New(cfg.Midtrans.ServerKey, midtransEnv)

	if cfg.Midtrans.IsProd {
		client.ServerKey = cfg.Midtrans.ServerKey
		client.Options.SetPaymentAppendNotification("https://example.com/notification")
		client.Options.SetPaymentOverrideNotification("https://example.com/notification")
	}

	log.Infof("Midtrans client initialized successfully")

	return &MidtransService{
		client: client,
		log:    log,
		config: cfg,
	}
}

// VerifyNotificationSignature memvalidasi signature key dari payload notifikasi.
func (s *MidtransService) VerifyNotificationSignature(payload MidtransNotificationPayload) error {

	sourceString := payload.OrderID + payload.StatusCode + payload.GrossAmount + s.config.Midtrans.ServerKey

	hasher := sha512.New()
	_, err := hasher.Write([]byte(sourceString))
	if err != nil {
		s.log.Errorf("Failed to write to SHA512 hasher", "error", err)
		return fmt.Errorf("failed to compute signature")
	}

	computedSignature := hex.EncodeToString(hasher.Sum(nil))

	if computedSignature != payload.SignatureKey {
		s.log.Warnf("Invalid Midtrans notification signature", "orderID", payload.OrderID, "computed", computedSignature, "received", payload.SignatureKey)
		return fmt.Errorf("invalid signature for order %s", payload.OrderID)
	}

	s.log.Infof("Midtrans notification signature verified successfully", "orderID", payload.OrderID)
	return nil
}

func (s *MidtransService) CreateQRISCharge(orderID string, amount int64) (*coreapi.ChargeResponse, error) {

	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeGopay,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
	}

	s.log.Infof("Creating QRIS charge for Order ID: %s with amount: %d", orderID, amount)

	chargeResp, err := s.client.ChargeTransaction(chargeReq)
	if err != nil {

		s.log.Errorf("Failed to create QRIS charge for Order ID: %s. Error: %v", orderID, err)
		return nil, err
	}

	s.log.Infof("Charge response: %+v", chargeResp)

	s.log.Infof("Successfully created QRIS charge for Order ID: %s. Transaction ID: %s", orderID, chargeResp.TransactionID)

	return chargeResp, nil
}

func (s *MidtransService) GetQRISCharge(orderID string) (*coreapi.TransactionStatusResponse, error) {
	resp, err := s.client.CheckTransaction(orderID)
	if err != nil {
		s.log.Errorf("Failed to check transaction status for Order ID: %s. Error: %v", orderID, err)
		return nil, err
	}

	s.log.Infof("Successfully retrieved transaction status for Order ID: %s. Status: %s", orderID, resp.TransactionStatus)
	return resp, nil
}

func (s *MidtransService) CancelTransaction(orderID string) (*coreapi.CancelResponse, error) {
	resp, err := s.client.CancelTransaction(orderID)
	if err != nil {
		s.log.Errorf("Failed to cancel transaction for Order ID: %s. Error: %v", orderID, err)
		return nil, err
	}

	s.log.Infof("Successfully cancelled transaction for Order ID: %s. Status: %s", orderID, resp.TransactionStatus)
	return resp, nil
}
