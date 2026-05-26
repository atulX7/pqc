package gitrepo

import "testing"

func TestValidatePublicGitHubURL(t *testing.T) {
	tests := map[string]string{
		"https://github.com/owner/repo":                 "https://github.com/owner/repo.git",
		"https://github.com/owner/repo/":                "https://github.com/owner/repo.git",
		"https://github.com/owner/repo.git":             "https://github.com/owner/repo.git",
		"https://github.com/owner/repo/tree/main":       "https://github.com/owner/repo.git",
		"https://github.com/owner/repo/blob/main/a.txt": "https://github.com/owner/repo.git",
	}
	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			got, err := ValidatePublicGitHubURL(input)
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Fatalf("expected %q, got %q", want, got)
			}
		})
	}
}

func TestValidatePublicGitHubURLRejectsUnsafeInputs(t *testing.T) {
	inputs := []string{
		"git@github.com:owner/repo.git",
		"http://github.com/owner/repo",
		"https://github.com/owner/repo?token=secret",
		"https://example.com/owner/repo",
		"https://github.com/owner",
		"https://github.com/owner/repo with space",
	}
	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			if _, err := ValidatePublicGitHubURL(input); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
