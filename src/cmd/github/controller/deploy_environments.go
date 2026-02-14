package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/github/be"
	githubHelper "github.com/eltiocaballoloco/sinaloa-cli/src/helpers/github"
)

func ReposDeployEnvironments(
	organization string,
	query string,
	envs string,
	manifestSection string,
	manifestName string,
	folders string,
	saveJSON bool,
	savePathJSON string,
) ([]byte, error) {

	// Parse comma-separated values
	envFilters := parseCommaSeparated(envs)
	folderList := parseCommaSeparated(folders)

	fmt.Printf("Starting deployment environment analysis...\n")
	fmt.Printf("Organization: %s\n", organization)
	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Environment filters: %v\n", envFilters)
	fmt.Printf("Folders to scan: %v\n", folderList)
	fmt.Printf("Manifest name: %s\n", manifestName)

	// Step 1: Fetch repositories
	fmt.Println("\n[1/4] Fetching repositories from GitHub...")

	// Show authentication method after first API call
	defer func() {
		authMethod := githubHelper.GetAuthMethod()
		if authMethod != "not authenticated yet" {
			fmt.Printf("ℹ️  GitHub authentication: %s\n", authMethod)
		}
	}()
	repos, err := githubHelper.ListRepositories(organization)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}
	fmt.Printf("Found %d repositories\n", len(repos))

	// Step 2: Filter repositories
	if query != "" {
		repos = githubHelper.FilterRepositories(repos, query)
		fmt.Printf("Filtered to %d repositories matching query\n", len(repos))
	}

	if len(repos) == 0 {
		return nil, fmt.Errorf("no repositories found matching criteria")
	}

	// Step 3: Process repositories concurrently
	fmt.Println("\n[2/4] Processing repositories and parsing manifests...")
	maxWorkers := 10 // Concurrent workers
	repoDataMap := be.ProcessRepositoriesConcurrently(
		repos,
		folderList,
		manifestName,
		envFilters,
		maxWorkers,
	)

	fmt.Printf("Successfully processed %d repositories with deployments\n", len(repoDataMap))

	if len(repoDataMap) == 0 {
		return nil, fmt.Errorf("no repositories found with matching deployments")
	}

	// Step 4: Build deployment matrix
	fmt.Println("\n[3/4] Building deployment matrix...")
	matrix := be.BuildDeploymentMatrix(
		repoDataMap,
		organization,
		query,
		envFilters,
		folderList,
	)

	fmt.Printf("Matrix generated with %d deploy keys and %d projects\n",
		matrix.Meta.Stats.TotalDeployKeys,
		matrix.Meta.Stats.TotalProjects,
	)

	// Step 5: Marshal to JSON
	fmt.Println("\n[4/4] Generating JSON output...")
	jsonData, err := json.MarshalIndent(matrix, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Save to file if requested
	if saveJSON && savePathJSON != "" {
		if err := os.WriteFile(savePathJSON, jsonData, 0644); err != nil {
			return nil, fmt.Errorf("failed to write JSON file: %w", err)
		}
		fmt.Printf("\n✓ JSON saved to: %s\n", savePathJSON)

		// Return summary instead of full JSON
		summary := fmt.Sprintf(`
Deployment Matrix Analysis Complete!
=====================================
Organization: %s
Repositories scanned: %d
Repositories with deployments: %d
Total projects: %d
Total deployments: %d
Total deploy keys: %d

Output saved to: %s
`,
			organization,
			matrix.Meta.Stats.TotalReposScanned,
			matrix.Meta.Stats.TotalReposWithDeployments,
			matrix.Meta.Stats.TotalProjects,
			matrix.Meta.Stats.TotalDeployments,
			matrix.Meta.Stats.TotalDeployKeys,
			savePathJSON,
		)

		return []byte(summary), nil
	}

	return jsonData, nil
}

// parseCommaSeparated parses a comma-separated string into a slice
func parseCommaSeparated(input string) []string {
	if input == "" {
		return []string{}
	}

	parts := strings.Split(input, ",")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
