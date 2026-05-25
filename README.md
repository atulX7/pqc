# PQC Readiness

PQC Readiness is an MVP command-line scanner for post-quantum cryptography migration planning.

It scans source code and configuration files, detects cryptographic usage, normalizes findings into a crypto inventory, calculates a 0-100 Quantum Risk Score, generates recommendations, and exports a JSON report.

Pipeline:

```text
Scan -> Detect -> Normalize -> Enrich -> Score -> Recommend -> Report
```

The MVP is read-only. It does not execute scanned code and does not automatically replace cryptography.

## Features

- Recursive local repository/folder scanning
- Rule-based crypto detection from `rules/crypto_rules.yaml`
- Structured crypto inventory records
- Business metadata inputs for sensitivity, exposure, criticality, secrecy lifetime, agility, vendor dependency, and migration complexity
- Quantum Risk Score from 0-100
- Risk levels: Low, Moderate, High, Very High, Critical
- Rule-based remediation recommendations
- JSON report export

## Install

Install Go 1.22 or newer, then build the CLI:

```bash
go build -o pqcscan ./cmd/pqcscan
```

## Usage

```bash
./pqcscan scan \
  --path ./sample_repo \
  --app-name "Clinical API" \
  --sensitivity PHI,GxP \
  --exposure internet_authenticated \
  --business-criticality high \
  --secrecy-lifetime-years 15 \
  --crypto-agility partially_configurable \
  --vendor-dependency cloud_with_roadmap \
  --migration-complexity multi_system \
  --output report.json
```

You can also run directly without building:

```bash
go run ./cmd/pqcscan scan \
  --path ./sample_repo \
  --app-name "Clinical API" \
  --sensitivity PHI,GxP \
  --exposure internet_authenticated \
  --business-criticality high \
  --secrecy-lifetime-years 15 \
  --crypto-agility partially_configurable \
  --vendor-dependency cloud_with_roadmap \
  --migration-complexity multi_system \
  --output report.json
```

## Metadata Values

Sensitivity examples:

```text
public, internal, confidential, pii, spi, phi, gxp, pharma_ip, financial, legal
```

Exposure:

```text
offline, internal, partner, internet_authenticated, internet_public
```

Business criticality:

```text
low, medium, high, mission_critical
```

Crypto agility:

```text
centralized_policy_driven, configurable, partially_configurable, hardcoded, vendor_controlled_or_legacy
```

Vendor dependency:

```text
none, cloud_with_roadmap, saas_unknown_roadmap, legacy_no_roadmap
```

Migration complexity:

```text
config_change, library_upgrade, code_change, multi_system, legacy_regulated_vendor_bound
```

## Detected Patterns

The MVP rule set includes RSA, RSA-2048, RSA-4096, ECC, ECDSA, ECDH, Diffie-Hellman, DH, DSA, JWT algorithms such as RS256 and ES256, SHA1, MD5, TLS ECDHE cipher suites, PEM private keys, SSH RSA/ECDSA keys, Java `KeyPairGenerator` usage, Python `cryptography` RSA/EC imports, and Node `crypto.createSign("RSA-SHA256")`.

Add or tune rules in [rules/crypto_rules.yaml](rules/crypto_rules.yaml).

## Report

The JSON report includes:

- `application_name`
- `scan_summary`
- `total_findings`
- `findings_by_algorithm`
- `findings_by_severity`
- `top_critical_assets`
- `full_crypto_inventory`
- `highest_quantum_risk_score`
- `risk_posture`
- `recommended_next_steps`

## Security Notes

- The scanner only reads files.
- It never executes scanned code.
- Binary files and large files above 10 MB are skipped.
- Private key material is masked in reports.
- `matched_text` is capped at 300 characters.
- Common generated/vendor directories are skipped, including `.git`, `node_modules`, `venv`, `__pycache__`, `dist`, `build`, `.next`, and `target`.

## Tests

```bash
go test ./...
```
