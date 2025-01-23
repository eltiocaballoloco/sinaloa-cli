package be

import (
	"encoding/json"
	"fmt"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/shared"
	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models/azure"
)

// GetDriveItems uses the helpers.ApiClient to request items from a specific path in OneDrive
func GetDriveItems(path string) (models.ApiResponse, error) {
	// Declare variables
	var endpoint string
	var apiGraph azure.OneDriveGraphResponseApiModel
	var items []azure.OneDriveItemModel

	// Load configuration from .env
	helpers.LoadConfig()

	// Initialize the GraphApiClient
	graphApiClient := shared.NewGraphApiClient(
		helpers.AppConfig.AZURE_CLIENT_ID,
		helpers.AppConfig.AZURE_CLIENT_SECRET,
		helpers.AppConfig.AZURE_TENANT_ID,
	)

	// Get the access token from the GraphApiClient
	accessToken, err := graphApiClient.GetAccessToken()
	if err != nil {
		errorMessage := fmt.Sprintf("GetDriveItems, internal error getting access token: %v\n", err)
		return models.NewApiResponse(false, 500, nil, errorMessage, nil), err
	}

	// Initialize the ApiClient with the Microsoft Graph API base URL, access token, and Bearer auth type
	baseURL := graphApiClient.BaseURL + "drives/"
	apiClient := helpers.NewApiClient(baseURL, accessToken, "Bearer")

	// Set the full endpoint URL using the DRIVE ID and path
	if path == "." {
		endpoint = fmt.Sprintf("%s/root/children", helpers.AppConfig.AZURE_DRIVE_ID)
	} else {
		endpoint = fmt.Sprintf("%s/root:/%s:/children", helpers.AppConfig.AZURE_DRIVE_ID, path)
	}

	// Use the existing request method from ApiClient to make the GET request
	apiResponse := apiClient.Request("GET", endpoint, nil)

	// Unmarshal the response body into the apiGraph struct
	err = json.Unmarshal(apiResponse.Body, &apiGraph)
	if err != nil {
		errorMessage := fmt.Sprintf("GetDriveItems, internal error unmarshalling response body from graph api: %v\n", err)
		return models.NewApiResponse(false, 500, nil, errorMessage, nil), err
	}

	// Create a slice of OneDriveItemModel structs
	for _, item := range apiGraph.Value {
		if item.Folder != nil {
			item.Type = "folder"
		} else {
			item.Type = "item"
		}
		// Concatenate path with name
		item.ParentReference.Path = fmt.Sprintf("%s/%s", item.ParentReference.Path, item.Name)
		items = append(items, item)
	}

	// Wrap items in a ResponseWrapper
	wrappedItems := azure.OneDriveWrapperModel{
		Values: items,
	}

	// Convert items slice to JSON bytes
	itemsJson, err := json.Marshal(wrappedItems)
	if err != nil {
		return models.NewApiResponse(false, 500, nil, "GetDriveItems, internal error marshalling oneDive items to JSON", nil), err
	}

	return models.NewApiResponse(true, 200, nil, "Data obtained from onedrive", itemsJson), nil
}
