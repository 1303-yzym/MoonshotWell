package cert

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// 加载根证书和私钥
func loadRootCertificate(certFile, keyFile string) (*x509.Certificate, *rsa.PrivateKey, error) {
	// 读取根证书
	certPEM, err := os.ReadFile(certFile)
	if err != nil {
		return nil, nil, err
	}

	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, nil, err
	}

	rootCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	// 读取根私钥
	keyPEM, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, nil, err
	}

	block, _ = pem.Decode(keyPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, nil, err
	}

	rootKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return rootCert, rootKey, nil
}
