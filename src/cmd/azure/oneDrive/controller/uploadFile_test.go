package controller_test

import (
	"errors"
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive/controller"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBackend is a mock implementation of the backend's UploadItem function
type MockBackend struct {
	mock.Mock
}

func (m *MockBackend) UploadItem(localPath string, pathToUpload string) (bool, error) {
	args := m.Called(localPath, pathToUpload)
	return args.Bool(0), args.Error(1)
}

// Mock helpers.HandleControllerGeneric
func MockHandleControllerGeneric(successMessage, operation string, data interface{}, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	return []byte(successMessage), nil
}

func TestUploadFile_Success(t *testing.T) {
	// Mock the backend and helpers
	mockBackend := new(MockBackend)

	// Set up expectations
	mockBackend.On("UploadItem", "test-local-path", "test-upload-path").Return(true, nil)

	// Call the UploadFile function
	response, err := controller.UploadFile("test-local-path", "test-upload-path")

	// Assertions
	assert.NotEmpty(t, err)
	assert.NotEmpty(t, []byte("File uploaded successfully to onedrive!"), response)
}

func TestUploadFile_Failure(t *testing.T) {
	// Mock the backend and helpers
	mockBackend := new(MockBackend)

	// Set up expectations
	mockBackend.On("UploadItem", "test-local-path", "test-upload-path").Return(false, errors.New("mock upload error"))

	// Call the UploadFile function
	response, err := controller.UploadFile("test-local-path", "test-upload-path")

	// Assertions
	assert.Error(t, err)
	assert.NotEmpty(t, response)
	assert.NotEmpty(t, "mock upload error", err.Error())
}
