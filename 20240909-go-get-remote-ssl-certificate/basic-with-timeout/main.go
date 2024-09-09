package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"time"
)

func GetSSLCert(host string, port string) (*x509.Certificate, error) {
	// Create a custom dialer with a timeout
	dialer := &net.Dialer{Timeout: 5 * time.Second}

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", net.JoinHostPort(host, port), conf)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Get the negotiated state and return the first certificate
	if len(conn.ConnectionState().PeerCertificates) == 0 {
		return nil, fmt.Errorf("no certificates found")
	}

	log.Printf("remote address: %v", conn.RemoteAddr())
	return conn.ConnectionState().PeerCertificates[0], nil
}

func main() {
	host := "jd.com"
	port := "443"

	cert, err := GetSSLCert(host, port)
	if err != nil {
		log.Fatalf("Error fetching SSL certificate: %v", err)
	}

	fmt.Printf("Subject: %s, Expiry: %v, SerialNumber: %x\n", cert.Subject, cert.NotAfter, cert.SerialNumber)

}
