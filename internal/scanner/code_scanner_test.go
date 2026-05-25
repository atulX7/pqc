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

func TestScanPathDedupesGenericWhenSpecificRuleMatches(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "CryptoFactory.java"), []byte(`return KeyPairGenerator.getInstance("RSA");`), 0o644); err != nil {
		t.Fatal(err)
	}

	loadedRules := []rules.Rule{
		{ID: "rsa-generic", Name: "RSA usage", Pattern: "RSA", Algorithm: "RSA", RiskType: "quantum_vulnerable_public_key", Severity: "high", Confidence: "low", Generic: true, Priority: 10},
		{ID: "java-rsa-keypair", Name: "Java RSA KeyPairGenerator", Pattern: `KeyPairGenerator.getInstance("RSA")`, Algorithm: "RSA", RiskType: "quantum_vulnerable_public_key", Severity: "high", Confidence: "high", Priority: 80},
	}

	result, err := ScanPath(dir, loadedRules)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("expected 1 deduped finding, got %d", len(result.Findings))
	}
	if result.Findings[0].RuleID != "java-rsa-keypair" || result.Findings[0].Confidence != "high" {
		t.Fatalf("expected high-confidence specific finding, got %+v", result.Findings[0])
	}
}

func TestScanPathUsesSpecificPyCryptodomeRules(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "cert.py"), []byte("from Crypto.PublicKey import RSA\nprivateKey = RSA.import_key(open(\"rsa_private_key.pem\", \"rb\").read())"), 0o644); err != nil {
		t.Fatal(err)
	}

	loadedRules := []rules.Rule{
		{ID: "rsa-generic", Name: "RSA usage", Pattern: "RSA", Algorithm: "RSA", RiskType: "quantum_vulnerable_public_key", Severity: "high", Confidence: "low", Generic: true, Priority: 10},
		{ID: "pycryptodome-rsa-import", Name: "PyCryptodome RSA import", Pattern: "from Crypto.PublicKey import RSA", Algorithm: "RSA", RiskType: "quantum_vulnerable_public_key", Severity: "high", Confidence: "high", Priority: 40},
		{ID: "pycryptodome-rsa-import-key", Name: "PyCryptodome RSA key import", Pattern: "RSA.import_key", Algorithm: "RSA", RiskType: "quantum_vulnerable_public_key", Severity: "high", Confidence: "high", Priority: 80},
	}

	result, err := ScanPath(dir, loadedRules)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 2 {
		t.Fatalf("expected 2 specific findings, got %d", len(result.Findings))
	}
	for _, finding := range result.Findings {
		if finding.RuleID == "rsa-generic" {
			t.Fatalf("generic RSA finding should have been suppressed: %+v", finding)
		}
		if finding.Confidence != "high" {
			t.Fatalf("expected high confidence, got %+v", finding)
		}
	}
}

func TestScanPathPrefersConcretePrimitiveOverModuleImport(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "cert.py"), []byte("from Crypto.Cipher import PKCS1_OAEP"), 0o644); err != nil {
		t.Fatal(err)
	}

	loadedRules := []rules.Rule{
		{ID: "pycryptodome-pkcs1-oaep", Name: "PyCryptodome RSA OAEP encryption", Pattern: "PKCS1_OAEP", Algorithm: "RSA", RiskType: "quantum_vulnerable_public_key", Severity: "high", Confidence: "high", Priority: 90},
		{ID: "pycryptodome-cipher-module", Name: "PyCryptodome cipher module", Pattern: "Crypto.Cipher", Algorithm: "RSA", RiskType: "quantum_vulnerable_public_key", Severity: "medium", Confidence: "high", Priority: 30},
	}

	result, err := ScanPath(dir, loadedRules)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Findings) != 1 {
		t.Fatalf("expected 1 prioritized finding, got %d", len(result.Findings))
	}
	if result.Findings[0].RuleID != "pycryptodome-pkcs1-oaep" {
		t.Fatalf("expected OAEP finding to win, got %+v", result.Findings[0])
	}
}
