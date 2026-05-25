package recommendation

import (
	"strings"
	"testing"

	"github.com/atulX7/pqc/internal/inventory"
)

func TestGenerateRecommendationCombinesRelevantGuidance(t *testing.T) {
	asset := inventory.CryptoAsset{
		Algorithm:             "RSA",
		SensitivityFlags:      []string{"PHI"},
		ExposureLevel:         "internet_public",
		CryptoAgilityLevel:    "hardcoded",
		VendorDependencyLevel: "legacy_no_roadmap",
		RiskLevel:             "Critical",
	}

	recommendation := Generate(asset)
	expectedPhrases := []string{
		"PQC migration inventory",
		"sensitive or regulated data",
		"externally exposed",
		"Improve crypto agility",
		"vendor PQC",
		"Phase 1 migration roadmap",
	}
	for _, phrase := range expectedPhrases {
		if !strings.Contains(recommendation, phrase) {
			t.Fatalf("expected recommendation to contain %q, got %q", phrase, recommendation)
		}
	}
}
