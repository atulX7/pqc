package app

import (
	"fmt"
	"os"
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
	Domains              []string
	DomainsFile          string
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
	findings := append([]inventory.Finding(nil), scanResult.Findings...)
	domains, err := domainsFromOptions(options)
	if err != nil {
		return report.JSONReport{}, err
	}
	if len(domains) > 0 {
		tlsResult := scanner.ScanTLSDomains(domains)
		findings = append(findings, tlsResult.Findings...)
		scanResult.UnreadableFiles = append(scanResult.UnreadableFiles, tlsResult.Errors...)
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

	assets := inventory.NormalizeFindings(findings, metadata)
	assets = scoring.ScoreAssets(assets)
	assets = recommendation.AddRecommendations(assets)

	return report.Build(options.ApplicationName, report.ScanSummary{
		FilesScanned:    scanResult.FilesScanned,
		FilesSkipped:    scanResult.FilesSkipped,
		UnreadableFiles: scanResult.UnreadableFiles,
	}, assets), nil
}

func domainsFromOptions(options ScanOptions) ([]string, error) {
	domains := append([]string(nil), options.Domains...)
	if strings.TrimSpace(options.DomainsFile) == "" {
		return normalizeDomains(domains), nil
	}
	data, err := os.ReadFile(options.DomainsFile)
	if err != nil {
		return nil, fmt.Errorf("read domains file: %w", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		domains = append(domains, line)
	}
	return normalizeDomains(domains), nil
}

func normalizeDomains(domains []string) []string {
	seen := map[string]bool{}
	normalized := make([]string, 0, len(domains))
	for _, domain := range domains {
		domain = strings.TrimSpace(domain)
		if domain == "" || seen[domain] {
			continue
		}
		seen[domain] = true
		normalized = append(normalized, domain)
	}
	return normalized
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
