package recommendation

import (
	"strings"

	"github.com/atulX7/pqc/internal/inventory"
)

func AddRecommendations(assets []inventory.CryptoAsset) []inventory.CryptoAsset {
	recommended := make([]inventory.CryptoAsset, 0, len(assets))
	for _, asset := range assets {
		asset.Recommendation = Generate(asset)
		recommended = append(recommended, asset)
	}
	return recommended
}

func Generate(asset inventory.CryptoAsset) string {
	var parts []string
	algorithm := strings.ToUpper(asset.Algorithm)
	if isQuantumVulnerableAlgorithm(algorithm) {
		parts = append(parts, "Add this "+algorithm+" usage to the PQC migration inventory and assess hybrid or post-quantum alternatives.")
	}
	if hasSensitiveData(asset.SensitivityFlags) {
		parts = append(parts, "Prioritize this asset because sensitive or regulated data may require long-term confidentiality and compliance planning.")
	}
	if asset.ExposureLevel == "internet_authenticated" || asset.ExposureLevel == "internet_public" {
		parts = append(parts, "Review externally exposed TLS, certificate, signing, and key exchange dependencies for PQC migration readiness.")
	}
	if asset.CryptoAgilityLevel == "hardcoded" || asset.CryptoAgilityLevel == "vendor_controlled_or_legacy" {
		parts = append(parts, "Improve crypto agility with centralized KMS/HSM usage or a configurable crypto policy layer.")
	}
	if asset.VendorDependencyLevel == "saas_unknown_roadmap" || asset.VendorDependencyLevel == "legacy_no_roadmap" {
		parts = append(parts, "Request a vendor PQC and hybrid cryptography roadmap with timelines and supported algorithms.")
	}
	if asset.RiskLevel == "Very High" || asset.RiskLevel == "Critical" {
		parts = append(parts, "Place this item into a Phase 1 migration roadmap with ownership, target architecture, and validation milestones.")
	}
	if len(parts) == 0 {
		return "Track this crypto usage in the inventory and reassess as PQC migration standards and platform support mature."
	}
	return strings.Join(parts, " ")
}

func isQuantumVulnerableAlgorithm(algorithm string) bool {
	switch algorithm {
	case "RSA", "RSA-2048", "RSA-4096", "ECC", "ECDSA", "ECDH", "DH", "DSA":
		return true
	default:
		return false
	}
}

func hasSensitiveData(flags []string) bool {
	for _, flag := range flags {
		switch strings.ToLower(strings.TrimSpace(flag)) {
		case "phi", "gxp", "pharma_ip", "spi":
			return true
		}
	}
	return false
}
