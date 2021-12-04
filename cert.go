package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

const (
	certValidityDuration time.Duration = 24 * time.Hour * 7 // 7 days
)

func SelfSignedCert(sans []string) (tls.Certificate, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to generate private key: %v", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(certValidityDuration)

	serialNumber, err := generateSerialNumber()
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to generate serial number: %v", err)
	}

	cert := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{Organization: []string{"fileserve self-signed"}},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	for _, san := range sans {
		if ip := net.ParseIP(san); ip != nil {
			cert.IPAddresses = append(cert.IPAddresses, ip)
		} else {
			cert.DNSNames = append(cert.DNSNames, strings.ToLower(san))
		}
	}

	der, err := x509.CreateCertificate(rand.Reader, &cert, &cert, &privKey.PublicKey, privKey)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("could not create certificate: %v", err)
	}

	return tls.Certificate{
		Certificate: [][]byte{der},
		PrivateKey:  privKey,
		Leaf:        &cert,
	}, nil
}

func generateSerialNumber() (n *big.Int, err error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128) //nolint:gomnd  //.
	return rand.Int(rand.Reader, serialNumberLimit)
}

func sans() []string {
	result := []string{"localhost", "127.0.0.1"}

	hostname, err := os.Hostname()
	if err == nil {
		result = append(result, hostname, hostname+".local", "*."+hostname+".local", hostname+".lan", "*."+hostname+".lan", hostname+".home", "*."+hostname+".home")
	}

	return result
}
