package be

import (
	"encoding/json"
	"fmt"

	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
)

type Application struct {
	Metadata struct {
		Name   string            `json:"name"`
		Labels map[string]string `json:"labels"`
	} `json:"metadata"`
}

type ApplicationListResponse struct {
	Items []Application `json:"items"`
}

type ApplicationStatus struct {
	Status struct {
		Sync struct {
			Status string `json:"status"` // e.g., "Synced", "OutOfSync", "Unknown"
		} `json:"sync"`
		Health struct {
			Status string `json:"status"` // e.g., "Healthy", "Degraded", "Missing"
		} `json:"health"`
	} `json:"status"`
}

var argoClient *helpers.ApiClient

func InitArgoClient(baseURL, token string) {
	argoClient = helpers.NewApiClient(baseURL, token, "Bearer")
}

func GetAppNames(gitID, env string) []string {
	endpoint := "/api/v1/applications"
	resp := argoClient.Request("GET", endpoint, nil)

	if !resp.Response {
		fmt.Printf("[Error] Fetching applications went wrong: %s\n", resp.Message)
		return nil
	}

	var apps ApplicationListResponse
	if err := json.Unmarshal(resp.Body, &apps); err != nil {
		fmt.Printf("[Error] Parsing response failed: %v\n", err)
		return nil
	}

	var matchingNames []string
	for _, app := range apps.Items {
		if app.Metadata.Labels["git_id"] == gitID && app.Metadata.Labels["profile"] == env {
			matchingNames = append(matchingNames, app.Metadata.Name)
		}
	}

	if len(matchingNames) == 0 {
		fmt.Println("[Info] No matching applications found")
	}
	return matchingNames
}

// TriggerArgoSync triggers a sync operation for the specified application.
func TriggerArgoSync(appName string) error {
	endpoint := fmt.Sprintf("/api/v1/applications/%s/sync", appName)

	// You can customize sync options here if needed
	body := map[string]interface{}{
		"revision": "", // empty means latest
		"prune":    false,
		"dryRun":   false,
	}

	resp := argoClient.Request("POST", endpoint, body)
	if !resp.Response {
		return fmt.Errorf("[Error] failed to trigger sync for app %s: %s", appName, resp.Message)
	}

	fmt.Printf("[Info] Sync triggered for application: %s\n", appName)
	return nil
}

// GetArgoAppStatus retrieves the sync and health status of the application.
func GetArgoAppStatus(appName string) (*ApplicationStatus, error) {
	endpoint := fmt.Sprintf("/api/v1/applications/%s", appName)
	resp := argoClient.Request("GET", endpoint, nil)

	if !resp.Response {
		return nil, fmt.Errorf("[Error] failed to get status for app %s: %s", appName, resp.Message)
	}

	var status ApplicationStatus
	if err := json.Unmarshal(resp.Body, &status); err != nil {
		return nil, fmt.Errorf("[Error] failed to parse status for app %s: %v", appName, err)
	}

	return &status, nil
}
