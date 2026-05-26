package report

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/atulX7/pqc/internal/inventory"
	"github.com/xuri/excelize/v2"
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
	ExecutiveSummary        ExecutiveSummary        `json:"executive_summary"`
	ScanSummary             ScanSummary             `json:"scan_summary"`
	TotalFindings           int                     `json:"total_findings"`
	FindingsByAlgorithm     map[string]int          `json:"findings_by_algorithm"`
	FindingsBySeverity      map[string]int          `json:"findings_by_severity"`
	FindingsByUsageType     map[string]int          `json:"findings_by_usage_type"`
	TopCriticalAssets       []CriticalAsset         `json:"top_critical_assets"`
	FullCryptoInventory     []inventory.CryptoAsset `json:"full_crypto_inventory"`
	HighestQuantumRiskScore float64                 `json:"highest_quantum_risk_score"`
	RiskPosture             string                  `json:"risk_posture"`
	RecommendedNextSteps    []string                `json:"recommended_next_steps"`
}

type ExecutiveSummary struct {
	OverallRiskPosture  string          `json:"overall_risk_posture"`
	TopCriticalAssets   []CriticalAsset `json:"top_10_critical_assets"`
	FindingsByAlgorithm map[string]int  `json:"findings_by_algorithm"`
	FindingsByUsageType map[string]int  `json:"findings_by_usage_type"`
	Actions30Days       []string        `json:"actions_30_days"`
	Actions60Days       []string        `json:"actions_60_days"`
	Actions90Days       []string        `json:"actions_90_days"`
}

func Build(applicationName string, summary ScanSummary, assets []inventory.CryptoAsset) JSONReport {
	report := JSONReport{
		ApplicationName:     applicationName,
		ScanSummary:         summary,
		TotalFindings:       len(assets),
		FindingsByAlgorithm: map[string]int{},
		FindingsBySeverity:  map[string]int{},
		FindingsByUsageType: map[string]int{},
		FullCryptoInventory: assets,
	}

	for _, asset := range assets {
		report.FindingsByAlgorithm[asset.Algorithm]++
		report.FindingsBySeverity[asset.Severity]++
		report.FindingsByUsageType[asset.CryptoUsageType]++
		if asset.QuantumRiskScore > report.HighestQuantumRiskScore {
			report.HighestQuantumRiskScore = asset.QuantumRiskScore
			report.RiskPosture = asset.RiskLevel
		}
	}
	if len(assets) == 0 {
		report.RiskPosture = "Low"
	}
	report.RecommendedNextSteps = nextSteps(report.RiskPosture, len(assets))
	report.TopCriticalAssets = topAssets(assets, 10)
	report.ExecutiveSummary = buildExecutiveSummary(report)
	return report
}

func WriteJSON(path string, data JSONReport) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0o644)
}

func WriteExcel(path string, data JSONReport) error {
	file := excelize.NewFile()
	defer file.Close()
	summarySheet := "Executive Summary"
	file.SetSheetName("Sheet1", summarySheet)
	writeSummarySheet(file, summarySheet, data)
	writeTopAssetsSheet(file, data)
	writeInventorySheet(file, data)
	return file.SaveAs(path)
}

func writeSummarySheet(file *excelize.File, sheet string, data JSONReport) {
	rows := [][]any{
		{"Application", data.ApplicationName},
		{"Overall risk posture", data.ExecutiveSummary.OverallRiskPosture},
		{"Highest quantum risk score", data.HighestQuantumRiskScore},
		{"Total findings", data.TotalFindings},
		{"Files scanned", data.ScanSummary.FilesScanned},
		{"Files skipped", data.ScanSummary.FilesSkipped},
		{},
		{"Findings by algorithm"},
	}
	row := 1
	for _, values := range rows {
		setRow(file, sheet, row, values)
		row++
	}
	for _, key := range sortedKeys(data.FindingsByAlgorithm) {
		setRow(file, sheet, row, []any{key, data.FindingsByAlgorithm[key]})
		row++
	}
	row++
	setRow(file, sheet, row, []any{"30-day actions"})
	row++
	for _, action := range data.ExecutiveSummary.Actions30Days {
		setRow(file, sheet, row, []any{action})
		row++
	}
	row++
	setRow(file, sheet, row, []any{"60-day actions"})
	row++
	for _, action := range data.ExecutiveSummary.Actions60Days {
		setRow(file, sheet, row, []any{action})
		row++
	}
	row++
	setRow(file, sheet, row, []any{"90-day actions"})
	row++
	for _, action := range data.ExecutiveSummary.Actions90Days {
		setRow(file, sheet, row, []any{action})
		row++
	}
	_ = file.SetColWidth(sheet, "A", "B", 28)
}

func writeTopAssetsSheet(file *excelize.File, data JSONReport) {
	sheet := "Top Assets"
	_, _ = file.NewSheet(sheet)
	setRow(file, sheet, 1, []any{"Algorithm", "Usage", "File", "Line", "Risk Score", "Risk Level", "Recommendation"})
	for index, asset := range data.TopCriticalAssets {
		setRow(file, sheet, index+2, []any{asset.Algorithm, asset.Usage, asset.File, asset.Line, asset.RiskScore, asset.RiskLevel, asset.Recommendation})
	}
	_ = file.SetColWidth(sheet, "A", "G", 24)
}

func writeInventorySheet(file *excelize.File, data JSONReport) {
	sheet := "Inventory"
	_, _ = file.NewSheet(sheet)
	setRow(file, sheet, 1, []any{"Application", "Source Type", "File/Domain", "Line", "Usage", "Algorithm", "Severity", "Risk Type", "Confidence", "Risk Score", "Risk Level", "Matched Text", "Recommendation"})
	for index, asset := range data.FullCryptoInventory {
		setRow(file, sheet, index+2, []any{
			asset.ApplicationName,
			asset.SourceType,
			asset.FilePath,
			asset.LineNumber,
			asset.CryptoUsageType,
			asset.Algorithm,
			asset.Severity,
			asset.RiskType,
			asset.Confidence,
			asset.QuantumRiskScore,
			asset.RiskLevel,
			asset.MatchedText,
			asset.Recommendation,
		})
	}
	_ = file.SetColWidth(sheet, "A", "M", 20)
}

func setRow(file *excelize.File, sheet string, row int, values []any) {
	for index, value := range values {
		cell, err := excelize.CoordinatesToCellName(index+1, row)
		if err != nil {
			continue
		}
		_ = file.SetCellValue(sheet, cell, value)
	}
}

func topAssets(assets []inventory.CryptoAsset, limit int) []CriticalAsset {
	copyAssets := append([]inventory.CryptoAsset(nil), assets...)
	sort.SliceStable(copyAssets, func(i, j int) bool {
		left := copyAssets[i]
		right := copyAssets[j]
		if left.QuantumRiskScore != right.QuantumRiskScore {
			return left.QuantumRiskScore > right.QuantumRiskScore
		}
		if severityRank(left.Severity) != severityRank(right.Severity) {
			return severityRank(left.Severity) > severityRank(right.Severity)
		}
		if confidenceRank(left.Confidence) != confidenceRank(right.Confidence) {
			return confidenceRank(left.Confidence) > confidenceRank(right.Confidence)
		}
		if usageRank(left.CryptoUsageType) != usageRank(right.CryptoUsageType) {
			return usageRank(left.CryptoUsageType) > usageRank(right.CryptoUsageType)
		}
		if left.FilePath != right.FilePath {
			return left.FilePath < right.FilePath
		}
		return left.LineNumber < right.LineNumber
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

func severityRank(value string) int {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "critical":
		return 4
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 0
	}
}

func confidenceRank(value string) int {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 0
	}
}

func usageRank(value string) int {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "private key material":
		return 6
	case "private key import", "key file reference", "public key import":
		return 5
	case "jwt signing", "digital signature":
		return 4
	case "tls/ssl configuration":
		return 3
	case "key generation", "encryption/decryption":
		return 2
	case "jwt token handling", "crypto library import":
		return 1
	default:
		return 0
	}
}

func buildExecutiveSummary(data JSONReport) ExecutiveSummary {
	return ExecutiveSummary{
		OverallRiskPosture:  data.RiskPosture,
		TopCriticalAssets:   data.TopCriticalAssets,
		FindingsByAlgorithm: data.FindingsByAlgorithm,
		FindingsByUsageType: data.FindingsByUsageType,
		Actions30Days: []string{
			"Validate high and critical findings with application owners.",
			"Confirm TLS certificate algorithms, key sizes, and expiry dates for internet-facing domains.",
			"Identify owners and systems of record for each crypto asset.",
		},
		Actions60Days: []string{
			"Prioritize remediation for externally exposed signatures, key exchange, certificates, and private key material.",
			"Define target hybrid or post-quantum patterns for affected platforms.",
			"Request vendor PQC roadmaps for third-party dependencies.",
		},
		Actions90Days: []string{
			"Create migration backlog items with owners, test plans, and rollback plans.",
			"Pilot crypto-agility changes for the highest-risk services.",
			"Establish recurring scans in CI or release readiness checks.",
		},
	}
}

func sortedKeys(values map[string]int) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
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
