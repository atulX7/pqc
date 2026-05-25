package report

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/atulX7/pqc/internal/inventory"
)

type ScanSummary struct {
	FilesScanned    int      `json:"files_scanned"`
	FilesSkipped    int      `json:"files_skipped"`
	UnreadableFiles []string `json:"unreadable_files,omitempty"`
}

type CriticalAsset struct {
	Algorithm      string  `json:"algorithm"`
	Usage          string  `json:"usage"`
	File           string  `json:"file"`
	Line           int     `json:"line"`
	RiskScore      float64 `json:"risk_score"`
	RiskLevel      string  `json:"risk_level"`
	Recommendation string  `json:"recommendation"`
}

type JSONReport struct {
	ApplicationName         string                  `json:"application_name"`
	ScanSummary            ScanSummary             `json:"scan_summary"`
	TotalFindings          int                     `json:"total_findings"`
	FindingsByAlgorithm    map[string]int          `json:"findings_by_algorithm"`
	FindingsBySeverity     map[string]int          `json:"findings_by_severity"`
	TopCriticalAssets      []CriticalAsset         `json:"top_critical_assets"`
	FullCryptoInventory    []inventory.CryptoAsset `json:"full_crypto_inventory"`
	HighestQuantumRiskScore float64                 `json:"highest_quantum_risk_score"`
	RiskPosture            string                  `json:"risk_posture"`
	RecommendedNextSteps   []string                `json:"recommended_next_steps"`
}

func Build(applicationName string, summary ScanSummary, assets []inventory.CryptoAsset) JSONReport {
	report := JSONReport{
		ApplicationName:      applicationName,
		ScanSummary:         summary,
		TotalFindings:       len(assets),
		FindingsByAlgorithm: map[string]int{},
		FindingsBySeverity:  map[string]int{},
		FullCryptoInventory: assets,
	}

	for _, asset := range assets {
		report.FindingsByAlgorithm[asset.Algorithm]++
		report.FindingsBySeverity[asset.Severity]++
		if asset.QuantumRiskScore > report.HighestQuantumRiskScore {
			report.HighestQuantumRiskScore = asset.QuantumRiskScore
			report.RiskPosture = asset.RiskLevel
		}
	}
	report.TopCriticalAssets = topAssets(assets, 5)
	report.RecommendedNextSteps = nextSteps(report.RiskPosture, len(assets))
	return report
}

func WriteJSON(path string, data JSONReport) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0o644)
}

func topAssets(assets []inventory.CryptoAsset, limit int) []CriticalAsset {
	copyAssets := append([]inventory.CryptoAsset(nil), assets...)
	sort.SliceStable(copyAssets, func(i, j int) bool {
		return copyAssets[i].QuantumRiskScore > copyAssets[j].QuantumRiskScore
	})
	if len(copyAssets) > limit {
		copyAssets = copyAssets[:limit]
	}
	top := make([]CriticalAsset, 0, len(copyAssets))
	for _, asset := range copyAssets {
		top = append(top, CriticalAsset{
			Algorithm:      asset.Algorithm,
			Usage:          asset.CryptoUsageType,
			File:           asset.FilePath,
			Line:           asset.LineNumber,
			RiskScore:      asset.QuantumRiskScore,
			RiskLevel:      asset.RiskLevel,
			Recommendation: asset.Recommendation,
		})
	}
	return top
}

func nextSteps(posture string, totalFindings int) []string {
	if totalFindings == 0 {
		return []string{"No crypto usage was detected by MVP rules. Expand rules and rescan high-value repositories."}
	}
	steps := []string{
		"Validate findings with application owners and security engineering.",
		"Assign ownership for high-risk crypto assets and document system dependencies.",
		"Define target PQC or hybrid migration patterns for signatures, key exchange, TLS, and certificates.",
	}
	if posture == "Very High" || posture == "Critical" {
		steps = append(steps, "Create a Phase 1 migration roadmap for the highest-risk assets.")
	}
	return steps
}
