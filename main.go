package main

import (
	"context"
	"fmt"
	"ip-port-scanner/scanner"
	"time"
)

func main() {
	portScanner := scanner.New(time.Millisecond * 200)

	result, _ := portScanner.Scan(context.Background(), "87.255.8.179", []scanner.Port{
		{Protocol: scanner.ProtocolTCP, Number: 80},
		{Protocol: scanner.ProtocolTCP, Number: 5444},
		{Protocol: scanner.ProtocolUDP, Number: 5444},
		{Protocol: scanner.ProtocolUDP, Number: 1337},
	})

	for _, r := range result {

		if r.PortStatus == scanner.PortStatusOpen {
			fmt.Printf("Порт %s:%d открыт\n", scanner.GetStringProtocol(r.Port.Protocol), r.Port.Number)
		} else {
			fmt.Printf("Порт %s:%d закрыт\n", scanner.GetStringProtocol(r.Port.Protocol), r.Port.Number)
		}
	}
}
