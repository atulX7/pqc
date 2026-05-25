package scanner

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"
)

func TestFindingsFromCertificateExtractsRSAKeyMetadata(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "example.test"},
		Issuer:       pkix.Name{CommonName: "Test CA"},
		NotAfter:     time.Now().Add(90 * 24 * time.Hour),
		PublicKey:    &key.PublicKey,
	}

	findings := FindingsFromCertificate("example.test:443", cert)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	finding := findings[0]
	if finding.SourceType != "tls_certificate" {
		t.Fatalf("unexpected source type: %s", finding.SourceType)
	}
	if finding.Algorithm != "RSA-2048" {
		t.Fatalf("expected RSA-2048, got %s", finding.Algorithm)
	}
	if finding.Confidence != "high" {
		t.Fatalf("expected high confidence, got %s", finding.Confidence)
	}
	if finding.MatchedText == "" {
		t.Fatal("expected certificate metadata in matched text")
	}
}
