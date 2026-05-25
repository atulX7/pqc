package scoring

import (
	"math"
	"strings"

	"github.com/atulX7/pqc/internal/inventory"
)

func ScoreAsset(asset inventory.CryptoAsset) inventory.CryptoAsset {
	score := CalculateQuantumRiskScore(asset)
	asset.QuantumRiskScore = score
	asset.RiskLevel = RiskLevel(score)
	return asset
}

func CalculateQuantumRiskScore(asset inventory.CryptoAsset) float64 {
	weightedTotal := 0.0
	weightedTotal += float64(algorithmScore(asset.Algorithm)) * 25
	weightedTotal += float64(sensitivityScore(asset.SensitivityFlags)) * 20
	weightedTotal += float64(secrecyLifetimeScore(asset.SecrecyLifetimeYears)) * 15
	weightedTotal += float64(lookupScore(asset.ExposureLevel, exposureScores, 2)) * 10
	weightedTotal += float64(lookupScore(asset.BusinessCriticality, criticalityScores, 1)) * 10
	weightedTotal += float64(lookupScore(asset.CryptoAgilityLevel, agilityScores, 3)) * 10
	weightedTotal += float64(lookupScore(asset.VendorDependencyLevel, vendorScores, 1)) * 5
	weightedTotal += float64(lookupScore(asset.MigrationComplexity, complexityScores, 3)) * 5
	return math.Round((weightedTotal/5)*100) / 100
}

func RiskLevel(score float64) string {
	switch {
	case score <= 20:
		return "Low"
	case score <= 40:
		return "Moderate"
	case score <= 60:
		return "High"
	case score <= 80:
		return "Very High"
	default:
		return "Critical"
	}
}

func ScoreAssets(assets []inventory.CryptoAsset) []inventory.CryptoAsset {
	scored := make([]inventory.CryptoAsset, 0, len(assets))
	for _, asset := range assets {
		scored = append(scored, ScoreAsset(asset))
	}
	return scored
}

func algorithmScore(algorithm string) int {
	switch strings.ToUpper(strings.TrimSpace(algorithm)) {
	case "ML-KEM", "ML-DSA", "SLH-DSA":
		return 0
	case "AES", "SHA256":
		return 1
	case "RSA-4096":
		return 3
	case "RSA", "RSA-2048", "ECC", "ECDSA", "ECDH", "DH":
		return 4
	case "SHA1":
		return 4
	case "DSA", "MD5", "UNKNOWN", "CUSTOM":
		return 5
	default:
		return 5
	}
}

func sensitivityScore(flags []string) int {
	maxScore := 0
	for _, flag := range flags {
		score := lookupScore(flag, sensitivityScores, 0)
		if score > maxScore {
			maxScore = score
		}
	}
	return maxScore
}

func secrecyLifetimeScore(years int) int {
	switch {
	case years < 1:
		return 1
	case years <= 3:
		return 2
	case years <= 7:
		return 3
	case years < 15:
		return 4
	default:
		return 5
	}
}

func lookupScore(value string, scores map[string]int, fallback int) int {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.ReplaceAll(normalized, "-", "_")
	if score, ok := scores[normalized]; ok {
		return score
	}
	return fallback
}

var sensitivityScores = map[string]int{
	"public": 0, "internal": 1, "confidential": 3, "pii": 4, "spi": 5,
	"phi": 5, "gxp": 5, "pharma_ip": 5, "financial": 4, "legal": 4,
}

var exposureScores = map[string]int{
	"offline": 1, "internal": 2, "partner": 3, "internet_authenticated": 4, "internet_public": 5,
}

var criticalityScores = map[string]int{
	"low": 1, "medium": 3, "high": 4, "mission_critical": 5,
}

var agilityScores = map[string]int{
	"centralized_policy_driven": 1, "configurable": 2, "partially_configurable": 3,
	"hardcoded": 4, "vendor_controlled_or_legacy": 5,
}

var vendorScores = map[string]int{
	"none": 1, "cloud_with_roadmap": 2, "saas_unknown_roadmap": 4, "legacy_no_roadmap": 5,
}

var complexityScores = map[string]int{
	"config_change": 1, "library_upgrade": 2, "code_change": 3,
	"multi_system": 4, "legacy_regulated_vendor_bound": 5,
}
