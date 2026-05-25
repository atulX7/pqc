package scanner

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/atulX7/pqc/internal/inventory"
)

type TLSDomainResult struct {
	Findings []inventory.Finding
	Errors   []string
}

func ScanTLSDomains(domains []string) TLSDomainResult {
	var result TLSDomainResult
	for _, domain := range domains {
		domain = normalizeDomain(domain)
		if domain == "" {
			continue
		}
		cert, err := fetchLeafCertificate(domain)
		if err != nil {
			result.Errors = append(result.Errors, domain+": "+err.Error())
			continue
		}
		result.Findings = append(result.Findings, FindingsFromCertificate(domain, cert)...)
	}
	return result
}

func fetchLeafCertificate(domain string) (*x509.Certificate, error) {
	host, _, err := net.SplitHostPort(domain)
	if err != nil {
		host = strings.TrimSuffix(domain, ":443")
		domain = net.JoinHostPort(host, "443")
	}
	dialer := net.Dialer{Timeout: 7 * time.Second}
	conn, err := tls.DialWithDialer(&dialer, "tcp", domain, &tls.Config{
		ServerName:         host,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false,
	})
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("no peer certificates")
	}
	return certs[0], nil
}

func FindingsFromCertificate(domain string, cert *x509.Certificate) []inventory.Finding {
	algorithm, keySize := publicKeyAlgorithmAndSize(cert.PublicKey)
	if algorithm == "UNKNOWN" {
		return nil
	}

	severity := "high"
	riskType := "quantum_vulnerable_public_key"
	if algorithm == "DSA" {
		severity = "critical"
		riskType = "deprecated_signature"
	}
	daysUntilExpiry := int(time.Until(cert.NotAfter).Hours() / 24)
	if daysUntilExpiry < 0 {
		severity = "critical"
	}

	return []inventory.Finding{{
		SourceType:  "tls_certificate",
		FilePath:    domain,
		LineNumber:  0,
		MatchedText: certificateSummary(cert, algorithm, keySize, daysUntilExpiry),
		RuleID:      "tls-certificate-public-key",
		RuleName:    "TLS certificate public key",
		Algorithm:   algorithm,
		Severity:    severity,
		RiskType:    riskType,
		Confidence:  "high",
	}}
}

func publicKeyAlgorithmAndSize(publicKey any) (string, int) {
	switch key := publicKey.(type) {
	case *rsa.PublicKey:
		return "RSA-" + strconv.Itoa(key.N.BitLen()), key.N.BitLen()
	case *ecdsa.PublicKey:
		return "ECDSA", key.Curve.Params().BitSize
	case ed25519.PublicKey:
		return "Ed25519", len(key) * 8
	case *dsa.PublicKey:
		return "DSA", key.P.BitLen()
	default:
		return "UNKNOWN", 0
	}
}

func certificateSummary(cert *x509.Certificate, algorithm string, keySize int, daysUntilExpiry int) string {
	return fmt.Sprintf("subject=%q issuer=%q public_key=%s key_size=%d not_after=%s days_until_expiry=%d",
		cert.Subject.String(),
		cert.Issuer.String(),
		algorithm,
		keySize,
		cert.NotAfter.Format(time.RFC3339),
		daysUntilExpiry,
	)
}

func normalizeDomain(domain string) string {
	domain = strings.TrimSpace(domain)
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimRight(domain, "/")
	return domain
}
