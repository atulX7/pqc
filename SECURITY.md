# Security Policy

## Supported Versions

This project is an MVP. Security fixes are handled on the active `dev` branch
until a stable release process is introduced.

## Reporting a Security Issue

If you find a security issue, please do not open a public issue with exploit
details or sensitive data.

Contact the repository owner with:

- a short description of the issue;
- affected files or endpoints;
- steps to reproduce, if safe to share;
- impact and suggested mitigation, if known.

## Scanner Safety Model

- The scanner is read-only.
- It does not execute scanned source code.
- Public GitHub scans perform shallow clones into temporary directories.
- Uploaded zip files are extracted into temporary directories.
- Temporary scan directories are removed after each web request.
- TLS scans perform a TLS handshake and read certificate metadata.
- Private key material is masked in reports.

## Responsible Use

Use this tool only on repositories, archives, and domains that you own or are
authorized to assess.
