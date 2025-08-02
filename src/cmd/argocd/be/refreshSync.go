package be

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
)

type ApplicationListResponse struct {
	Metadata struct {
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Items []Application `json:"items"`
}

type Application struct {
	Metadata struct {
		Name              string            `json:"name"`
		Namespace         string            `json:"namespace"`
		UID               string            `json:"uid"`
		ResourceVersion   string            `json:"resourceVersion"`
		Generation        int64             `json:"generation"`
		CreationTimestamp string            `json:"creationTimestamp"`
		Labels            map[string]string `json:"labels"`
		Finalizers        []string          `json:"finalizers"`
	} `json:"metadata"`
	Spec struct {
		Source struct {
			RepoURL        string `json:"repoURL"`
			Path           string `json:"path"`
			TargetRevision string `json:"targetRevision"`
			Plugin         struct {
				Env []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"env"`
			} `json:"plugin"`
		} `json:"source"`
		Destination struct {
			Namespace string `json:"namespace"`
			Name      string `json:"name"`
		} `json:"destination"`
		Project string `json:"project"`
	} `json:"spec"`
	Status struct {
		Sync struct {
			Status     string `json:"status"` // e.g., "Synced", "OutOfSync", "Unknown"
			ComparedTo struct {
				Source struct {
					RepoURL        string `json:"repoURL"`
					Path           string `json:"path"`
					TargetRevision string `json:"targetRevision"`
				} `json:"source"`
				Destination struct {
					Namespace string `json:"namespace"`
					Name      string `json:"name"`
				} `json:"destination"`
			} `json:"comparedTo"`
		} `json:"sync"`
		Health struct {
			Status             string `json:"status"` // e.g., "Healthy", "Degraded", "Missing"
			LastTransitionTime string `json:"lastTransitionTime"`
		} `json:"health"`
		Conditions []struct {
			Type               string `json:"type"`
			Message            string `json:"message"`
			LastTransitionTime string `json:"lastTransitionTime"`
		} `json:"conditions"`
		ReconciledAt         string `json:"reconciledAt"`
		ResourceHealthSource string `json:"resourceHealthSource"`
		ControllerNamespace  string `json:"controllerNamespace"`
	} `json:"status"`
}

// Simplified struct for status polling (backward compatibility)
type ApplicationStatus struct {
	Status struct {
		Sync struct {
			Status     string `json:"status"` // e.g., "Synced", "OutOfSync", "Unknown"
			ComparedTo struct {
				Source struct {
					RepoURL        string `json:"repoURL"`
					Path           string `json:"path"`
					TargetRevision string `json:"targetRevision"`
					Plugin         struct {
						Env []struct {
							Name  string `json:"name"`
							Value string `json:"value"`
						} `json:"env"`
					} `json:"plugin"`
				} `json:"source"`
				Destination struct {
					Namespace string `json:"namespace"`
					Name      string `json:"name"`
				} `json:"destination"`
			} `json:"comparedTo"`
		} `json:"sync"`
		Health struct {
			Status             string `json:"status"` // e.g., "Healthy", "Degraded", "Missing"
			LastTransitionTime string `json:"lastTransitionTime"`
		} `json:"health"`
		Conditions []struct {
			Type               string `json:"type"`
			Message            string `json:"message"`
			LastTransitionTime string `json:"lastTransitionTime"`
		} `json:"conditions"`
		OperationState *struct {
			Operation struct {
				Sync struct {
					Revision string `json:"revision"`
				} `json:"sync"`
				InitiatedBy struct {
					Username string `json:"username"`
				} `json:"initiatedBy"`
				Retry struct{} `json:"retry"`
			} `json:"operation"`
			Phase      string `json:"phase"`   // e.g., "Succeeded", "Error", "Failed", "Running"
			Message    string `json:"message"` // Detailed error message
			SyncResult struct {
				Revision string `json:"revision"`
				Source   struct {
					RepoURL        string `json:"repoURL"`
					Path           string `json:"path"`
					TargetRevision string `json:"targetRevision"`
					Plugin         struct {
						Env []struct {
							Name  string `json:"name"`
							Value string `json:"value"`
						} `json:"env"`
					} `json:"plugin"`
				} `json:"source"`
			} `json:"syncResult"`
			StartedAt  string `json:"startedAt"`
			FinishedAt string `json:"finishedAt"`
		} `json:"operationState,omitempty"` // Pointer to handle null values
		ReconciledAt         string   `json:"reconciledAt"`
		Summary              struct{} `json:"summary"`
		ResourceHealthSource string   `json:"resourceHealthSource"`
		ControllerNamespace  string   `json:"controllerNamespace"`
		SourceHydrator       struct{} `json:"sourceHydrator"`
	} `json:"status"`
}

type ArgoCDLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ArgoCDLoginResponse struct {
	Token string `json:"token"`
}

var argoClient *helpers.ApiClient

// LoginToArgoCD performs login to ArgoCD and returns the authentication token
func LoginToArgoCD(baseURL, username, password string) (string, error) {
	// Create a temporary client without authentication for login
	loginClient := helpers.NewApiClient(baseURL, "", "None")

	payload := ArgoCDLoginRequest{
		Username: username,
		Password: password,
	}

	res := loginClient.Request("POST", "/api/v1/session", payload)
	if !res.Response {
		return "", fmt.Errorf("ArgoCD login failed: %s", res.Message)
	}

	var loginResp ArgoCDLoginResponse
	err := json.Unmarshal(res.Body, &loginResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse ArgoCD login response: %v", err)
	}

	return loginResp.Token, nil
}

// InitArgoClientWithLogin performs login and initializes the ArgoCD client
func InitArgoClientWithLogin(baseURL, username, password string) error {
	token, err := LoginToArgoCD(baseURL, username, password)
	if err != nil {
		return err
	}

	argoClient = helpers.NewApiClient(baseURL, token, "Bearer")
	fmt.Printf("[Info] Successfully authenticated to ArgoCD\n")
	return nil
}

func GetAppNames(gitID, gitlabPath, env string) []string {
	var matchingNames []string
	var apps ApplicationListResponse

	// Construct the endpoint with proper URL formatting
	var endpoint string
	if gitlabPath != "" {
		endpoint = fmt.Sprintf("/api/v1/applications?repo=https://gitlab.com/%s.git", gitlabPath)
	} else {
		endpoint = "/api/v1/applications"
	}

	resp := argoClient.Request("GET", endpoint, nil)

	if !resp.Response {
		fmt.Printf("[Error] Fetching applications went wrong: %s\n", resp.Message)
		return nil
	}

	if err := json.Unmarshal(resp.Body, &apps); err != nil {
		fmt.Printf("[Error] Parsing response failed: %v\n", err)
		return nil
	}

	// Filter applications based on git_id and profile labels
	for _, app := range apps.Items {
		// Check if labels exist and match criteria
		if app.Metadata.Labels != nil {
			if app.Metadata.Labels["git_id"] == gitID && app.Metadata.Labels["profile"] == env {
				// Check if contains the env
				if strings.Contains(app.Metadata.Name, env) {
					matchingNames = append(matchingNames, app.Metadata.Name)
					fmt.Printf("[Info] Found matching app: %s (git_id: %s, profile: %s)\n",
						app.Metadata.Name, app.Metadata.Labels["git_id"], app.Metadata.Labels["profile"])
				}
			}
		}
	}

	if len(matchingNames) == 0 {
		fmt.Printf("[Info] No matching applications found for git_id: %s, profile: %s\n", gitID, env)
	}

	return matchingNames
}

// TriggerArgoHardRefreshAndSync triggers a hard refresh followed by sync
func TriggerArgoHardRefreshAndSync(appName string) error {
	// First, trigger hard refresh
	if err := TriggerArgoHardRefresh(appName); err != nil {
		return err
	}

	// Wait a moment for the refresh to complete
	time.Sleep(60 * time.Second)

	// Then trigger sync
	return TriggerArgoSync(appName)
}

// TriggerArgoHardRefresh triggers a hard refresh for the specified application
func TriggerArgoHardRefresh(appName string) error {
	endpoint := fmt.Sprintf("/api/v1/applications/%s", appName)

	resp := argoClient.Request("GET", endpoint+"?refresh=hard", nil)
	if !resp.Response {
		return fmt.Errorf("[Error] failed to trigger hard refresh for app %s: %s", appName, resp.Message)
	}

	fmt.Printf("[Info] Hard refresh triggered for application: %s\n", appName)
	return nil
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
