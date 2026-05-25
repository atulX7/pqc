package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/atulX7/pqc/internal/rules"
)

func TestScanPathFindsCryptoAndMasksPrivateKeys(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "service.py"), []byte(`JWT_ALGORITHM = "RS256"`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "secret.key"), []byte("-----BEGIN RSA PRIVATE KEY-----\nabc\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "node_modules"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "node_modules", "ignored.js"), []byte("RS256"), 0o644); err != nil {
		t.Fatal(err)
	}

	loadedRules := []rules.Rule{
		{ID: "jwt-rs256", Name: "JWT RS256 signing", Pattern: "RS256", Algorithm: "RSA", RiskType: "quantum_vulnerable_signature", Severity: "high"},
		{ID: "rsa-private-key", Name: "RSA private key committed", Pattern: "-----BEGIN RSA PRIVATE KEY-----", Algorithm: "RSA", RiskType: "secret_exposure", Severity: "critical"},
	}

	result, err := ScanPath(dir, loadedRules)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(result.Findings))
	}

	var masked bool
	for _, finding := range result.Findings {
		if finding.RuleID == "rsa-private-key" && finding.MatchedText == "[MASKED PRIVATE KEY MATERIAL]" {
			masked = true
		}
		if finding.FilePath == "node_modules/ignored.js" {
			t.Fatal("scanner should skip node_modules")
		}
	}
	if !masked {
		t.Fatal("expected private key material to be masked")
	}
}

func TestScanPathDoesNotMatchCryptoPatternInsideIdentifier(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "EmotiSimUnity.cs"), []byte(`var resp = JsonUtility.FromJson<EmotionResponse>(www.downloadHandler.text);`), 0o644); err != nil {
		t.Fatal(err)
	}

	loadedRules := []rules.Rule{
		{ID: "dh-generic", Name: "DH usage", Pattern: "DH", Algorithm: "DH", RiskType: "quantum_vulnerable_key_exchange", Severity: "high"},
	}

	result, err := ScanPath(dir, loadedRules)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(result.Findings))
	}
}

func TestScanPathMatchesStandaloneCryptoPattern(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "tls.yaml"), []byte(`key_exchange: DH`), 0o644); err != nil {
		t.Fatal(err)
	}

	loadedRules := []rules.Rule{
		{ID: "dh-generic", Name: "DH usage", Pattern: "DH", Algorithm: "DH", RiskType: "quantum_vulnerable_key_exchange", Severity: "high"},
	}

	result, err := ScanPath(dir, loadedRules)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(result.Findings))
	}
}
