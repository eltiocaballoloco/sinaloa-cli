package controller

import (
	"fmt"
	"strings"
	"time"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/argocd/be"
)

func RefreshSync(
	gitId string,
	env string,
	regions string,
	argocdUrl string,
	argocdToken string,
) error {
	// 1. Init client and get all apps matching gitId/env
	be.InitArgoClient("https://"+argocdUrl, argocdToken)
	appNames := be.GetAppNames(gitId, env)
	if len(appNames) == 0 {
		return fmt.Errorf("[Error] no applications found for gitId: %s and env: %s", gitId, env)
	}

	// 2. If no regions specified, sync all apps in appNames as-is
	if strings.TrimSpace(regions) == "" {
		for _, appName := range appNames {
			fmt.Printf("[Info] Syncing app: %s\n", appName)
			if err := syncAppWithPolling(appName); err != nil {
				return err
			}
		}
	} else {
		// 3. Otherwise, sync apps by ordered regions list
		orderedRegions := strings.Split(regions, ",")

		for _, region := range orderedRegions {
			for _, appName := range appNames {
				if strings.HasSuffix(appName, "-"+region) {
					fmt.Printf("[Info] Syncing app: %s\n", appName)
					if err := syncAppWithPolling(appName); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func syncAppWithPolling(appName string) error {
	if err := be.TriggerArgoSync(appName); err != nil {
		return fmt.Errorf("[Error] Failed to trigger sync for %s: %v", appName, err)
	}

	for {
		status, err := be.GetArgoAppStatus(appName)
		if err != nil {
			return fmt.Errorf("[Error] Failed to get status for %s: %v", appName, err)
		}

		if status.Status.Sync.Status == "Synced" && status.Status.Health.Status == "Healthy" {
			fmt.Printf("[Info] App %s successfully synced and healthy\n", appName)
			break
		}

		if status.Status.Health.Status == "Degraded" || status.Status.Sync.Status == "Failed" {
			return fmt.Errorf("[Error] Sync failed for %s", appName)
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}
