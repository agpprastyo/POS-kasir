package escpos

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Printer struct {
	conn net.Conn
}

func NewPrinter(connectionString string) (*Printer, error) {
	// Clean up connection string
	address := strings.TrimPrefix(connectionString, "socket://")
	address = strings.TrimPrefix(address, "tcp://")

	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to printer: %w", err)
	}

	return &Printer{conn: conn}, nil
}

func (p *Printer) Close() error {
	return p.conn.Close()
}

func (p *Printer) Write(data []byte) (int, error) {
	return p.conn.Write(data)
}

func (p *Printer) WriteString(s string) (int, error) {
	return p.conn.Write([]byte(s))
}

func (p *Printer) Init() error {
	_, err := p.Write(Init)
	return err
}

func (p *Printer) Cut() error {
	// Feed a few lines before cutting
	p.Feed(3)
	_, err := p.Write(Cut)
	return err
}

func (p *Printer) Feed(n int) error {
	for i := 0; i < n; i++ {
		if _, err := p.Write([]byte{LF}); err != nil {
			return err
		}
	}
	return nil
}

func (p *Printer) SetAlign(align []byte) error {
	_, err := p.Write(align)
	return err
}

func (p *Printer) SetBold(on bool) error {
	if on {
		_, err := p.Write(BoldOn)
		return err
	}
	_, err := p.Write(BoldOff)
	return err
}

func (p *Printer) SetSize(size []byte) error {
	_, err := p.Write(size)
	return err
}
