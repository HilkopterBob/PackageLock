package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type CertGeneratorParams struct {
	fx.In

	Logger *zap.Logger
}

type CertGenerator struct {
	logger *zap.Logger
}

func NewCertGenerator(params CertGeneratorParams) *CertGenerator {
	return &CertGenerator{
		logger: params.Logger,
	}
}

// CreateSelfSignedCert generates a self-signed RSA certificate and private key.
func (cg *CertGenerator) CreateSelfSignedCert(certFile, keyFile string) error {
	// Generate a private key using RSA (2048-bit key size)
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		cg.logger.Warn("Failed to generate private key", zap.Error(err))
		return err
	}

	// Create a certificate template
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // 1 year

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		cg.logger.Warn("Failed to generate serial number", zap.Error(err))
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

	// Ensure the "certs" directory exists
	err = os.MkdirAll("certs/", os.ModePerm)
	if err != nil {
		cg.logger.Fatal("Cannot create 'certs/' directory", zap.Error(err))
		return err
	}

	// Generate a self-signed certificate using the RSA private key
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		cg.logger.Warn("Failed to create certificate", zap.Error(err))
		return err
	}

	// Save the certificate to certFile
	certOut, err := os.Create(certFile)
	if err != nil {
		cg.logger.Warn("Failed to open cert.pem for writing", zap.Error(err))
		return err
	}
	defer func() {
		if err := certOut.Close(); err != nil {
			cg.logger.Warn("Cannot close certificate file", zap.Error(err))
		}
	}()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		cg.logger.Warn("Failed to write certificate to file", zap.Error(err))
		return err
	}

	// Save the RSA private key to keyFile
	keyOut, err := os.Create(keyFile)
	if err != nil {
		cg.logger.Warn("Failed to open key.pem for writing", zap.Error(err))
		return err
	}
	defer func() {
		if err := keyOut.Close(); err != nil {
			cg.logger.Warn("Cannot close key file", zap.Error(err))
		}
	}()

	// Marshal the RSA private key
	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}); err != nil {
		cg.logger.Warn("Failed to write private key to file", zap.Error(err))
		return err
	}

	cg.logger.Info("Successfully created self-signed RSA certificate and private key",
		zap.String("certFile", certFile),
		zap.String("keyFile", keyFile),
	)
	return nil
}

// Module exports the certs module.
var Module = fx.Options(
	fx.Provide(NewCertGenerator),
)
