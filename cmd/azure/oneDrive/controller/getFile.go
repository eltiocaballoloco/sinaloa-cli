package controller

import (
	"context"
	"fmt"

	"github.com/eltiocaballoloco/sinaloa-cli/cmd/azure/shared"
)

func GetFileList(urlPath string) ([]byte, error) {
	// Initialize the GraphServiceClient
	graphClient, err := shared.NewGraphClient()
	if err != nil {
		return nil, fmt.Errorf("error initializing Graph client: %w", err)
	}

	// Fetch the root drive
	// https://learn.microsoft.com/en-us/graph/api/drive-get?view=graph-rest-1.0&tabs=http
	drive, err := graphClient.Me().Drive().Get(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching OneDrive root: %w", err)
	}

	// Use WithUrl to specify the full URL of the item
	itemRequest := graphClient.Drives().
		ByDriveId(*drive.GetId()).
		Items().
		WithUrl(urlPath)

	itemRequest.Get(context.Background(), nil)

	// Fetch children of the specified item
	// https://learn.microsoft.com/en-us/graph/api/driveitem-list-children?view=graph-rest-1.0&tabs=http
	// children, err := itemRequest.
	// if err != nil {
	// 	return nil, fmt.Errorf("error fetching files from URL '%s': %w", urlPath, err)
	// }

	// // Process and return file details
	// var fileList []byte
	// for _, file := range children.GetValue() {
	// 	fmt.Printf("Name: %s, ID: %s\n", *file.GetName(), *file.GetId())
	// 	fileList = append(fileList, []byte(*file.GetName())...)
	// 	fileList = append(fileList, '\n')
	// }

	// return itemRequest, nil

	return nil, nil
}
