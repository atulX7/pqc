package rules

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Rule struct {
	ID         string `yaml:"id"`
	Name       string `yaml:"name"`
	Pattern    string `yaml:"pattern"`
	Algorithm  string `yaml:"algorithm"`
	RiskType   string `yaml:"risk_type"`
	Severity   string `yaml:"severity"`
	Confidence string `yaml:"confidence,omitempty"`
	Generic    bool   `yaml:"generic,omitempty"`
	Priority   int    `yaml:"priority,omitempty"`
}

type ruleFile struct {
	Rules []Rule `yaml:"rules"`
}

func LoadRules(path string) ([]Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("open rules file: %w", err)
	}

	var parsed ruleFile
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return nil, fmt.Errorf("parse rules file: %w", err)
	}
	if len(parsed.Rules) == 0 {
		return nil, fmt.Errorf("no rules loaded from %s", path)
	}

	for index := range parsed.Rules {
		if err := normalizeRule(&parsed.Rules[index]); err != nil {
			return nil, fmt.Errorf("invalid rule %d: %w", index+1, err)
		}
	}
	return parsed.Rules, nil
}

func normalizeRule(rule *Rule) error {
	rule.ID = strings.TrimSpace(rule.ID)
	rule.Name = strings.TrimSpace(rule.Name)
	rule.Pattern = strings.TrimSpace(rule.Pattern)
	rule.Algorithm = strings.TrimSpace(rule.Algorithm)
	rule.RiskType = strings.TrimSpace(rule.RiskType)
	rule.Severity = strings.TrimSpace(rule.Severity)
	rule.Confidence = strings.ToLower(strings.TrimSpace(rule.Confidence))

	if rule.ID == "" {
		return fmt.Errorf("id is required")
	}
	if rule.Pattern == "" {
		return fmt.Errorf("pattern is required for %s", rule.ID)
	}
	if rule.Algorithm == "" {
		return fmt.Errorf("algorithm is required for %s", rule.ID)
	}
	if rule.Confidence == "" {
		if rule.Generic {
			rule.Confidence = "low"
		} else {
			rule.Confidence = "high"
		}
	}
	if rule.Priority == 0 {
		if rule.Generic {
			rule.Priority = 10
		} else {
			rule.Priority = 50
		}
	}
	switch rule.Confidence {
	case "high", "medium", "low":
	default:
		return fmt.Errorf("confidence must be high, medium, or low for %s", rule.ID)
	}
	return nil
}
