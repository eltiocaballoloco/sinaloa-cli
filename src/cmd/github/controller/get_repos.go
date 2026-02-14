package controller

import (
	"encoding/json"
	"fmt"
	"os"

	githubHelper "github.com/eltiocaballoloco/sinaloa-cli/src/helpers/github"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models/github"
)

func GetRepos(organization string, query string, saveJSON bool, savePathJSON string) ([]byte, error) {
	// Fetch repositories from GitHub
	fmt.Printf("Fetching repositories from organization: %s\n", organization)
	repos, err := githubHelper.ListRepositories(organization)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}

	fmt.Printf("Found %d repositories\n", len(repos))

	// Filter repositories if query is provided
	if query != "" {
		repos = githubHelper.FilterRepositories(repos, query)
		fmt.Printf("Filtered to %d repositories matching query: %s\n", len(repos), query)
	}

	// Convert to simpler structure for output
	var outputRepos []github.Repository
	for _, repo := range repos {
		outputRepos = append(outputRepos, github.Repository{
			ID:              fmt.Sprintf("%d", repo.ID),
			Name:            repo.Name,
			FullName:        repo.FullName,
			HTMLURL:         repo.HTMLURL,
			Description:     repo.Description,
			DefaultBranch:   repo.DefaultBranch,
			Private:         repo.Private,
			Fork:            repo.Fork,
			Archived:        repo.Archived,
			CreatedAt:       repo.CreatedAt,
			UpdatedAt:       repo.UpdatedAt,
			PushedAt:        repo.PushedAt,
			Size:            repo.Size,
			Language:        repo.Language,
			ForksCount:      repo.ForksCount,
			StargazersCount: repo.StargazersCount,
			WatchersCount:   repo.WatchersCount,
			OpenIssuesCount: repo.OpenIssuesCount,
		})
	}

	// Marshal to JSON
	var jsonData []byte
	if saveJSON {
		jsonData, err = json.MarshalIndent(outputRepos, "", "  ")
	} else {
		jsonData, err = json.MarshalIndent(outputRepos, "", "  ")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Save to file if requested
	if saveJSON && savePathJSON != "" {
		if err := os.WriteFile(savePathJSON, jsonData, 0644); err != nil {
			return nil, fmt.Errorf("failed to write JSON file: %w", err)
		}
		fmt.Printf("JSON saved to: %s\n", savePathJSON)
		return []byte(fmt.Sprintf("Successfully saved %d repositories to %s", len(outputRepos), savePathJSON)), nil
	}

	return jsonData, nil
}

