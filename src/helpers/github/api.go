package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models/github"
)

var (
	client = &http.Client{
		Timeout: 30 * time.Second,
	}

	// Track which auth method was used (for logging)
	authMethodUsed string
)

// GetGitHubToken retrieves the GitHub token from config or gh CLI automatically
// Priority:
// 1. GITHUB_TOKEN environment variable (via helpers.AppConfig)
// 2. gh CLI credentials (automatic fallback)
func GetGitHubToken() (string, error) {
	// First, try to get token from config (GITHUB_TOKEN env var)
	if helpers.AppConfig.GITHUB_TOKEN != "" {
		if authMethodUsed == "" {
			authMethodUsed = "GITHUB_TOKEN environment variable"
			if helpers.AppConfig.SINALOA_DEBUG {
				fmt.Printf("[DEBUG] Using GitHub authentication from: %s\n", authMethodUsed)
			}
		}
		return helpers.AppConfig.GITHUB_TOKEN, nil
	}

	// Automatic fallback to gh CLI (for local development)
	// This allows seamless usage without manual configuration
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("GitHub authentication failed. Please either:\n  1. Set GITHUB_TOKEN environment variable, or\n  2. Authenticate with 'gh auth login'\nError: %w", err)
	}

	token := strings.TrimSpace(string(output))
	if token == "" {
		return "", fmt.Errorf("GitHub token is empty. Please authenticate with 'gh auth login' or set GITHUB_TOKEN env var")
	}

	if authMethodUsed == "" {
		authMethodUsed = "gh CLI (automatic fallback)"
		if helpers.AppConfig.SINALOA_DEBUG {
			fmt.Printf("[DEBUG] Using GitHub authentication from: %s\n", authMethodUsed)
		}
	}

	return token, nil
}

// GetAuthMethod returns the authentication method that was used
func GetAuthMethod() string {
	if authMethodUsed == "" {
		return "not authenticated yet"
	}
	return authMethodUsed
}

// GitHubAPICall makes a generic GitHub API call
func GitHubAPICall(endpoint string, result interface{}) error {
	token, err := GetGitHubToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.github.com%s", endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// ListRepositories fetches repositories from an organization
func ListRepositories(org string) ([]github.GitHubAPIRepository, error) {
	var repos []github.GitHubAPIRepository
	page := 1
	perPage := 100

	for {
		endpoint := fmt.Sprintf("/orgs/%s/repos?page=%d&per_page=%d&type=all", org, page, perPage)
		var pageRepos []github.GitHubAPIRepository

		if err := GitHubAPICall(endpoint, &pageRepos); err != nil {
			return nil, err
		}

		if len(pageRepos) == 0 {
			break
		}

		repos = append(repos, pageRepos...)

		if len(pageRepos) < perPage {
			break
		}

		page++
	}

	return repos, nil
}

// GetRepoTree fetches the repository tree structure
func GetRepoTree(repoFullName string, sha string) (*github.GitTree, error) {
	endpoint := fmt.Sprintf("/repos/%s/git/trees/%s?recursive=1", repoFullName, sha)
	var tree github.GitTree

	if err := GitHubAPICall(endpoint, &tree); err != nil {
		return nil, err
	}

	return &tree, nil
}

// GetFileContent fetches the content of a file from GitHub
func GetFileContent(repoFullName string, path string, ref string) ([]byte, error) {
	endpoint := fmt.Sprintf("/repos/%s/contents/%s?ref=%s", repoFullName, path, ref)
	var content github.GitHubContent

	if err := GitHubAPICall(endpoint, &content); err != nil {
		return nil, err
	}

	if content.Type != "file" {
		return nil, fmt.Errorf("path %s is not a file", path)
	}

	// Decode base64 content
	decoded, err := base64.StdEncoding.DecodeString(content.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file content: %w", err)
	}

	return decoded, nil
}

// GetDefaultBranchCommit gets the latest commit SHA for the default branch
func GetDefaultBranchCommit(repoFullName string, branch string) (string, error) {
	endpoint := fmt.Sprintf("/repos/%s/commits/%s", repoFullName, branch)
	var commit github.GitHubCommit

	if err := GitHubAPICall(endpoint, &commit); err != nil {
		return "", err
	}

	return commit.SHA, nil
}

// FilterRepositories filters repositories based on query patterns
func FilterRepositories(repos []github.GitHubAPIRepository, query string) []github.GitHubAPIRepository {
	if query == "" {
		return repos
	}

	patterns := strings.Split(query, ",")
	var filtered []github.GitHubAPIRepository

	for _, repo := range repos {
		for _, pattern := range patterns {
			pattern = strings.TrimSpace(pattern)
			if strings.Contains(repo.Name, pattern) {
				filtered = append(filtered, repo)
				break
			}
		}
	}

	return filtered
}
