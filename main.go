package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	hub := NewHub()
	handler := newServer(hub)

	ip := getLANIP()
	port := "8080"
	addr := "0.0.0.0:" + port

	url1 := "http://localhost:" + port
	url2 := "http://" + ip + ":" + port
	label1 := "  Player 1 (host)  "
	label2 := "  Player 2 (LAN)   "
	innerWidth := len(label2) + len(url2) + 2
	if w := len(label1) + len(url1) + 2; w > innerWidth {
		innerWidth = w
	}
	bar := strings.Repeat("─", innerWidth)

	fmt.Println()
	fmt.Printf("  Tic-Tac-Toe\n")
	fmt.Printf("  %s\n", bar)
	fmt.Printf("%s%s\n", label1, url1)
	fmt.Printf("%s%s\n", label2, url2)
	fmt.Printf("  %s\n", bar)
	fmt.Printf("  Both players can use the LAN IP — even the host.\n")
	fmt.Printf("  Ctrl+C to stop the server.\n")
	fmt.Println()

	srv := &http.Server{Addr: addr, Handler: handler}

	// Handle OS interrupt/terminate signals.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down.")
		srv.Close()
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

// getLANIP returns the machine's outbound LAN IPv4 address.
// It dials a UDP address to let the OS pick the right interface — no packet is sent.
func getLANIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
