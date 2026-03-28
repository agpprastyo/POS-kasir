package printer_test

import (
	"POS-kasir/internal/orders"
	orders_repo "POS-kasir/internal/orders/repository"
	"POS-kasir/internal/payment_methods"
	"POS-kasir/internal/printer"
	"POS-kasir/internal/settings"
	user_repo "POS-kasir/internal/user/repository"
	"POS-kasir/mocks"
	"POS-kasir/pkg/escpos"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/mock/gomock"
)

// Manual Mocks

// MockPrinter
type MockPrinter struct {
	mock.Mock
}

func (m *MockPrinter) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPrinter) Write(data []byte) (int, error) {
	args := m.Called(data)
	return args.Int(0), args.Error(1)
}

func (m *MockPrinter) WriteString(s string) (int, error) {
	args := m.Called(s)
	return args.Int(0), args.Error(1)
}

func (m *MockPrinter) Init() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPrinter) Cut() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPrinter) Feed(n int) error {
	args := m.Called(n)
	return args.Error(0)
}

func (m *MockPrinter) SetAlign(align []byte) error {
	args := m.Called(align)
	return args.Error(0)
}

func (m *MockPrinter) SetBold(on bool) error {
	args := m.Called(on)
	return args.Error(0)
}

func (m *MockPrinter) SetSize(size []byte) error {
	args := m.Called(size)
	return args.Error(0)
}

// MockSettingsService
type MockSettingsService struct {
	mock.Mock
}

func (m *MockSettingsService) GetBranding(ctx context.Context) (*settings.BrandingSettingsResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*settings.BrandingSettingsResponse), args.Error(1)
}

func (m *MockSettingsService) UpdateBranding(ctx context.Context, req settings.UpdateBrandingRequest) (*settings.BrandingSettingsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*settings.BrandingSettingsResponse), args.Error(1)
}

func (m *MockSettingsService) GetPrinterSettings(ctx context.Context) (*settings.PrinterSettingsResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*settings.PrinterSettingsResponse), args.Error(1)
}

func (m *MockSettingsService) UpdatePrinterSettings(ctx context.Context, req settings.UpdatePrinterSettingsRequest) (*settings.PrinterSettingsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*settings.PrinterSettingsResponse), args.Error(1)
}

func (m *MockSettingsService) UpdateLogo(ctx context.Context, data []byte, filename string, contentType string) (string, error) {
	args := m.Called(ctx, data, filename, contentType)
	return args.String(0), args.Error(1)
}

// Helper for logger mocks
func allowAllLoggerCalls(mockLogger *mocks.MockILogger) {
	mockLogger.EXPECT().Warn(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
}

// newMockPrinter creates a fresh MockPrinter with all printer calls allowed.
func newMockPrinter() *MockPrinter {
	mp := new(MockPrinter)
	mp.On("Init").Return(nil)
	mp.On("SetAlign", mock.Anything).Return(nil)
	mp.On("SetBold", mock.Anything).Return(nil)
	mp.On("SetSize", mock.Anything).Return(nil)
	mp.On("WriteString", mock.Anything).Return(0, nil)
	mp.On("Cut").Return(nil)
	mp.On("Close").Return(nil)
	return mp
}

func TestPrinterService_PrintInvoice(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	userID := uuid.New()
	payMethodID := int32(1)

	branding := &settings.BrandingSettingsResponse{AppName: "Test App"}
	printerSettings := &settings.PrinterSettingsResponse{Connection: "tcp://127.0.0.1:9100"}

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrderService := mocks.NewMockIOrderService(ctrl)
		mockSettingsService := new(MockSettingsService)
		mockPayment := mocks.NewMockIPaymentMethodService(ctrl)
		mockUserRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockPrinter := newMockPrinter()

		printerFactory := func(conn string) (escpos.Printer, error) { return mockPrinter, nil }
		service := printer.NewPrinterService(mockOrderService, mockSettingsService, mockPayment, mockUserRepo, mockLogger, printerFactory)

		allowAllLoggerCalls(mockLogger)

		order := orders.OrderDetailResponse{
			ID: orderID, UserID: &userID, Status: orders_repo.OrderStatusPaid,
			GrossTotal: 50000, NetTotal: 50000, PaymentMethodID: &payMethodID,
			Items: []orders.OrderItemResponse{
				{ProductName: "Item 1", Quantity: 2, PriceAtSale: 10000, Subtotal: 20000},
				{ProductName: "Item 2", Quantity: 1, PriceAtSale: 30000, Subtotal: 30000},
			},
		}

		mockSettingsService.On("GetPrinterSettings", ctx).Return(printerSettings, nil).Once()
		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(&order, nil)
		mockSettingsService.On("GetBranding", ctx).Return(branding, nil).Once()
		mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{Username: "Cashier1"}, nil)
		mockPayment.EXPECT().ListPaymentMethods(ctx).Return([]payment_methods.PaymentMethodResponse{{ID: 1, Name: "Cash"}}, nil)

		err := service.PrintInvoice(ctx, orderID)
		assert.NoError(t, err)
		mockSettingsService.AssertExpectations(t)
	})

	t.Run("GetPrinterSettingsError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockLogger := mocks.NewMockILogger(ctrl)
		mockSettingsService := new(MockSettingsService)
		service := printer.NewPrinterService(nil, mockSettingsService, nil, nil, mockLogger, nil)

		mockSettingsService.On("GetPrinterSettings", ctx).Return(nil, errors.New("settings db error")).Once()

		err := service.PrintInvoice(ctx, orderID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get printer settings")
	})

	t.Run("GetOrderError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrderService := mocks.NewMockIOrderService(ctrl)
		mockSettingsService := new(MockSettingsService)
		mockLogger := mocks.NewMockILogger(ctrl)
		printerFactory := func(conn string) (escpos.Printer, error) { return nil, nil }
		service := printer.NewPrinterService(mockOrderService, mockSettingsService, nil, nil, mockLogger, printerFactory)

		mockSettingsService.On("GetPrinterSettings", ctx).Return(printerSettings, nil).Once()
		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(nil, errors.New("db error"))

		err := service.PrintInvoice(ctx, orderID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get order")
	})

	t.Run("GetBrandingError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrderService := mocks.NewMockIOrderService(ctrl)
		mockSettingsService := new(MockSettingsService)
		mockLogger := mocks.NewMockILogger(ctrl)
		printerFactory := func(conn string) (escpos.Printer, error) { return nil, nil }
		service := printer.NewPrinterService(mockOrderService, mockSettingsService, nil, nil, mockLogger, printerFactory)

		order := orders.OrderDetailResponse{ID: orderID}
		mockSettingsService.On("GetPrinterSettings", ctx).Return(printerSettings, nil).Once()
		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(&order, nil)
		mockSettingsService.On("GetBranding", ctx).Return(nil, errors.New("branding error")).Once()

		err := service.PrintInvoice(ctx, orderID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get branding")
	})

	t.Run("PrinterConnectionError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrderService := mocks.NewMockIOrderService(ctrl)
		mockSettingsService := new(MockSettingsService)
		mockLogger := mocks.NewMockILogger(ctrl)
		allowAllLoggerCalls(mockLogger)

		printerFactory := func(conn string) (escpos.Printer, error) { return nil, errors.New("connection failed") }
		service := printer.NewPrinterService(mockOrderService, mockSettingsService, nil, nil, mockLogger, printerFactory)

		order := orders.OrderDetailResponse{ID: orderID}
		failSettings := &settings.PrinterSettingsResponse{Connection: "fail"}
		mockSettingsService.On("GetPrinterSettings", ctx).Return(failSettings, nil).Once()
		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(&order, nil)
		mockSettingsService.On("GetBranding", ctx).Return(branding, nil).Once()

		err := service.PrintInvoice(ctx, orderID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection failed")
	})

	t.Run("SkippedFE", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockLogger := mocks.NewMockILogger(ctrl)
		allowAllLoggerCalls(mockLogger)
		mockSettingsService := new(MockSettingsService)
		service := printer.NewPrinterService(nil, mockSettingsService, nil, nil, mockLogger, nil)

		feSettings := &settings.PrinterSettingsResponse{Connection: "tcp://...", PrintMethod: "FE"}
		mockSettingsService.On("GetPrinterSettings", ctx).Return(feSettings, nil).Once()

		err := service.PrintInvoice(ctx, orderID)
		assert.NoError(t, err)
	})

	t.Run("WithDiscountAndCash", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrderService := mocks.NewMockIOrderService(ctrl)
		mockSettingsService := new(MockSettingsService)
		mockPayment := mocks.NewMockIPaymentMethodService(ctrl)
		mockUserRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockPrinter := newMockPrinter()
		allowAllLoggerCalls(mockLogger)

		printerFactory := func(conn string) (escpos.Printer, error) { return mockPrinter, nil }
		service := printer.NewPrinterService(mockOrderService, mockSettingsService, mockPayment, mockUserRepo, mockLogger, printerFactory)

		cashReceived := int64(100000)
		changeDue := int64(5000)
		order := orders.OrderDetailResponse{
			ID: orderID, UserID: &userID, Status: orders_repo.OrderStatusPaid,
			GrossTotal: 100000, DiscountAmount: 5000, NetTotal: 95000,
			PaymentMethodID: &payMethodID,
			CashReceived:    &cashReceived,
			ChangeDue:       &changeDue,
			Items: []orders.OrderItemResponse{
				{
					ProductName: "Nasi Goreng", Quantity: 2, PriceAtSale: 25000, Subtotal: 50000,
					Options: []orders.OrderItemOptionResponse{{OptionName: "Pedas"}},
				},
				{ProductName: "Es Teh", Quantity: 5, PriceAtSale: 10000, Subtotal: 50000},
			},
		}

		mockSettingsService.On("GetPrinterSettings", ctx).Return(printerSettings, nil).Once()
		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(&order, nil)
		mockSettingsService.On("GetBranding", ctx).Return(branding, nil).Once()
		mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{Username: "Cashier1"}, nil)
		mockPayment.EXPECT().ListPaymentMethods(ctx).Return([]payment_methods.PaymentMethodResponse{{ID: 1, Name: "Cash"}}, nil)

		err := service.PrintInvoice(ctx, orderID)
		assert.NoError(t, err)
	})

	t.Run("NilUserAndPaymentMethod", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrderService := mocks.NewMockIOrderService(ctrl)
		mockSettingsService := new(MockSettingsService)
		mockPayment := mocks.NewMockIPaymentMethodService(ctrl)
		mockUserRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockPrinter := newMockPrinter()
		allowAllLoggerCalls(mockLogger)

		printerFactory := func(conn string) (escpos.Printer, error) { return mockPrinter, nil }
		service := printer.NewPrinterService(mockOrderService, mockSettingsService, mockPayment, mockUserRepo, mockLogger, printerFactory)

		order := orders.OrderDetailResponse{
			ID: orderID, UserID: nil, Status: orders_repo.OrderStatusOpen,
			GrossTotal: 20000, NetTotal: 20000, PaymentMethodID: nil,
			Items: []orders.OrderItemResponse{
				{ProductName: "Water", Quantity: 1, PriceAtSale: 20000, Subtotal: 20000},
			},
		}

		mockSettingsService.On("GetPrinterSettings", ctx).Return(printerSettings, nil).Once()
		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(&order, nil)
		mockSettingsService.On("GetBranding", ctx).Return(branding, nil).Once()

		err := service.PrintInvoice(ctx, orderID)
		assert.NoError(t, err)
	})

	t.Run("GetUserError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrderService := mocks.NewMockIOrderService(ctrl)
		mockSettingsService := new(MockSettingsService)
		mockPayment := mocks.NewMockIPaymentMethodService(ctrl)
		mockUserRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockPrinter := newMockPrinter()
		allowAllLoggerCalls(mockLogger)

		printerFactory := func(conn string) (escpos.Printer, error) { return mockPrinter, nil }
		service := printer.NewPrinterService(mockOrderService, mockSettingsService, mockPayment, mockUserRepo, mockLogger, printerFactory)

		order := orders.OrderDetailResponse{
			ID: orderID, UserID: &userID, Status: orders_repo.OrderStatusPaid,
			GrossTotal: 20000, NetTotal: 20000, PaymentMethodID: &payMethodID,
			Items: []orders.OrderItemResponse{
				{ProductName: "Water", Quantity: 1, PriceAtSale: 20000, Subtotal: 20000},
			},
		}

		mockSettingsService.On("GetPrinterSettings", ctx).Return(printerSettings, nil).Once()
		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(&order, nil)
		mockSettingsService.On("GetBranding", ctx).Return(branding, nil).Once()
		mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{}, errors.New("user not found"))
		mockPayment.EXPECT().ListPaymentMethods(ctx).Return([]payment_methods.PaymentMethodResponse{{ID: 1, Name: "Cash"}}, nil)

		err := service.PrintInvoice(ctx, orderID)
		assert.NoError(t, err) // Continues despite user fetch error
	})

	t.Run("ListPaymentMethodsError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrderService := mocks.NewMockIOrderService(ctrl)
		mockSettingsService := new(MockSettingsService)
		mockPayment := mocks.NewMockIPaymentMethodService(ctrl)
		mockUserRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockPrinter := newMockPrinter()
		allowAllLoggerCalls(mockLogger)

		printerFactory := func(conn string) (escpos.Printer, error) { return mockPrinter, nil }
		service := printer.NewPrinterService(mockOrderService, mockSettingsService, mockPayment, mockUserRepo, mockLogger, printerFactory)

		order := orders.OrderDetailResponse{
			ID: orderID, UserID: &userID, Status: orders_repo.OrderStatusPaid,
			GrossTotal: 20000, NetTotal: 20000, PaymentMethodID: &payMethodID,
			Items: []orders.OrderItemResponse{
				{ProductName: "Water", Quantity: 1, PriceAtSale: 20000, Subtotal: 20000},
			},
		}

		mockSettingsService.On("GetPrinterSettings", ctx).Return(printerSettings, nil).Once()
		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(&order, nil)
		mockSettingsService.On("GetBranding", ctx).Return(branding, nil).Once()
		mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{Username: "Cashier1"}, nil)
		mockPayment.EXPECT().ListPaymentMethods(ctx).Return(nil, errors.New("db error"))

		err := service.PrintInvoice(ctx, orderID)
		assert.NoError(t, err) // Payment name defaults to "Unknown" on error
	})

	t.Run("PaymentMethodNotMatched", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrderService := mocks.NewMockIOrderService(ctrl)
		mockSettingsService := new(MockSettingsService)
		mockPayment := mocks.NewMockIPaymentMethodService(ctrl)
		mockUserRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		mockPrinter := newMockPrinter()
		allowAllLoggerCalls(mockLogger)

		printerFactory := func(conn string) (escpos.Printer, error) { return mockPrinter, nil }
		service := printer.NewPrinterService(mockOrderService, mockSettingsService, mockPayment, mockUserRepo, mockLogger, printerFactory)

		wrongPayMethodID := int32(99)
		order := orders.OrderDetailResponse{
			ID: orderID, UserID: &userID, Status: orders_repo.OrderStatusPaid,
			GrossTotal: 20000, NetTotal: 20000, PaymentMethodID: &wrongPayMethodID,
			Items: []orders.OrderItemResponse{
				{ProductName: "Water", Quantity: 1, PriceAtSale: 20000, Subtotal: 20000},
			},
		}

		mockSettingsService.On("GetPrinterSettings", ctx).Return(printerSettings, nil).Once()
		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(&order, nil)
		mockSettingsService.On("GetBranding", ctx).Return(branding, nil).Once()
		mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{Username: "Cashier1"}, nil)
		mockPayment.EXPECT().ListPaymentMethods(ctx).Return([]payment_methods.PaymentMethodResponse{{ID: 1, Name: "Cash"}}, nil)

		err := service.PrintInvoice(ctx, orderID)
		assert.NoError(t, err) // Falls through loop without match, paymentMethodName stays "Unknown"
	})

	t.Run("InitError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrderService := mocks.NewMockIOrderService(ctrl)
		mockSettingsService := new(MockSettingsService)
		mockPayment := mocks.NewMockIPaymentMethodService(ctrl)
		mockUserRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		allowAllLoggerCalls(mockLogger)

		failPrinter := new(MockPrinter)
		failPrinter.On("Init").Return(errors.New("init error"))
		failPrinter.On("Close").Return(nil)

		printerFactory := func(conn string) (escpos.Printer, error) { return failPrinter, nil }
		service := printer.NewPrinterService(mockOrderService, mockSettingsService, mockPayment, mockUserRepo, mockLogger, printerFactory)

		order := orders.OrderDetailResponse{
			ID: orderID, GrossTotal: 10000, NetTotal: 10000,
			Items: []orders.OrderItemResponse{{ProductName: "X", Quantity: 1, PriceAtSale: 10000, Subtotal: 10000}},
		}

		mockSettingsService.On("GetPrinterSettings", ctx).Return(printerSettings, nil).Once()
		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(&order, nil)
		mockSettingsService.On("GetBranding", ctx).Return(branding, nil).Once()

		err := service.PrintInvoice(ctx, orderID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "init error")
	})
}

func TestPrinterService_TestPrint(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSettingsService := new(MockSettingsService)
		mockLogger := mocks.NewMockILogger(ctrl)

		mp := new(MockPrinter)
		mp.On("Init").Return(nil)
		mp.On("SetAlign", mock.Anything).Return(nil)
		mp.On("SetBold", mock.Anything).Return(nil)
		mp.On("WriteString", mock.Anything).Return(0, nil)
		mp.On("Cut").Return(nil)
		mp.On("Close").Return(nil)

		printerFactory := func(conn string) (escpos.Printer, error) { return mp, nil }
		service := printer.NewPrinterService(nil, mockSettingsService, nil, nil, mockLogger, printerFactory)

		printerSettings := &settings.PrinterSettingsResponse{Connection: "tcp://127.0.0.1:9100"}
		mockSettingsService.On("GetPrinterSettings", ctx).Return(printerSettings, nil).Once()

		err := service.TestPrint(ctx)
		assert.NoError(t, err)
		mp.AssertExpectations(t)
	})

	t.Run("GetPrinterSettingsError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSettingsService := new(MockSettingsService)
		mockLogger := mocks.NewMockILogger(ctrl)
		service := printer.NewPrinterService(nil, mockSettingsService, nil, nil, mockLogger, nil)

		mockSettingsService.On("GetPrinterSettings", ctx).Return(nil, errors.New("settings error")).Once()

		err := service.TestPrint(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get printer settings")
	})

	t.Run("PrinterConnectionError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSettingsService := new(MockSettingsService)
		mockLogger := mocks.NewMockILogger(ctrl)
		allowAllLoggerCalls(mockLogger)

		printerFactory := func(conn string) (escpos.Printer, error) { return nil, errors.New("conn fail") }
		service := printer.NewPrinterService(nil, mockSettingsService, nil, nil, mockLogger, printerFactory)

		printerSettings := &settings.PrinterSettingsResponse{Connection: "tcp://127.0.0.1:9100"}
		mockSettingsService.On("GetPrinterSettings", ctx).Return(printerSettings, nil).Once()

		err := service.TestPrint(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conn fail")
	})

	t.Run("InitError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSettingsService := new(MockSettingsService)
		mockLogger := mocks.NewMockILogger(ctrl)

		failPrinter := new(MockPrinter)
		failPrinter.On("Init").Return(errors.New("init failed"))
		failPrinter.On("Close").Return(nil)

		printerFactory := func(conn string) (escpos.Printer, error) { return failPrinter, nil }
		service := printer.NewPrinterService(nil, mockSettingsService, nil, nil, mockLogger, printerFactory)

		printerSettings := &settings.PrinterSettingsResponse{Connection: "tcp://127.0.0.1:9100"}
		mockSettingsService.On("GetPrinterSettings", ctx).Return(printerSettings, nil).Once()

		err := service.TestPrint(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "init failed")
	})
}

func TestPrinterService_GetInvoiceData(t *testing.T) {
	ctx := context.Background()
	orderID := uuid.New()
	userID := uuid.New()
	payMethodID := int32(1)

	branding := &settings.BrandingSettingsResponse{AppName: "Test App"}

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrderService := mocks.NewMockIOrderService(ctrl)
		mockSettingsService := new(MockSettingsService)
		mockPayment := mocks.NewMockIPaymentMethodService(ctrl)
		mockUserRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		allowAllLoggerCalls(mockLogger)

		printerFactory := func(conn string) (escpos.Printer, error) { return nil, nil }
		service := printer.NewPrinterService(mockOrderService, mockSettingsService, mockPayment, mockUserRepo, mockLogger, printerFactory)

		order := orders.OrderDetailResponse{
			ID: orderID, UserID: &userID, Status: orders_repo.OrderStatusPaid,
			GrossTotal: 50000, NetTotal: 50000, PaymentMethodID: &payMethodID,
			Items: []orders.OrderItemResponse{
				{ProductName: "Item 1", Quantity: 1, PriceAtSale: 50000, Subtotal: 50000},
			},
		}

		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(&order, nil)
		mockSettingsService.On("GetBranding", ctx).Return(branding, nil).Once()
		mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{Username: "Cashier1"}, nil)
		mockPayment.EXPECT().ListPaymentMethods(ctx).Return([]payment_methods.PaymentMethodResponse{{ID: 1, Name: "Cash"}}, nil)

		data, filename, err := service.GetInvoiceData(ctx, orderID)

		assert.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.Contains(t, filename, "invoice_")
		assert.Contains(t, string(data), "Test App")

		mockSettingsService.AssertExpectations(t)
	})

	t.Run("PrepareDataError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrderService := mocks.NewMockIOrderService(ctrl)
		mockSettingsService := new(MockSettingsService)
		mockLogger := mocks.NewMockILogger(ctrl)

		printerFactory := func(conn string) (escpos.Printer, error) { return nil, nil }
		service := printer.NewPrinterService(mockOrderService, mockSettingsService, nil, nil, mockLogger, printerFactory)

		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(nil, errors.New("order not found"))

		data, filename, err := service.GetInvoiceData(ctx, orderID)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Empty(t, filename)
	})

	t.Run("WithDiscountAndOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrderService := mocks.NewMockIOrderService(ctrl)
		mockSettingsService := new(MockSettingsService)
		mockPayment := mocks.NewMockIPaymentMethodService(ctrl)
		mockUserRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		allowAllLoggerCalls(mockLogger)

		printerFactory := func(conn string) (escpos.Printer, error) { return nil, nil }
		service := printer.NewPrinterService(mockOrderService, mockSettingsService, mockPayment, mockUserRepo, mockLogger, printerFactory)

		cashReceived := int64(50000)
		changeDue := int64(5000)
		order := orders.OrderDetailResponse{
			ID: orderID, UserID: &userID, Status: orders_repo.OrderStatusPaid,
			GrossTotal: 50000, DiscountAmount: 5000, NetTotal: 45000,
			PaymentMethodID: &payMethodID,
			CashReceived:    &cashReceived, ChangeDue: &changeDue,
			Items: []orders.OrderItemResponse{
				{
					ProductName: "Nasi Goreng", Quantity: 2, PriceAtSale: 25000, Subtotal: 50000,
					Options: []orders.OrderItemOptionResponse{
						{OptionName: "Pedas"},
						{OptionName: ""},
					},
				},
			},
		}

		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(&order, nil)
		mockSettingsService.On("GetBranding", ctx).Return(branding, nil).Once()
		mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{Username: "Cashier1"}, nil)
		mockPayment.EXPECT().ListPaymentMethods(ctx).Return([]payment_methods.PaymentMethodResponse{{ID: 1, Name: "Cash"}}, nil)

		data, filename, err := service.GetInvoiceData(ctx, orderID)
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
		assert.Contains(t, filename, "invoice_")
		content := string(data)
		assert.Contains(t, content, "Discount")
		assert.Contains(t, content, "Pedas")
	})

	t.Run("LongProductNamesPaddingEdge", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockOrderService := mocks.NewMockIOrderService(ctrl)
		mockSettingsService := new(MockSettingsService)
		mockPayment := mocks.NewMockIPaymentMethodService(ctrl)
		mockUserRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockILogger(ctrl)
		allowAllLoggerCalls(mockLogger)

		printerFactory := func(conn string) (escpos.Printer, error) { return nil, nil }
		service := printer.NewPrinterService(mockOrderService, mockSettingsService, mockPayment, mockUserRepo, mockLogger, printerFactory)

		// Very long product name to trigger padding < 1 edge case in writeTotalLine
		order := orders.OrderDetailResponse{
			ID: orderID, UserID: &userID, Status: orders_repo.OrderStatusPaid,
			GrossTotal: 999999999, NetTotal: 999999999, PaymentMethodID: &payMethodID,
			Items: []orders.OrderItemResponse{
				{
					ProductName: "A Very Long Product Name That Exceeds 32 Characters Easily", Quantity: 999, PriceAtSale: 999999, Subtotal: 999999999,
				},
			},
		}

		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(&order, nil)
		mockSettingsService.On("GetBranding", ctx).Return(branding, nil).Once()
		mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{Username: "Cashier1"}, nil)
		mockPayment.EXPECT().ListPaymentMethods(ctx).Return([]payment_methods.PaymentMethodResponse{{ID: 1, Name: "Cash"}}, nil)

		data, _, err := service.GetInvoiceData(ctx, orderID)
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
	})
}

func TestPrinterService_DiscoverPrinters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := mocks.NewMockILogger(ctrl)

	service := printer.NewPrinterService(nil, nil, nil, nil, mockLogger, nil)
	ctx := context.Background()

	// This calls the real DiscoverPrinters which scans the network.
	// It should not error on any machine with a network interface.
	_, err := service.DiscoverPrinters(ctx)
	if err != nil {
		// The only expected error is "no local network interfaces found" in CI/containers
		assert.Contains(t, err.Error(), "no local network interfaces found")
	}
	// We don't assert specific printers found — just that no panic occurs
}
