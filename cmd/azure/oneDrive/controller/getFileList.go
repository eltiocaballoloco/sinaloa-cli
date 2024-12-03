package controller

import (
	"context"
	"fmt"

	"github.com/eltiocaballoloco/sinaloa-cli/cmd/azure/shared"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
)

func GetFileList(path string, userId string) ([]string, error) {
	// Initialize the GraphServiceClient
	graphClient, err := shared.NewGraphClient()
	if err != nil {
		return nil, fmt.Errorf("error initializing Graph client: %w", err)
	}

	fmt.Printf("Debug: path='%s'\n", path)

	// Construct the full URL for the folder
	fullPath := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/drive/root:/%s:/children", userId, path)

	// Use WithUrl to fetch the folder contents
	request := graphClient.Drives().WithUrl(fullPath)
	response, err := request.Get(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching children for path '%s': %w", path, err)
	}

	// Cast response to DriveItemCollectionResponseable
	children, ok := response.(msgraphmodels.DriveItemCollectionResponseable)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}

	// Prepare a list to store file/folder names
	var itemsList []string

	// Process and add items to the list
	if children.GetValue() != nil {
		for _, item := range children.GetValue() {
			itemName := *item.GetName()
			itemType := "File"
			if item.GetFolder() != nil {
				itemType = "Folder"
			}
			itemsList = append(itemsList, fmt.Sprintf("%s (%s)", itemName, itemType))
		}
	} else {
		return nil, fmt.Errorf("no items found in the specified folder")
	}

	return itemsList, nil
}
