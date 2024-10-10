package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"packagelock/logger"
	"time"
)

// CreateSelfSignedCert generates a self-signed RSA certificate and private key
func CreateSelfSignedCert(certFile, keyFile string) error {
	// Generate a private key using RSA (2048-bit key size)
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logger.Logger.Warnf("failed to generate private key: %v", err)
		return err
	}

	// Create a certificate template
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // 1 year

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		logger.Logger.Warnf("failed to generate serial number: %v", err)
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Self-Signed Co"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	err = os.MkdirAll("certs/", os.ModePerm)
	if err != nil {
		logger.Logger.Panicf("Cannot create 'certs/' directory, got: %s", err)
	}

	// Generate a self-signed certificate using the RSA private key
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		logger.Logger.Warnf("failed to create certificate: %v", err)
		return err
	}

	// Save the certificate to certFile
	certOut, err := os.Create(certFile)
	if err != nil {
		logger.Logger.Warnf("failed to open cert.pem for writing: %v", err)
	}

	// INFO: If the parrent throws an err and this defer is called
	// and fileOut.Close() throws an error to, the original error will be overwritten.
	defer func() {
		deferredErr := certOut.Close()
		if deferredErr != nil {
			logger.Logger.Warnf("Cannot close Cert File, got: %s", deferredErr)
			return
		}
	}()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		logger.Logger.Warnf("failed to write certificate to cert.pem: %v", err)
		return err
	}

	// Save the RSA private key to keyFile
	keyOut, err := os.Create(keyFile)
	if err != nil {
		logger.Logger.Warnf("failed to open key.pem for writing: %v", err)
		return err
	}

	// INFO: If the parrent throws an err and this defer is called
	// and fileOut.Close() throws an error to, the original error will be overwritten.
	defer func() {
		deferredErr := keyOut.Close()
		if deferredErr != nil {
			logger.Logger.Warnf("Cannot close Cert File, got: %s", deferredErr)
			return
		}
	}()

	// Marshal the RSA private key
	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}); err != nil {
		logger.Logger.Warnf("failed to write private key to key.pem: %v", err)
		return err
	}

	logger.Logger.Infof("Successfully created self-signed RSA certificate and private key.\n%s \n%s", certFile, keyFile)
	return nil
}
