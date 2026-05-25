package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/atulX7/pqc/internal/app"
	"github.com/atulX7/pqc/internal/report"
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
	options := app.DefaultScanOptions()
	flags.StringVar(&options.Path, "path", options.Path, "local repository or folder to scan")
	flags.StringVar(&options.ApplicationName, "app-name", options.ApplicationName, "application name")
	flags.StringVar(&options.Sensitivity, "sensitivity", options.Sensitivity, "comma-separated sensitivity flags")
	flags.StringVar(&options.Exposure, "exposure", options.Exposure, "exposure level")
	flags.StringVar(&options.BusinessCriticality, "business-criticality", options.BusinessCriticality, "business criticality")
	flags.IntVar(&options.SecrecyLifetimeYears, "secrecy-lifetime-years", options.SecrecyLifetimeYears, "secrecy lifetime in years")
	flags.StringVar(&options.CryptoAgility, "crypto-agility", options.CryptoAgility, "crypto agility level")
	flags.StringVar(&options.VendorDependency, "vendor-dependency", options.VendorDependency, "vendor dependency level")
	flags.StringVar(&options.MigrationComplexity, "migration-complexity", options.MigrationComplexity, "migration complexity")
	flags.StringVar(&options.RulesPath, "rules", options.RulesPath, "rules YAML path")
	output := flags.String("output", "report.json", "JSON output path")

	if err := flags.Parse(args); err != nil {
		return err
	}

	jsonReport, err := app.RunScan(options)
	if err != nil {
		return err
	}

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
