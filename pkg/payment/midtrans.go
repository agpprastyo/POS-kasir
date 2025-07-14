package payment

import (
	"POS-kasir/config"
	"POS-kasir/pkg/logger"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type IMidtrans interface {
	CreateQRISCharge(orderID string, amount int64) (*coreapi.ChargeResponse, error)
}

type MidtransService struct {
	client coreapi.Client
	log    *logger.Logger
}

func NewMidtransService(cfg *config.AppConfig, log *logger.Logger) IMidtrans {
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

	log.Info("Midtrans client initialized successfully")

	return &MidtransService{
		client: client,
		log:    log,
	}
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
		s.log.Errorf("Failed to create Midtrans charge: %v", err)
		return nil, err
	}

	s.log.Infof("Successfully created QRIS charge for Order ID: %s. Transaction ID: %s", orderID, chargeResp.TransactionID)

	return chargeResp, nil
}
