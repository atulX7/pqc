package gitrepo

import (
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

const cloneTimeout = 60 * time.Second

func ValidatePublicGitHubURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", fmt.Errorf("github repo URL is required")
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse github repo URL: %w", err)
	}
	if parsed.Scheme != "https" {
		return "", fmt.Errorf("only https://github.com URLs are supported")
	}
	if !strings.EqualFold(parsed.Hostname(), "github.com") {
		return "", fmt.Errorf("only github.com repositories are supported")
	}
	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("credentials, query parameters, and fragments are not allowed")
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("expected github URL format https://github.com/owner/repo")
	}
	owner := strings.TrimSpace(parts[0])
	repo := strings.TrimSuffix(strings.TrimSpace(parts[1]), ".git")
	if !safeGitHubPathPart(owner) || !safeGitHubPathPart(repo) {
		return "", fmt.Errorf("expected github URL format https://github.com/owner/repo using a valid GitHub owner and repo name")
	}
	normalized := &url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   "/" + owner + "/" + repo + ".git",
	}
	return normalized.String(), nil
}

func ClonePublicRepo(ctx context.Context, repoURL string, branch string, destination string) error {
	normalizedURL, err := ValidatePublicGitHubURL(repoURL)
	if err != nil {
		return err
	}
	args := []string{"clone", "--depth=1", "--single-branch"}
	branch = strings.TrimSpace(branch)
	if branch != "" {
		if !safeBranch(branch) {
			return fmt.Errorf("branch contains unsupported characters")
		}
		args = append(args, "--branch", branch)
	}
	args = append(args, normalizedURL, destination)

	cloneCtx, cancel := context.WithTimeout(ctx, cloneTimeout)
	defer cancel()
	cmd := exec.CommandContext(cloneCtx, "git", args...)
	output, err := cmd.CombinedOutput()
	if cloneCtx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("git clone timed out")
	}
	if err != nil {
		return fmt.Errorf("git clone failed: %s", sanitizeGitOutput(string(output)))
	}
	return nil
}

func safeGitHubPathPart(value string) bool {
	if value == "" || value == "." || value == ".." {
		return false
	}
	for _, char := range value {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '.' || char == '_' || char == '-' {
			continue
		}
		return false
	}
	return true
}

func safeBranch(value string) bool {
	if value == "" || strings.HasPrefix(value, "-") || strings.Contains(value, "..") || strings.ContainsAny(value, " \t\n\r~^:?*[\\") {
		return false
	}
	return true
}

func sanitizeGitOutput(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "unknown git error"
	}
	if len(value) > 500 {
		return value[:500]
	}
	return value
}
