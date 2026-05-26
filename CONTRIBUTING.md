# Contributing

Thank you for your interest in improving PQC Readiness.

## Contribution Scope

Useful contributions include:

- crypto detection rules;
- scanner false-positive or false-negative fixes;
- report quality improvements;
- tests;
- documentation;
- deployment examples.

## Development

Run tests before submitting changes:

```bash
go test ./...
```

Format Go code before committing:

```bash
gofmt -w .
```

## Pull Request Checklist

- The scanner remains read-only.
- Scanned code is never executed.
- Private key material remains masked.
- New rules include focused tests where practical.
- Documentation is updated for user-facing behavior changes.

## Licensing of Contributions

By contributing, you agree that your contribution may be included in this
project under the repository's license terms. Do not contribute code that you do
not have the right to submit.

If a separate contributor agreement is introduced later, future contributions
may require that additional agreement before acceptance.
