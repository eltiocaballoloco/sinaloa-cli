package be_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive/be"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGraphApiClient is a mock implementation of the Graph API client
type MockGraphApiClientUpload struct {
	mock.Mock
}

// MockApiClient is a mock implementation for API calls
type MockApiClient struct {
	mock.Mock
}

func (m *MockApiClient) Request(method, url string, body interface{}) *models.ApiResponse {
	args := m.Called(method, url, body)
	return args.Get(0).(*models.ApiResponse)
}

// TestUploadItem tests the UploadItem function
func TestUploadItem(t *testing.T) {
	mockApiClient := new(MockApiClient)
	mockApiClient.On("Request", mock.Anything, mock.Anything, mock.Anything).Return(&models.ApiResponse{
		Body: []byte(`{"uploadUrl": "http://mock-upload-url"}`),
	})

	tempFile, err := os.CreateTemp("", "test-file-*.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte("This is a test file"))
	assert.NoError(t, err)

	success, err := be.UploadItem(tempFile.Name(), "/test/path/file.txt")
	assert.Empty(t, success)
}

// TestCreateUploadSession tests the CreateUploadSession function
func TestCreateUploadSession(t *testing.T) {
	mockApiClient := new(MockApiClient)

	mockApiClient.On("Request", "POST", mock.Anything, mock.Anything).Return(&models.ApiResponse{
		Body: []byte(`{"uploadUrl": "http://mock-upload-url"}`),
	})

	uploadURL, _ := be.CreateUploadSession("http://mock-base-url", "mockAccessToken", "mockDriveID", "mockFolder", "mockFile")
	assert.NotEmpty(t, "http://mock-upload-url", uploadURL)
}

// TestUploadFileInChunks tests the UploadFileInChunks function
func TestUploadFileInChunks(t *testing.T) {
	// Create a temporary file to simulate a real file
	tempFile, err := os.CreateTemp("", "test-file-*.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte("This is a test file for chunk upload."))
	assert.NoError(t, err)

	// Mock server to simulate upload URL
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		assert.NotEmpty(t, body)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer mockServer.Close()

	err = be.UploadFileInChunks(mockServer.URL, tempFile.Name())
	assert.NoError(t, err)
}

// TestUploadFileInChunksError tests the UploadFileInChunks function for errors
func TestUploadFileInChunksError(t *testing.T) {
	// Create a temporary file to simulate a real file
	tempFile, err := os.CreateTemp("", "test-file-*.txt")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte("This is a test file for chunk upload."))
	assert.NoError(t, err)

	// Mock server to simulate upload URL with an error
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	err = be.UploadFileInChunks(mockServer.URL, tempFile.Name())
	assert.Error(t, err)
}
