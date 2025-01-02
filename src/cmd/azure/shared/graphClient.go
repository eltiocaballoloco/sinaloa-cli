package shared

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"

	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
)

type GraphClient struct {
	clientSecretCredential *azidentity.ClientSecretCredential
	appClient              *msgraphsdk.GraphServiceClient
}

// NewGraphClient initializes a new GraphClient and returns its instance
// https://learn.microsoft.com/en-us/graph/tutorials/go-app-only?tabs=aad
func NewGraphClient() (*msgraphsdk.GraphServiceClient, error) {
	client := &GraphClient{}
	err := client.initializeGraphForAppAuth()
	if err != nil {
		return nil, err
	}
	return client.GetClient(), nil
}

// initializeGraphForAppAuth sets up the Graph client for application authentication
func (g *GraphClient) initializeGraphForAppAuth() error {
	// Load configuration from .env
	helpers.LoadConfig()

	// Retrieve Azure secrets from AppConfig
	clientId := helpers.AppConfig.AZURE_CLIENT_ID
	tenantId := helpers.AppConfig.AZURE_TENANT_ID
	clientSecret := helpers.AppConfig.AZURE_CLIENT_SECRET

	if clientId == "" || tenantId == "" || clientSecret == "" {
		return fmt.Errorf("Graph Client is missing required environment variables: AZURE_CLIENT_ID, AZURE_TENANT_ID, or AZURE_CLIENT_SECRET")
	}

	// Initialize Azure Identity credential
	credential, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		return fmt.Errorf("On Graph Client failed to create client secret credential: %w", err)
	}

	g.clientSecretCredential = credential

	// Create an authentication provider using Azure Identity
	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(credential, []string{
		"https://graph.microsoft.com/.default",
	})
	if err != nil {
		return fmt.Errorf("On Graph Client failed to create auth provider: %w", err)
	}

	// Create a request adapter
	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return fmt.Errorf("On Graph Client failed to create request adapter: %w", err)
	}

	// Initialize GraphServiceClient
	g.appClient = msgraphsdk.NewGraphServiceClient(adapter)

	return nil
}

// GetClient exposes the GraphServiceClient for external use
func (g *GraphClient) GetClient() *msgraphsdk.GraphServiceClient {
	return g.appClient
}
