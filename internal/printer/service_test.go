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
func allowAllLoggerCalls(mockLogger *mocks.MockFieldLogger) {
	mockLogger.EXPECT().Warn(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any()).AnyTimes()
}

func TestPrinterService_PrintInvoice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderService := mocks.NewMockIOrderService(ctrl)
	mockSettingsService := new(MockSettingsService)
	mockPayment := mocks.NewMockIPaymentMethodService(ctrl)
	mockUserRepo := mocks.NewMockUserRepo(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockPrinter := new(MockPrinter)

	// Factory that returns the mock printer
	printerFactory := func(conn string) (escpos.Printer, error) {
		if conn == "fail" {
			return nil, errors.New("connection failed")
		}
		return mockPrinter, nil
	}

	service := printer.NewPrinterService(mockOrderService, mockSettingsService, mockPayment, mockUserRepo, mockLogger, printerFactory)

	ctx := context.Background()
	orderID := uuid.New()
	userID := uuid.New()
	payMethodID := int32(1)

	branding := &settings.BrandingSettingsResponse{AppName: "Test App"}
	printerSettings := &settings.PrinterSettingsResponse{Connection: "tcp://127.0.0.1:9100"}

	t.Run("Success", func(t *testing.T) {
		allowAllLoggerCalls(mockLogger)

		order := orders.OrderDetailResponse{
			ID:              orderID,
			UserID:          &userID,
			Status:          orders_repo.OrderStatusPaid,
			GrossTotal:      50000,
			NetTotal:        50000,
			PaymentMethodID: &payMethodID,
			Items: []orders.OrderItemResponse{
				{ProductName: "Item 1", Quantity: 2, PriceAtSale: 10000, Subtotal: 20000},
				{ProductName: "Item 2", Quantity: 1, PriceAtSale: 30000, Subtotal: 30000},
			},
		}

		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(&order, nil)
		mockSettingsService.On("GetBranding", ctx).Return(branding, nil).Once()
		mockSettingsService.On("GetPrinterSettings", ctx).Return(printerSettings, nil).Once()
		mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{Username: "Cashier1"}, nil)
		mockPayment.EXPECT().ListPaymentMethods(ctx).Return([]payment_methods.PaymentMethodResponse{{ID: 1, Name: "Cash"}}, nil)

		// Expect printer calls
		mockPrinter.On("Init").Return(nil)
		mockPrinter.On("SetAlign", mock.Anything).Return(nil)
		mockPrinter.On("SetBold", mock.Anything).Return(nil)
		mockPrinter.On("SetSize", mock.Anything).Return(nil)
		mockPrinter.On("WriteString", mock.Anything).Return(0, nil)
		mockPrinter.On("Cut").Return(nil)
		mockPrinter.On("Close").Return(nil)

		err := service.PrintInvoice(ctx, orderID)
		assert.NoError(t, err)
		mockPrinter.AssertExpectations(t)
		mockSettingsService.AssertExpectations(t)
	})

	t.Run("GetOrderError", func(t *testing.T) {
		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(nil, errors.New("db error"))

		err := service.PrintInvoice(ctx, orderID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get order")
	})

	t.Run("PrinterConnectionError", func(t *testing.T) {
		allowAllLoggerCalls(mockLogger)
		order := orders.OrderDetailResponse{ID: orderID}

		mockOrderService.EXPECT().GetOrder(ctx, orderID).Return(&order, nil)
		mockSettingsService.On("GetBranding", ctx).Return(branding, nil).Once()
		mockSettingsService.On("GetPrinterSettings", ctx).Return(&settings.PrinterSettingsResponse{Connection: "fail"}, nil).Once()

		err := service.PrintInvoice(ctx, orderID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "connection failed")
	})
}

func TestPrinterService_TestPrint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSettingsService := new(MockSettingsService)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockPrinter := new(MockPrinter)

	printerFactory := func(conn string) (escpos.Printer, error) {
		return mockPrinter, nil
	}

	service := printer.NewPrinterService(nil, mockSettingsService, nil, nil, mockLogger, printerFactory) // nil for unused deps
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		printerSettings := &settings.PrinterSettingsResponse{Connection: "tcp://127.0.0.1:9100"}
		mockSettingsService.On("GetPrinterSettings", ctx).Return(printerSettings, nil).Once()

		mockPrinter.On("Init").Return(nil)
		mockPrinter.On("SetAlign", mock.Anything).Return(nil)
		mockPrinter.On("SetBold", mock.Anything).Return(nil)
		mockPrinter.On("WriteString", mock.Anything).Return(0, nil)
		mockPrinter.On("Cut").Return(nil)
		mockPrinter.On("Close").Return(nil)

		err := service.TestPrint(ctx)
		assert.NoError(t, err)
		mockPrinter.AssertExpectations(t)
	})
}
