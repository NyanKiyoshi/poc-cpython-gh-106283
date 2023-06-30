package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http/httptrace"
	"os"
	"time"
)

// createHttpTracer creates a client tracer object that prints debug info
// about what the dialer is doing.
func createHttpTracer() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		// When the connection to the server is initiated.
		ConnectStart: func(network, addr string) {
			log.Printf("Attempting to connect to %s...", addr)
		},
		// When the DNS resolution is complete.
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			log.Printf("Found %d IPs to try out", len(dnsInfo.Addrs))
		},
	}
}

// main runs a Go HTTP Client against user-provided flags which shows the Golang
// behavior against a domain name returning many IP addresses.
func main() {
	log.SetFlags(log.Lmicroseconds)
	timeout := flag.Duration("timeout", time.Second*10, "The connect timeout value to set (defaults to 1s)")
	domain := flag.String("domain", "poc.test", "The domain to connect to (defaults to poc.test)")
	flag.Parse()

	ctx := httptrace.WithClientTrace(context.Background(), createHttpTracer())
	dialer := net.Dialer{Timeout: *timeout}

	// Attempt to connect to the domain name
	// and measure the total time it takes to give-up.
	t0 := time.Now()
	conn, err := dialer.DialContext(ctx, "tcp", *domain+":80")

	if conn != nil {
		defer conn.Close()
	}

	t1 := time.Now()
	if err == nil {
		// counterintuitive: we expect and want to fail,
		// if we succeed then the domain is most likely not malicious.
		log.Printf("Unexpectedly succeeded to connect to remote address.")
		os.Exit(1)
	}
	log.Printf("Failed to get URL (in %s): %v", t1.Sub(t0), err)
}
