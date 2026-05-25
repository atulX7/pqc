package rules

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Rule struct {
	ID        string
	Name      string
	Pattern   string
	Algorithm string
	RiskType  string
	Severity  string
}

func LoadRules(path string) ([]Rule, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open rules file: %w", err)
	}
	defer file.Close()

	var loaded []Rule
	var current *Rule
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || line == "rules:" {
			continue
		}

		if strings.HasPrefix(line, "- ") {
			if current != nil {
				loaded = append(loaded, *current)
			}
			current = &Rule{}
			line = strings.TrimSpace(strings.TrimPrefix(line, "- "))
			if line == "" {
				continue
			}
		}

		if current == nil {
			continue
		}

		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		setRuleField(current, strings.TrimSpace(key), cleanYAMLValue(value))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read rules file: %w", err)
	}
	if current != nil {
		loaded = append(loaded, *current)
	}
	if len(loaded) == 0 {
		return nil, fmt.Errorf("no rules loaded from %s", path)
	}
	return loaded, nil
}

func setRuleField(rule *Rule, key string, value string) {
	switch key {
	case "id":
		rule.ID = value
	case "name":
		rule.Name = value
	case "pattern":
		rule.Pattern = value
	case "algorithm":
		rule.Algorithm = value
	case "risk_type":
		rule.RiskType = value
	case "severity":
		rule.Severity = value
	}
}

func cleanYAMLValue(raw string) string {
	value := strings.TrimSpace(raw)
	value = strings.Trim(value, "\"'")
	return value
}
