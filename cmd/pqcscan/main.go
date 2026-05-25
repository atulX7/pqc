package main

import (
	"flag"
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

func main() {
	if len(os.Args) < 2 || os.Args[1] != "scan" {
		fmt.Fprintln(os.Stderr, "usage: pqcscan scan --path ./repo --app-name \"Clinical API\" --output report.json")
		os.Exit(2)
	}
	if err := runScan(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func runScan(args []string) error {
	flags := flag.NewFlagSet("scan", flag.ContinueOnError)
	scanPath := flags.String("path", ".", "local repository or folder to scan")
	appName := flags.String("app-name", "Unknown Application", "application name")
	sensitivity := flags.String("sensitivity", "internal", "comma-separated sensitivity flags")
	exposure := flags.String("exposure", "internal", "exposure level")
	businessCriticality := flags.String("business-criticality", "medium", "business criticality")
	secrecyLifetimeYears := flags.Int("secrecy-lifetime-years", 1, "secrecy lifetime in years")
	cryptoAgility := flags.String("crypto-agility", "partially_configurable", "crypto agility level")
	vendorDependency := flags.String("vendor-dependency", "none", "vendor dependency level")
	migrationComplexity := flags.String("migration-complexity", "code_change", "migration complexity")
	output := flags.String("output", "report.json", "JSON output path")
	rulesPath := flags.String("rules", "rules/crypto_rules.yaml", "rules YAML path")

	if err := flags.Parse(args); err != nil {
		return err
	}

	loadedRules, err := rules.LoadRules(*rulesPath)
	if err != nil {
		return err
	}

	absScanPath, err := filepath.Abs(*scanPath)
	if err != nil {
		return fmt.Errorf("resolve scan path: %w", err)
	}
	scanResult, err := scanner.ScanPath(absScanPath, loadedRules)
	if err != nil {
		return err
	}

	metadata := inventory.BusinessMetadata{
		ApplicationName:       *appName,
		SensitivityFlags:      splitCSV(*sensitivity),
		ExposureLevel:         *exposure,
		BusinessCriticality:   *businessCriticality,
		SecrecyLifetimeYears:  *secrecyLifetimeYears,
		CryptoAgilityLevel:    *cryptoAgility,
		VendorDependencyLevel: *vendorDependency,
		MigrationComplexity:   *migrationComplexity,
	}

	assets := inventory.NormalizeFindings(scanResult.Findings, metadata)
	assets = scoring.ScoreAssets(assets)
	assets = recommendation.AddRecommendations(assets)

	jsonReport := report.Build(*appName, report.ScanSummary{
		FilesScanned:    scanResult.FilesScanned,
		FilesSkipped:    scanResult.FilesSkipped,
		UnreadableFiles: scanResult.UnreadableFiles,
	}, assets)

	if err := report.WriteJSON(*output, jsonReport); err != nil {
		return fmt.Errorf("write JSON report: %w", err)
	}
	fmt.Printf("PQC scan complete: %d findings, highest risk %.2f (%s). Report written to %s\n",
		jsonReport.TotalFindings,
		jsonReport.HighestQuantumRiskScore,
		jsonReport.RiskPosture,
		*output,
	)
	return nil
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
