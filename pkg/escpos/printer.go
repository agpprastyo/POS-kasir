package escpos

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Printer interface {
	Close() error
	Write(data []byte) (int, error)
	WriteString(s string) (int, error)
	Init() error
	Cut() error
	Feed(n int) error
	SetAlign(align []byte) error
	SetBold(on bool) error
	SetSize(size []byte) error
}

type networkPrinter struct {
	conn net.Conn
}

func NewPrinter(connectionString string) (Printer, error) {
	// Clean up connection string
	address := strings.TrimPrefix(connectionString, "socket://")
	address = strings.TrimPrefix(address, "tcp://")

	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to printer: %w", err)
	}

	return &networkPrinter{conn: conn}, nil
}

func (p *networkPrinter) Close() error {
	return p.conn.Close()
}

func (p *networkPrinter) Write(data []byte) (int, error) {
	return p.conn.Write(data)
}

func (p *networkPrinter) WriteString(s string) (int, error) {
	return p.conn.Write([]byte(s))
}

func (p *networkPrinter) Init() error {
	_, err := p.Write(Init)
	return err
}

func (p *networkPrinter) Cut() error {
	// Feed a few lines before cutting
	p.Feed(3)
	_, err := p.Write(Cut)
	return err
}

func (p *networkPrinter) Feed(n int) error {
	for i := 0; i < n; i++ {
		if _, err := p.Write([]byte{LF}); err != nil {
			return err
		}
	}
	return nil
}

func (p *networkPrinter) SetAlign(align []byte) error {
	_, err := p.Write(align)
	return err
}

func (p *networkPrinter) SetBold(on bool) error {
	if on {
		_, err := p.Write(BoldOn)
		return err
	}
	_, err := p.Write(BoldOff)
	return err
}

func (p *networkPrinter) SetSize(size []byte) error {
	_, err := p.Write(size)
	return err
}
