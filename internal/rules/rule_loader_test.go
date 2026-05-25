package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRules(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rules.yaml")
	content := []byte(`rules:
  - id: jwt-rs256
    name: JWT RS256 signing
    pattern: RS256
    algorithm: RSA
    risk_type: quantum_vulnerable_signature
    severity: high
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadRules(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(loaded))
	}
	if loaded[0].ID != "jwt-rs256" || loaded[0].Pattern != "RS256" {
		t.Fatalf("unexpected rule: %+v", loaded[0])
	}
}

func TestLoadRulesPreservesQuotedPatternsAndConfidence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rules.yaml")
	content := []byte(`rules:
  - id: java-rsa-keypair
    name: Java RSA KeyPairGenerator
    pattern: 'KeyPairGenerator.getInstance("RSA")'
    algorithm: RSA
    risk_type: quantum_vulnerable_public_key
    severity: high
  - id: node-rsa-sign
    name: Node RSA SHA256 signing
    pattern: 'crypto.createSign("RSA-SHA256")'
    algorithm: RSA
    risk_type: quantum_vulnerable_signature
    severity: high
  - id: dh-generic
    name: DH usage
    pattern: DH
    algorithm: DH
    risk_type: quantum_vulnerable_key_exchange
    severity: high
    confidence: low
    generic: true
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadRules(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded[0].Pattern != `KeyPairGenerator.getInstance("RSA")` {
		t.Fatalf("quoted Java pattern was not preserved: %q", loaded[0].Pattern)
	}
	if loaded[1].Pattern != `crypto.createSign("RSA-SHA256")` {
		t.Fatalf("quoted Node pattern was not preserved: %q", loaded[1].Pattern)
	}
	if loaded[0].Confidence != "high" {
		t.Fatalf("expected exact rule high confidence, got %q", loaded[0].Confidence)
	}
	if !loaded[2].Generic || loaded[2].Confidence != "low" {
		t.Fatalf("expected generic low-confidence rule, got %+v", loaded[2])
	}
}
