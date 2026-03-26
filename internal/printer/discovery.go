package printer

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type DiscoveredPrinter struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
	Name string `json:"name"`
}

func DiscoverPrinters(ctx context.Context) ([]DiscoveredPrinter, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var networks []*net.IPNet
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				networks = append(networks, ipnet)
			}
		}
	}

	if len(networks) == 0 {
		return nil, fmt.Errorf("no local network interfaces found")
	}

	var found []DiscoveredPrinter
	var mu sync.Mutex
	var wg sync.WaitGroup

	// We only scan the first non-loopback IPv4 interface subnet for brevity and speed
	// Thermal printers are almost always on port 9100
	ipnet := networks[0]
	ip := ipnet.IP.To4()
	mask := ipnet.Mask

	// Iterate through all IPs in the /24 subnet (simplification)
	// Real subnet scanning would be more complex, but /24 is common.
	baseIP := ip.Mask(mask)
	
	// Max 255 concurrent scans
	for i := 1; i < 255; i++ {
		targetIP := make(net.IP, len(baseIP))
		copy(targetIP, baseIP)
		targetIP[3] = byte(i)

		wg.Add(1)
		go func(target string) {
			defer wg.Done()
			
			// Short timeout for discovery
			d := net.Dialer{Timeout: 300 * time.Millisecond}
			address := fmt.Sprintf("%s:9100", target)
			conn, err := d.DialContext(ctx, "tcp", address)
			if err == nil {
				conn.Close()
				mu.Lock()
				found = append(found, DiscoveredPrinter{
					IP:   target,
					Port: 9100,
					Name: fmt.Sprintf("Printer (%s)", target),
				})
				mu.Unlock()
			}
		}(targetIP.String())
	}

	wg.Wait()
	return found, nil
}
