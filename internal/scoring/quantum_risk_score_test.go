package scoring

import (
	"testing"

	"github.com/atulX7/pqc/internal/inventory"
)

func TestCalculateQuantumRiskScoreCriticalScenario(t *testing.T) {
	asset := inventory.CryptoAsset{
		Algorithm:             "RSA",
		SensitivityFlags:      []string{"PHI", "GxP"},
		SecrecyLifetimeYears:  15,
		ExposureLevel:         "internet_authenticated",
		BusinessCriticality:   "high",
		CryptoAgilityLevel:    "partially_configurable",
		VendorDependencyLevel: "cloud_with_roadmap",
		MigrationComplexity:   "multi_system",
	}

	score := CalculateQuantumRiskScore(asset)
	if score != 83 {
		t.Fatalf("expected score 83, got %.2f", score)
	}
	if RiskLevel(score) != "Critical" {
		t.Fatalf("expected Critical, got %s", RiskLevel(score))
	}
}

func TestRiskLevelBoundaries(t *testing.T) {
	cases := map[float64]string{
		20: "Low",
		40: "Moderate",
		60: "High",
		80: "Very High",
		81: "Critical",
	}
	for score, expected := range cases {
		if actual := RiskLevel(score); actual != expected {
			t.Fatalf("score %.2f expected %s, got %s", score, expected, actual)
		}
	}
}
