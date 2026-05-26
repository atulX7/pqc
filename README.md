# PQC Readiness

PQC Readiness is an MVP assessment tool for post-quantum cryptography migration planning.

It performs read-only scans of source repositories and TLS domains, builds a crypto inventory, calculates a Quantum Risk Score, and produces JSON and Excel-ready reports with executive summary data and recommended 30/60/90-day actions.

Pipeline:

```text
Scan -> Detect -> Normalize -> Enrich -> Score -> Recommend -> Report
```

The scanner does not execute scanned code and does not automatically change cryptography.

## What This Helps Answer

- Where are RSA, ECC, ECDSA, ECDH, DH, DSA, weak hashes, JWT signing, TLS certificates, or private keys used?
- Which crypto assets are highest priority for post-quantum readiness review?
- Which systems may need hybrid or post-quantum migration planning?
- Which TLS certificates use RSA/ECDSA keys and when do they expire?
- What should application and security owners validate first?

This is an initial readiness assessment tool, not a formal compliance audit or a full SAST replacement. Findings should be reviewed with application and security owners before remediation decisions.

## Features

- Recursive source-code and configuration scanning
- Public GitHub repository scanning from the web UI
- Uploaded repository `.zip` scanning from the web UI
- TLS/domain certificate scanning from CLI and web UI
- Rule-based crypto detection from `rules/crypto_rules.yaml`
- Real YAML rule parsing with `gopkg.in/yaml.v3`
- Confidence on findings: `high`, `medium`, `low`
- Deduplication and rule priority to reduce over-reporting
- Private key material masking
- Structured crypto inventory records
- Business metadata inputs for sensitivity, exposure, criticality, secrecy lifetime, crypto agility, vendor dependency, and migration complexity
- Quantum Risk Score from 0-100
- Risk levels: Low, Moderate, High, Very High, Critical
- Executive summary with top assets and 30/60/90-day actions
- JSON report output
- Excel workbook output from the CLI
- Docker and Kubernetes deployment files

## Current Scan Sources

| Source | CLI | Web UI |
| --- | --- | --- |
| Local folder/repository | Yes | Sample only |
| Uploaded repository zip | No | Yes |
| Public GitHub repository URL | No | Yes |
| TLS domain | Yes | Yes |
| Multiple TLS domains | Yes | Yes |

Public GitHub scanning performs a shallow clone into a temporary directory, scans the checked-out files, and removes the clone after the request completes.

## Install

Install Go 1.22 or newer. This repo currently builds and tests with the installed Go toolchain in the local environment.

Build the CLI:

```bash
go build -o pqcscan ./cmd/pqcscan
```

Run tests:

```bash
go test ./...
```

## CLI Usage

Scan a local repository and export JSON plus Excel:

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
  --output report.json \
  --excel-output report.xlsx
```

Run directly without building:

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

Scan one TLS domain:

```bash
./pqcscan scan \
  --path ./sample_repo \
  --app-name "TLS Assessment" \
  --domain github.com \
  --output tls-report.json
```

Scan multiple TLS domains:

```bash
cat > domains.txt <<'EOF'
github.com
google.com
cloudflare.com
microsoft.com
amazon.com
openai.com
EOF

./pqcscan scan \
  --path ./sample_repo \
  --app-name "TLS Assessment" \
  --domains-file domains.txt \
  --output tls-report.json
```

## Web UI Usage

Run the web UI locally:

```bash
go run ./cmd/pqcweb
```

Open:

```text
http://localhost:8080
```

The web UI has four scan modes:

- `Sample`: scans the included sample repository
- `Zip`: scans an uploaded repository `.zip`
- `GitHub`: scans a public GitHub repository URL
- `Domains`: scans TLS certificates for one or more domains

The results page shows:

- Risk posture
- Highest risk score
- Total findings
- Files scanned
- Findings by algorithm
- Findings by usage type
- Top assets
- Inventory preview
- 30/60/90-day actions
- Full JSON report with a `Copy JSON` button

### Public GitHub Scan

Use a public repository URL such as:

```text
https://github.com/secana/fsharp-rsa-jwt.git
```

Supported URL forms include:

```text
https://github.com/owner/repo
https://github.com/owner/repo.git
https://github.com/owner/repo/tree/main
https://github.com/owner/repo/blob/main/file.ext
```

Only public `https://github.com/owner/repo` style repositories are supported in this MVP.

### TLS Domain Scan

In the `Domains` tab, enter one domain per line:

```text
github.com
google.com
cloudflare.com
microsoft.com
amazon.com
openai.com
```

Click `Scan Domains`.

Expected output includes `tls_certificate` inventory records with certificate subject, issuer, public key algorithm, key size, expiry, and days until expiry.

## Docker

Build the web container:

```bash
docker build -t pqc-readiness:dev .
```

Run locally:

```bash
docker run --rm -p 8080:8080 pqc-readiness:dev
```

Open:

```text
http://localhost:8080
```

## Kubernetes

Deploy to a local Kubernetes cluster such as Docker Desktop Kubernetes:

```bash
kubectl apply -k k8s
kubectl -n pqc get pods,svc
```

If the service is not reachable directly, start a port-forward:

```bash
kubectl -n pqc port-forward svc/pqc-readiness 8080:8080
```

Then open:

```text
http://localhost:8080
```

See [docs/kubernetes.md](docs/kubernetes.md) for endpoint and deployment details.

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

The MVP rule set includes patterns for:

- RSA, RSA-2048, RSA-4096
- ECC, ECDSA, ECDH
- Diffie-Hellman, DH
- DSA
- JWT algorithms such as RS256, RS384, RS512, ES256, ES384, ES512
- Java `KeyPairGenerator.getInstance("RSA")`
- Node `crypto.createSign("RSA-SHA256")`
- Python `cryptography` RSA/EC imports
- PyCryptodome RSA imports, key imports, OAEP encryption, and PKCS#1 signatures
- .NET/F# RSA APIs such as `RSA.Create`, `ImportRSAPrivateKey`, `ImportSubjectPublicKeyInfo`, `SigningCredentials`, and `RsaSecurityKey`
- TLS ECDHE/certificate usage
- PEM private key material
- SSH RSA/ECDSA keys
- SHA1 and MD5

Add or tune rules in [rules/crypto_rules.yaml](rules/crypto_rules.yaml).

## Supported File Types

The source scanner includes common application, config, and key file extensions, including:

```text
.go, .py, .js, .ts, .java, .cs, .fs, .fsx, .fsi, .fsproj,
.yaml, .yml, .json, .xml, .properties, .conf, .tf,
.pem, .key, .crt, .cert, .priv, .pub
```

Large files above 10 MB, binary files, generated folders, and dependency/vendor folders are skipped.

## Report Contents

The JSON report includes:

- `application_name`
- `executive_summary`
- `scan_summary`
- `total_findings`
- `findings_by_algorithm`
- `findings_by_severity`
- `findings_by_usage_type`
- `top_critical_assets`
- `full_crypto_inventory`
- `highest_quantum_risk_score`
- `risk_posture`
- `recommended_next_steps`

Inventory records include:

- source type
- file path or domain
- line number
- usage type
- algorithm
- severity
- risk type
- confidence
- matched text
- business metadata
- risk score
- recommendation

The CLI can also export an Excel workbook with:

- Executive summary
- Top assets
- Full crypto inventory

## Security Notes

- The scanner is read-only.
- It never executes scanned source code.
- Public GitHub scans perform shallow clones into temporary directories.
- Uploaded zip files are extracted into temporary directories.
- Temporary scan directories are removed after each web request.
- TLS domain scanning only performs a TLS handshake and reads certificate metadata.
- Private key material is masked in reports.
- `matched_text` is capped at 300 characters.
- Common generated/vendor directories are skipped, including `.git`, `node_modules`, `venv`, `__pycache__`, `dist`, `build`, `.next`, and `target`.

## License

This project is source-available for learning, research, and non-commercial evaluation.

Commercial use requires prior written permission from the copyright holder. Commercial use includes paid consulting, customer assessments, SaaS or managed-service use, redistribution, sublicensing, or inclusion in a commercial product or platform.

See [LICENSE](LICENSE) and [NOTICE](NOTICE) for full terms.

This is not an OSI open-source license. If you need broad open-source rights, commercial rights, or redistribution rights, contact the repository owner for a separate license.

## First-Customer Positioning

This tool is best positioned as an initial PQC readiness inventory and migration planning accelerator.

Suggested wording:

```text
This MVP performs a read-only PQC readiness assessment across source repositories and TLS domains. It identifies classical public-key cryptography, weak hashes, JWT signing, certificate algorithms, and exposed key material, then produces a prioritized crypto inventory for validation by application and security owners.
```

Good first pilot scope:

- 3-5 repositories
- 5-20 TLS domains
- business metadata for each application
- JSON/Excel report review
- prioritized 30/60/90-day migration planning discussion

Do not position this MVP as a formal audit sign-off, complete SAST replacement, or guarantee of zero false positives/false negatives.

## Tests

```bash
go test ./...
```
