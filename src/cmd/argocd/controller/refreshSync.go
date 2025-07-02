package controller

import (
	"fmt"
	"strings"
	"time"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/argocd/be"
)

func RefreshSync(
	gitId string,
	gitlabPath string,
	env string,
	regions string,
	argocdUrl string,
	argocdUsername string,
	argocdPassword string,
) (string, []string, error) {
	// 1. Login to ArgoCD and init client
	if err := be.InitArgoClientWithLogin("https://"+argocdUrl, argocdUsername, argocdPassword); err != nil {
		return "error", nil, fmt.Errorf("[Error] Failed to authenticate to ArgoCD: %v", err)
	}

	// 2. Get all apps matching gitId/env
	appNames := be.GetAppNames(gitId, gitlabPath, env)
	if len(appNames) == 0 {
		return "error", nil, fmt.Errorf("[Error] no applications found for gitId: %s and env: %s", gitId, env)
	}

	// 3. If no regions specified, sync all apps in appNames as-is
	if strings.TrimSpace(regions) == "" {
		for _, appName := range appNames {
			fmt.Printf("[Info] Refreshing and syncing app: %s\n", appName)
			if status, err := syncAppWithPolling(appName); err != nil {
				return status, appNames, err
			}
		}
	} else {
		// 4. Otherwise, sync apps by ordered regions list
		orderedRegions := strings.Split(regions, ",")
		// For each region, sync all apps that end with "-<region>"
		for _, region := range orderedRegions {
			for _, appName := range appNames {
				if strings.HasSuffix(appName, "-"+region) {
					fmt.Printf("[Info] Refreshing and Syncing app: %s\n", appName)
					if status, err := syncAppWithPolling(appName); err != nil {
						return status, appNames, err
					}
				}
			}
		}
	}
	return "error", appNames, nil
}

func syncAppWithPolling(appName string) (string, error) {
	// First perform hard refresh, then sync
	if err := be.TriggerArgoHardRefreshAndSync(appName); err != nil {
		return "error", fmt.Errorf("[Error] Failed to trigger hard refresh and sync for %s: %v", appName, err)
	}

	for {
		status, err := be.GetArgoAppStatus(appName)
		if err != nil {
			return "error", fmt.Errorf("[Error] Failed to get status for %s: %v", appName, err)
		}

		// Check operationState first (most important for sync operations)
		if status.Status.OperationState != nil {
			// Log current operationState for debugging with detailed info
			phase := status.Status.OperationState.Phase
			fmt.Printf("[Info] App %s - OperationState Phase: '%s'\n", appName, phase)

			if phase == "Error" || phase == "Failed" {
				return "error", fmt.Errorf("[Error] Sync operation failed for %s: %s", appName, status.Status.OperationState.Message)
			}

			// Check if sync completed successfully with operationState
			if status.Status.OperationState.Phase == "Succeeded" &&
				status.Status.Sync.Status == "Synced" &&
				status.Status.Health.Status == "Healthy" {
				fmt.Printf("[Info] App %s successfully synced and healthy\n", appName)
				break
			}
		} else {
			// Fallback to basic sync/health check when operationState is null
			if status.Status.Sync.Status == "Synced" && status.Status.Health.Status == "Healthy" {
				fmt.Printf("[Info] App %s successfully synced and healthy (no operationState)\n", appName)
				break
			}
		}

		// Check for degraded health or failed sync status
		if status.Status.Health.Status == "Degraded" || status.Status.Sync.Status == "Failed" {
			return "error", fmt.Errorf("[Error] Sync failed for %s - Health: %s, Sync: %s",
				appName, status.Status.Health.Status, status.Status.Sync.Status)
		}

		// Log current status for debugging
		if status.Status.OperationState != nil {
			fmt.Printf("[Info] App %s - Phase: %s, Sync: %s, Health: %s\n",
				appName, status.Status.OperationState.Phase, status.Status.Sync.Status, status.Status.Health.Status)
		} else {
			fmt.Printf("[Info] App %s - Sync: %s, Health: %s (no operationState)\n",
				appName, status.Status.Sync.Status, status.Status.Health.Status)
		}

		time.Sleep(5 * time.Second)
	}

	return "ok", nil
}
