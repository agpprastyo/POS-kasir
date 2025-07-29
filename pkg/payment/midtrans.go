package payment

import (
	"POS-kasir/config"
	"POS-kasir/internal/dto"
	"POS-kasir/pkg/logger"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type IMidtrans interface {
	CreateQRISCharge(orderID string, amount int64) (*coreapi.ChargeResponse, error)
	VerifyNotificationSignature(payload dto.MidtransNotificationPayload) error
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
func (s *MidtransService) VerifyNotificationSignature(payload dto.MidtransNotificationPayload) error {
	// 1. Buat string sumber sesuai dokumentasi Midtrans: order_id + status_code + gross_amount + server_key
	sourceString := payload.OrderID + payload.StatusCode + payload.GrossAmount + s.config.Midtrans.ServerKey

	// 2. Hitung hash SHA-512 dari string sumber
	hasher := sha512.New()
	_, err := hasher.Write([]byte(sourceString))
	if err != nil {
		s.log.Errorf("Failed to write to SHA512 hasher", "error", err)
		return fmt.Errorf("failed to compute signature")
	}

	computedSignature := hex.EncodeToString(hasher.Sum(nil))

	// 3. Bandingkan signature yang dihitung dengan yang diterima dari Midtrans
	if computedSignature != payload.SignatureKey {
		s.log.Warnf("Invalid Midtrans notification signature", "orderID", payload.OrderID, "computed", computedSignature, "received", payload.SignatureKey)
		return fmt.Errorf("invalid signature for order %s", payload.OrderID)
	}

	s.log.Infof("Midtrans notification signature verified successfully", "orderID", payload.OrderID)
	return nil
}

func (s *MidtransService) CreateQRISCharge(orderID string, amount int64) (*coreapi.ChargeResponse, error) {

	chargeReq := &coreapi.ChargeReq{

		PaymentType: coreapi.PaymentTypeQris,
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
