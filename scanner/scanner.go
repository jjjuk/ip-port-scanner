package scanner

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type IPAddress string

type Protocol int

const (
	ProtocolTCP Protocol = iota
	ProtocolUDP
)

type Port struct {
	Protocol Protocol
	Number   int
}

type PortStatus int

const (
	PortStatusOpen PortStatus = iota
	PortStatusClosed
	PortStatusUnknown
)

type PortScanResult struct {
	Port       Port
	PortStatus PortStatus
}

type Scanner interface {
	Scan(ctx context.Context, target IPAddress, ports *[]Port) ([]PortScanResult, error)
}

type PortScanner struct {
	timeout time.Duration
}

func New(timeout time.Duration) *PortScanner {
	return &PortScanner{timeout: timeout}
}

func (s *PortScanner) getStringProtocol(p Protocol) string {
	if p == ProtocolTCP {
		return "tcp"
	}
	return "udp"
}

func dial(ctx context.Context, network, address string, timeout time.Duration) (net.Conn, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return net.DialTimeout(network, address, timeout)
}

func (s *PortScanner) checkPort(ctx context.Context, wg *sync.WaitGroup, ip IPAddress, port Port, result chan<- PortScanResult) {
	defer wg.Done()

	conn, err := dial(ctx, s.getStringProtocol(port.Protocol), fmt.Sprintf("%s:%d", ip, port.Number), s.timeout)

	if err == nil {
		conn.Close()
		result <- PortScanResult{Port: port, PortStatus: PortStatusOpen}
	} else {
		result <- PortScanResult{Port: port, PortStatus: PortStatusClosed}
		// TODO: Error handling to get UNKNOWN status if needed
	}
}

func (s *PortScanner) Scan(ctx context.Context, ip IPAddress, ports []Port) ([]PortScanResult, error) {
	var wg sync.WaitGroup

	resultCh := make(chan PortScanResult, len(ports))

	for _, port := range ports {
		wg.Add(1)
		go s.checkPort(ctx, &wg, ip, port, resultCh)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	results := make([]PortScanResult, 0, len(ports))
	for r := range resultCh {
		results = append(results, r)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return results, nil
	}

}
