package inventory

import "strings"

func NormalizeFindings(findings []Finding, metadata BusinessMetadata) []CryptoAsset {
	assets := make([]CryptoAsset, 0, len(findings))
	for _, finding := range findings {
		assets = append(assets, CryptoAsset{
			ApplicationName:       metadata.ApplicationName,
			SourceType:            finding.SourceType,
			FilePath:              finding.FilePath,
			LineNumber:            finding.LineNumber,
			CryptoUsageType:       InferCryptoUsageType(finding.MatchedText, finding.RuleName),
			Algorithm:             finding.Algorithm,
			Severity:              finding.Severity,
			RiskType:              finding.RiskType,
			Confidence:            finding.Confidence,
			MatchedText:           finding.MatchedText,
			SensitivityFlags:      metadata.SensitivityFlags,
			ExposureLevel:         metadata.ExposureLevel,
			BusinessCriticality:   metadata.BusinessCriticality,
			SecrecyLifetimeYears:  metadata.SecrecyLifetimeYears,
			CryptoAgilityLevel:    metadata.CryptoAgilityLevel,
			VendorDependencyLevel: metadata.VendorDependencyLevel,
			MigrationComplexity:   metadata.MigrationComplexity,
		})
	}
	return assets
}

func InferCryptoUsageType(matchedText string, ruleName string) string {
	value := strings.ToLower(matchedText + " " + ruleName)
	switch {
	case strings.Contains(value, "signingcredentials") ||
		strings.Contains(value, "rsasecuritykey"):
		return "JWT signing"
	case strings.Contains(value, "jwtsecuritytokenhandler") ||
		strings.Contains(value, "jwtsecuritytoken"):
		return "JWT token handling"
	case strings.Contains(value, "jwt") || strings.Contains(value, "rs256") || strings.Contains(value, "es256"):
		return "JWT signing"
	case strings.Contains(value, "importrsaprivatekey"):
		return "private key import"
	case strings.Contains(value, "importsubjectpublickeyinfo"):
		return "public key import"
	case strings.Contains(value, "from crypto.publickey import rsa") ||
		strings.Contains(value, "cryptography rsa import") ||
		strings.Contains(value, "pycryptodome rsa import"):
		return "crypto library import"
	case strings.Contains(value, "rsa_private_key.pem") ||
		strings.Contains(value, "rsa_public_key.pem") ||
		strings.Contains(value, "private_key.pem") ||
		strings.Contains(value, "public_key.pem"):
		return "key file reference"
	case strings.Contains(value, "rsa.import_key") ||
		strings.Contains(value, "import_key") ||
		strings.Contains(value, "key import"):
		return "key import/loading"
	case strings.Contains(value, "keypairgenerator"):
		return "key generation"
	case strings.Contains(value, "rsa.create"):
		return "key generation"
	case strings.Contains(value, "pkcs1_15") || strings.Contains(value, "crypto.signature"):
		return "digital signature"
	case strings.Contains(value, "pkcs1_oaep"):
		return "encryption/decryption"
	case strings.Contains(value, "createsign") || strings.Contains(value, "sign"):
		return "digital signature"
	case strings.Contains(value, "encrypt") || strings.Contains(value, "decrypt"):
		return "encryption/decryption"
	case strings.Contains(value, "tls") || strings.Contains(value, "ssl") || strings.Contains(value, "ecdhe"):
		return "TLS/SSL configuration"
	case strings.Contains(value, "private key"):
		return "private key material"
	default:
		return "unknown crypto usage"
	}
}
