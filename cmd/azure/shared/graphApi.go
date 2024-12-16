package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GraphApiClient struct {
	BaseURL      string
	HTTPClient   *http.Client
	ClientID     string
	ClientSecret string
	TenantID     string
}

func NewGraphApiClient(clientID string, clientSecret string, tenantID string) *GraphApiClient {
	return &GraphApiClient{
		BaseURL:      "https://graph.microsoft.com/v1.0/",
		HTTPClient:   &http.Client{Timeout: 30 * time.Second},
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TenantID:     tenantID,
	}
}

func (client *GraphApiClient) GetAccessToken() (string, error) {
	url := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", client.TenantID)
	data := "client_id=" + client.ClientID +
		"&scope=https://graph.microsoft.com/.default" +
		"&client_secret=" + client.ClientSecret +
		"&grant_type=client_credentials"
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return "", fmt.Errorf("creating token request failed: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing token request failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decoding token response failed: %v", err)
	}

	accessToken, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("token response did not contain access_token")
	}

	return accessToken, nil
}
