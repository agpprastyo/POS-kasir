package printer

import (
	"POS-kasir/internal/orders"
	"POS-kasir/internal/repository"
	"POS-kasir/internal/settings"
	"POS-kasir/pkg/escpos"
	"POS-kasir/pkg/logger"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type IPrinterService interface {
	PrintInvoice(ctx context.Context, orderID uuid.UUID) error
	TestPrint(ctx context.Context) error
}

type PrinterService struct {
	orderService    orders.IOrderService
	settingsService settings.ISettingsService
	store           repository.Store
	log             logger.ILogger
}

func NewPrinterService(orderService orders.IOrderService, settingsService settings.ISettingsService, store repository.Store, log logger.ILogger) IPrinterService {
	return &PrinterService{
		orderService:    orderService,
		settingsService: settingsService,
		store:           store,
		log:             log,
	}
}

func (s *PrinterService) PrintInvoice(ctx context.Context, orderID uuid.UUID) error {
	// 1. Fetch Data
	order, err := s.orderService.GetOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	branding, err := s.settingsService.GetBranding(ctx)
	if err != nil {
		return fmt.Errorf("failed to get branding: %w", err)
	}

	printerSettings, err := s.settingsService.GetPrinterSettings(ctx)
	if err != nil {
		return fmt.Errorf("failed to get printer settings: %w", err)
	}

	// 2a. Fetch Cashier Name
	var cashierName string = "Unknown"
	if order.UserID != nil {
		user, err := s.store.GetUserByID(ctx, *order.UserID)
		if err == nil {
			cashierName = user.Username
		} else {
			s.log.Warn("Failed to fetch cashier name", "userID", order.UserID, "error", err)
		}
	}

	// 2. Connect to Printer
	p, err := escpos.NewPrinter(printerSettings.Connection)
	if err != nil {
		s.log.Error("Failed to connect to printer", "connection", printerSettings.Connection, "error", err)
		return err
	}
	defer p.Close()

	if err := p.Init(); err != nil {
		return err
	}

	// 3. Format Receipt
	// Header
	p.SetAlign(escpos.AlignCenter)
	p.SetBold(true)
	p.SetSize(escpos.DoubleHeightOn)
	p.WriteString(branding.AppName + "\n")
	p.SetSize(escpos.NormalSize)
	p.SetBold(false)
	p.WriteString(fmt.Sprintf("%s\n", time.Now().Format("02 Jan 2006 15:04")))
	p.WriteString(fmt.Sprintf("Order #%s\n", order.ID.String()[len(order.ID.String())-4:])) // Short ID
	p.WriteString("Cashier: " + cashierName + "\n")
	p.WriteString("--------------------------------\n")

	// Items
	p.SetAlign(escpos.AlignLeft)
	for _, item := range order.Items {
		// Name
		p.WriteString(item.ProductName + "\n")

		// Qty x Price = Total
		qtyLine := fmt.Sprintf("%dx %s", item.Quantity, formatCurrency(item.PriceAtSale))
		totalLine := formatCurrency(item.Subtotal)

		// Simple padding for 58mm (approx 32 chars)
		padding := 32 - len(qtyLine) - len(totalLine)
		if padding < 1 {
			padding = 1
		}

		p.WriteString(qtyLine + strings.Repeat(" ", padding) + totalLine + "\n")

		// Options
		for _, opt := range item.Options {
			if opt.OptionName != "" {
				p.WriteString("   + " + opt.OptionName + "\n")
			}
		}
	}

	p.WriteString("--------------------------------\n")

	// Totals
	writeTotalLine(p, "Subtotal", formatCurrency(order.GrossTotal))
	if order.DiscountAmount > 0 {
		writeTotalLine(p, "Discount", "-"+formatCurrency(order.DiscountAmount))
	}
	p.SetBold(true)
	writeTotalLine(p, "TOTAL", formatCurrency(order.NetTotal))
	p.SetBold(false)

	p.WriteString("--------------------------------\n")

	// Payment Info
	if order.PaymentMethodID != nil {
		var paymentMethodName string = "Unknown"
		methods, err := s.store.ListPaymentMethods(ctx)
		if err == nil {
			for _, m := range methods {
				if m.ID == *order.PaymentMethodID {
					paymentMethodName = m.Name
					break
				}
			}
		}

		p.WriteString("Payment: " + paymentMethodName + "\n")

		if order.CashReceived != nil && *order.CashReceived > 0 {
			writeTotalLine(p, "Cash", formatCurrency(*order.CashReceived))
			if order.ChangeDue != nil {
				writeTotalLine(p, "Change", formatCurrency(*order.ChangeDue))
			}
		}
	} else {
		p.WriteString("UNPAID\n")
	}

	// Footer
	p.SetAlign(escpos.AlignCenter)
	p.WriteString("\n")
	p.WriteString("© 2025 " + branding.AppName + "\n")
	p.WriteString("All rights reserved\n")
	p.WriteString("\n")
	// Cut
	return p.Cut()
}

func (s *PrinterService) TestPrint(ctx context.Context) error {
	printerSettings, err := s.settingsService.GetPrinterSettings(ctx)
	if err != nil {
		return fmt.Errorf("failed to get printer settings: %w", err)
	}

	p, err := escpos.NewPrinter(printerSettings.Connection)
	if err != nil {
		s.log.Error("Failed to connect to printer", "connection", printerSettings.Connection, "error", err)
		return err
	}
	defer p.Close()

	if err := p.Init(); err != nil {
		return err
	}

	p.SetAlign(escpos.AlignCenter)
	p.SetBold(true)
	p.WriteString("TEST PRINT SUCCESS\n")
	p.SetBold(false)
	p.WriteString("POS Kasir System\n")
	p.WriteString(time.Now().Format("02 Jan 2006 15:04:05") + "\n")
	p.WriteString("\n\n")

	return p.Cut()
}

func writeTotalLine(p *escpos.Printer, label, value string) {
	padding := 32 - len(label) - len(value)
	if padding < 1 {
		padding = 1
	}
	p.WriteString(label + strings.Repeat(" ", padding) + value + "\n")
}

func formatCurrency(amount int64) string {
	return fmt.Sprintf("Rp %d", amount)
}
