package report

import (
	"path/filepath"
	"testing"

	"github.com/atulX7/pqc/internal/inventory"
)

func TestBuildIncludesExecutiveSummaryAndUsageCounts(t *testing.T) {
	assets := []inventory.CryptoAsset{{
		ApplicationName:  "Clinical API",
		CryptoUsageType:  "TLS/SSL configuration",
		Algorithm:        "RSA-2048",
		Severity:         "high",
		RiskType:         "quantum_vulnerable_public_key",
		QuantumRiskScore: 83,
		RiskLevel:        "Critical",
		Recommendation:   "Rotate certificate.",
		FilePath:         "example.com:443",
		Confidence:       "high",
	}}

	report := Build("Clinical API", ScanSummary{FilesScanned: 1}, assets)
	if report.ExecutiveSummary.OverallRiskPosture != "Critical" {
		t.Fatalf("unexpected executive posture: %s", report.ExecutiveSummary.OverallRiskPosture)
	}
	if report.FindingsByUsageType["TLS/SSL configuration"] != 1 {
		t.Fatalf("expected usage type count")
	}
	if len(report.ExecutiveSummary.Actions30Days) == 0 || len(report.TopCriticalAssets) != 1 {
		t.Fatalf("expected executive actions and top assets")
	}
}

func TestTopAssetsPrioritizesSeverityConfidenceAndUsageOnRiskTie(t *testing.T) {
	assets := []inventory.CryptoAsset{
		{
			CryptoUsageType:  "JWT token handling",
			Algorithm:        "RSA",
			Severity:         "low",
			Confidence:       "low",
			QuantumRiskScore: 83,
			RiskLevel:        "Critical",
			FilePath:         "jwt.fsx",
			LineNumber:       78,
		},
		{
			CryptoUsageType:  "private key material",
			Algorithm:        "RSA",
			Severity:         "critical",
			Confidence:       "high",
			QuantumRiskScore: 83,
			RiskLevel:        "Critical",
			FilePath:         "key.priv",
			LineNumber:       1,
		},
		{
			CryptoUsageType:  "JWT signing",
			Algorithm:        "RSA",
			Severity:         "high",
			Confidence:       "high",
			QuantumRiskScore: 83,
			RiskLevel:        "Critical",
			FilePath:         "jwt.fsx",
			LineNumber:       50,
		},
	}

	top := topAssets(assets, 3)
	if top[0].Usage != "private key material" {
		t.Fatalf("expected private key material first, got %+v", top)
	}
	if top[2].Usage != "JWT token handling" {
		t.Fatalf("expected JWT token handling last, got %+v", top)
	}
}

func TestWriteExcelCreatesWorkbook(t *testing.T) {
	report := Build("Clinical API", ScanSummary{FilesScanned: 1}, nil)
	path := filepath.Join(t.TempDir(), "report.xlsx")
	if err := WriteExcel(path, report); err != nil {
		t.Fatal(err)
	}
}
