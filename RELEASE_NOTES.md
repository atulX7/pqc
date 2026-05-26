# Release Notes

## v0.1

Initial source-available MVP release of PQC Readiness.

### Highlights

- CLI scanner for post-quantum cryptography readiness inventory
- Web UI for sample scans, repository zip scans, public GitHub repository scans, and TLS domain scans
- Docker and Kubernetes deployment support
- JSON report output
- Excel workbook export from the CLI
- Executive summary with risk posture, top assets, findings by algorithm, findings by usage type, and recommended 30/60/90-day actions

### Scanner Capabilities

- Detects classical public-key crypto and related usage patterns, including RSA, RSA-2048, RSA-4096, ECC, ECDSA, ECDH, DH, DSA, MD5, SHA1, JWT signing algorithms, TLS certificate keys, and private key material
- Supports common source, configuration, certificate, and key file extensions
- Includes F#/.NET rules for `RSA.Create`, `ImportRSAPrivateKey`, `ImportSubjectPublicKeyInfo`, `SigningCredentials`, `RsaSecurityKey`, `JwtSecurityToken`, and `JwtSecurityTokenHandler`
- Includes Python/PyCryptodome and Python `cryptography` detection patterns
- Masks private key material in reports
- Adds confidence levels to findings
- Uses rule priority and deduplication to reduce over-reporting

### Web UI

- Tabbed assessment workflow: Sample, Zip, GitHub, Domains
- Business context inputs for sensitivity, exposure, criticality, secrecy lifetime, crypto agility, vendor dependency, and migration complexity
- Report dashboard with metrics, finding distributions, top assets, inventory preview, actions, and full JSON
- Copy JSON button for generated reports

### TLS Domain Scanning

- CLI supports `--domain` and `--domains-file`
- Web UI supports one or more TLS domains in the Domains tab
- TLS findings include certificate subject, issuer, public key algorithm, key size, expiry, and days until expiry

### Deployment

- Single-container web deployment
- Local Docker support
- Kubernetes manifests under `k8s/`
- Runtime image includes `git` for public GitHub repository scans

### License

This release is source-available for learning, research, and non-commercial evaluation. Commercial use requires prior written permission from the copyright holder. See `LICENSE` and `NOTICE`.

### Known Limitations

- This is an MVP readiness tool, not a formal audit or compliance certification
- Findings may include false positives or false negatives
- Private GitHub repository authentication is not implemented
- CI/SARIF integrations are not implemented
- Risk scoring is metadata-driven and should be validated by application and security owners
- Results should be reviewed before remediation or migration decisions

### Validation

Validated with:

```bash
go test ./...
```
