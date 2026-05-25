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
