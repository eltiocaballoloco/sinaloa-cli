package shared

import (
	"encoding/json"
	"fmt"

	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
)

type DockerHubLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DockerHubLoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// LoginToDockerHub logs into Docker Hub and returns the auth token and refresh token.
func LoginToDockerHub(username, password string) (string, string, error) {
	client := helpers.NewApiClient("https://hub.docker.com", "", "None")
	payload := DockerHubLoginRequest{
		Username: username,
		Password: password,
	}

	res := client.Request("POST", "/v2/users/login", payload)
	if !res.Response {
		return "", "", fmt.Errorf("login failed: %s", res.Message)
	}

	var loginResp DockerHubLoginResponse
	err := json.Unmarshal(res.Body, &loginResp)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse login response: %v", err)
	}

	return loginResp.Token, loginResp.RefreshToken, nil
}
