package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/atulX7/pqc/internal/inventory"
	"github.com/atulX7/pqc/internal/rules"
)

const maxFileSizeBytes = 10 * 1024 * 1024

var skippedDirs = map[string]bool{
	".git": true, "node_modules": true, "venv": true, "__pycache__": true,
	"dist": true, "build": true, ".next": true, "target": true,
}

var scannedExtensions = map[string]bool{
	".go": true, ".py": true, ".js": true, ".ts": true, ".java": true, ".cs": true,
	".rb": true, ".php": true, ".yaml": true, ".yml": true, ".json": true,
	".xml": true, ".properties": true, ".env": true, ".pem": true, ".key": true,
	".crt": true, ".conf": true, ".tf": true,
}

type ScanResult struct {
	Findings        []inventory.Finding
	FilesScanned    int
	FilesSkipped    int
	UnreadableFiles []string
}

func ScanPath(root string, loadedRules []rules.Rule) (ScanResult, error) {
	var result ScanResult
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			result.UnreadableFiles = append(result.UnreadableFiles, path)
			return nil
		}
		if entry.IsDir() {
			if skippedDirs[entry.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if !shouldScan(path) {
			result.FilesSkipped++
			return nil
		}
		findings, err := scanFile(root, path, loadedRules)
		if err != nil {
			result.FilesSkipped++
			result.UnreadableFiles = append(result.UnreadableFiles, path)
			return nil
		}
		result.FilesScanned++
		result.Findings = append(result.Findings, findings...)
		return nil
	})
	if err != nil {
		return result, fmt.Errorf("walk scan path: %w", err)
	}
	return result, nil
}

func shouldScan(path string) bool {
	name := filepath.Base(path)
	if name == ".env" {
		return true
	}
	return scannedExtensions[strings.ToLower(filepath.Ext(path))]
}

func scanFile(root string, path string, loadedRules []rules.Rule) ([]inventory.Finding, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.Size() > maxFileSizeBytes {
		return nil, fmt.Errorf("file too large")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	head := make([]byte, 512)
	n, err := file.Read(head)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if bytes.IndexByte(head[:n], 0) >= 0 || !utf8.Valid(head[:n]) {
		return nil, fmt.Errorf("binary file")
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	relPath, err := filepath.Rel(root, path)
	if err != nil {
		relPath = path
	}

	var findings []inventory.Finding
	lineScanner := bufio.NewScanner(file)
	lineScanner.Buffer(make([]byte, 1024), 1024*1024)
	lineNumber := 0
	for lineScanner.Scan() {
		lineNumber++
		line := lineScanner.Text()
		var lineFindings []inventory.Finding
		for _, rule := range loadedRules {
			if matchesRulePattern(line, rule.Pattern) {
				lineFindings = append(lineFindings, inventory.Finding{
					SourceType:  sourceTypeFor(path),
					FilePath:    filepath.ToSlash(relPath),
					LineNumber:  lineNumber,
					MatchedText: sanitizeMatch(line),
					RuleID:      rule.ID,
					RuleName:    rule.Name,
					Algorithm:   rule.Algorithm,
					Severity:    rule.Severity,
					RiskType:    rule.RiskType,
					Confidence:  confidenceForRule(rule, line),
					Priority:    rule.Priority,
				})
			}
		}
		findings = append(findings, dedupeLineFindings(lineFindings)...)
	}
	if err := lineScanner.Err(); err != nil {
		return nil, err
	}
	return findings, nil
}

func matchesRulePattern(line string, pattern string) bool {
	lineLower := strings.ToLower(line)
	patternLower := strings.ToLower(pattern)
	if pattern == "" {
		return false
	}
	if isExactPattern(pattern) {
		return strings.Contains(lineLower, patternLower)
	}

	start := 0
	for {
		index := strings.Index(lineLower[start:], patternLower)
		if index < 0 {
			return false
		}
		index += start
		beforeOK := index == 0 || !isTokenChar(rune(lineLower[index-1]))
		after := index + len(patternLower)
		afterOK := after == len(lineLower) || !isTokenChar(rune(lineLower[after]))
		if beforeOK && afterOK {
			return true
		}
		start = index + 1
	}
}

func isExactPattern(pattern string) bool {
	return strings.ContainsAny(pattern, "\"'().:_/- ") || len(pattern) > 8
}

func isTokenChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '_' ||
		char == '-'
}

func confidenceForRule(rule rules.Rule, line string) string {
	if !rule.Generic {
		return "high"
	}
	value := strings.ToLower(line)
	switch strings.ToUpper(rule.Algorithm) {
	case "RSA":
		if strings.Contains(value, "encrypt") || strings.Contains(value, "sign") || strings.Contains(value, "key") || strings.Contains(value, "cert") || strings.Contains(value, "tls") {
			return "medium"
		}
	case "DH", "ECDH":
		if strings.Contains(value, "key") || strings.Contains(value, "exchange") || strings.Contains(value, "tls") || strings.Contains(value, "cipher") {
			return "medium"
		}
	case "MD5", "SHA1":
		if strings.Contains(value, "hash") || strings.Contains(value, "digest") || strings.Contains(value, "checksum") {
			return "medium"
		}
	}
	if rule.Confidence != "" {
		return rule.Confidence
	}
	return "low"
}

func dedupeLineFindings(findings []inventory.Finding) []inventory.Finding {
	if len(findings) <= 1 {
		return findings
	}
	maxPriority := 0
	for _, finding := range findings {
		if finding.Priority > maxPriority {
			maxPriority = finding.Priority
		}
	}

	seen := map[string]bool{}
	var deduped []inventory.Finding
	for _, finding := range findings {
		if finding.Priority < maxPriority {
			continue
		}
		key := strings.Join([]string{finding.FilePath, fmt.Sprint(finding.LineNumber), finding.RuleID, finding.Algorithm, finding.RiskType}, "|")
		if seen[key] {
			continue
		}
		seen[key] = true
		deduped = append(deduped, finding)
	}
	return deduped
}

func sourceTypeFor(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".pem", ".key", ".crt":
		return "certificate_or_key"
	case ".yaml", ".yml", ".json", ".xml", ".properties", ".env", ".conf", ".tf":
		return "configuration"
	default:
		return "source_code"
	}
}

func sanitizeMatch(value string) string {
	text := strings.TrimSpace(value)
	upper := strings.ToUpper(text)
	if strings.Contains(upper, "PRIVATE KEY") {
		return "[MASKED PRIVATE KEY MATERIAL]"
	}
	if len(text) > 300 {
		return text[:300]
	}
	return text
}
