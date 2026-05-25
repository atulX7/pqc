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
	case strings.Contains(value, "jwt") || strings.Contains(value, "rs256") || strings.Contains(value, "es256"):
		return "JWT signing"
	case strings.Contains(value, "keypairgenerator"):
		return "key generation"
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
