package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/atulX7/pqc/internal/inventory"
	"github.com/atulX7/pqc/internal/recommendation"
	"github.com/atulX7/pqc/internal/report"
	"github.com/atulX7/pqc/internal/rules"
	"github.com/atulX7/pqc/internal/scanner"
	"github.com/atulX7/pqc/internal/scoring"
)

type ScanOptions struct {
	Path                 string
	ApplicationName      string
	Sensitivity          string
	Exposure             string
	BusinessCriticality  string
	SecrecyLifetimeYears int
	CryptoAgility        string
	VendorDependency     string
	MigrationComplexity  string
	RulesPath            string
}

func DefaultScanOptions() ScanOptions {
	return ScanOptions{
		Path:                 ".",
		ApplicationName:      "Unknown Application",
		Sensitivity:          "internal",
		Exposure:             "internal",
		BusinessCriticality:  "medium",
		SecrecyLifetimeYears: 1,
		CryptoAgility:        "partially_configurable",
		VendorDependency:     "none",
		MigrationComplexity:  "code_change",
		RulesPath:            "rules/crypto_rules.yaml",
	}
}

func RunScan(options ScanOptions) (report.JSONReport, error) {
	if options.Path == "" {
		options.Path = "."
	}
	if options.ApplicationName == "" {
		options.ApplicationName = "Unknown Application"
	}
	if options.RulesPath == "" {
		options.RulesPath = "rules/crypto_rules.yaml"
	}

	loadedRules, err := rules.LoadRules(options.RulesPath)
	if err != nil {
		return report.JSONReport{}, err
	}

	absScanPath, err := filepath.Abs(options.Path)
	if err != nil {
		return report.JSONReport{}, fmt.Errorf("resolve scan path: %w", err)
	}
	scanResult, err := scanner.ScanPath(absScanPath, loadedRules)
	if err != nil {
		return report.JSONReport{}, err
	}

	metadata := inventory.BusinessMetadata{
		ApplicationName:       options.ApplicationName,
		SensitivityFlags:      splitCSV(options.Sensitivity),
		ExposureLevel:         options.Exposure,
		BusinessCriticality:   options.BusinessCriticality,
		SecrecyLifetimeYears:  options.SecrecyLifetimeYears,
		CryptoAgilityLevel:    options.CryptoAgility,
		VendorDependencyLevel: options.VendorDependency,
		MigrationComplexity:   options.MigrationComplexity,
	}

	assets := inventory.NormalizeFindings(scanResult.Findings, metadata)
	assets = scoring.ScoreAssets(assets)
	assets = recommendation.AddRecommendations(assets)

	return report.Build(options.ApplicationName, report.ScanSummary{
		FilesScanned:    scanResult.FilesScanned,
		FilesSkipped:    scanResult.FilesSkipped,
		UnreadableFiles: scanResult.UnreadableFiles,
	}, assets), nil
}

func splitCSV(value string) []string {
	var values []string
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			values = append(values, item)
		}
	}
	return values
}
