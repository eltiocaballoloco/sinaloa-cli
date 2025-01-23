package azure_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/models/azure"

	"github.com/stretchr/testify/assert"
)

func TestOneDriveItemModelUnmarshalJSON(t *testing.T) {
	// Arrange: JSON input with the MicrosoftGraphDownloadUrl field
	jsonInput := `{
		"createdDateTime": "2025-01-01T12:00:00Z",
		"lastModifiedDateTime": "2025-01-02T12:00:00Z",
		"id": "12345",
		"name": "testfile.txt",
		"webUrl": "https://example.com/testfile.txt",
		"size": 1024,
		"cTag": "c:12345",
		"@microsoft.graph.downloadUrl": "https://download.example.com/testfile.txt",
		"parentReference": {
			"driveType": "personal",
			"driveId": "67890",
			"id": "abcdef",
			"name": "Parent Folder",
			"path": "/drive/root:/Parent Folder",
			"siteId": "site-123"
		},
		"fileSystemInfo": {
			"createdDateTime": "2025-01-01T12:00:00Z",
			"lastModifiedDateTime": "2025-01-02T12:00:00Z"
		},
		"file": {
			"mimeType": "text/plain",
			"hashes": {
				"quickXorHash": "abc123"
			}
		},
		"type": "file"
	}`

	// Act: Unmarshal into a OneDriveItemModel
	var item azure.OneDriveItemModel
	err := json.Unmarshal([]byte(jsonInput), &item)

	// Assert: Validate unmarshaled data
	assert.NoError(t, err, "Unmarshal should not return an error")
	assert.Equal(t, "12345", item.ID, "ID should match")
	assert.Equal(t, "testfile.txt", item.Name, "Name should match")
	assert.Equal(t, "https://download.example.com/testfile.txt", item.DownloadUrl, "DownloadUrl should match the @microsoft.graph.downloadUrl")
	assert.Equal(t, "text/plain", item.File.MimeType, "MimeType should match")
	assert.Equal(t, "abc123", item.File.Hashes.QuickXorHash, "QuickXorHash should match")
	assert.Equal(t, "personal", item.ParentReference.DriveType, "Parent DriveType should match")
}

func TestOneDriveGraphResponseApiModel(t *testing.T) {
	// Arrange: JSON input for the response model
	jsonInput := `{
		"@odata.context": "https://graph.microsoft.com/v1.0/$metadata#driveItem",
		"@microsoft.graph.tips": "Some tips",
		"value": [{
			"createdDateTime": "2025-01-01T12:00:00Z",
			"lastModifiedDateTime": "2025-01-02T12:00:00Z",
			"id": "12345",
			"name": "testfile.txt",
			"webUrl": "https://example.com/testfile.txt",
			"size": 1024
		}]
	}`

	// Act: Unmarshal into a OneDriveGraphResponseApiModel
	var response azure.OneDriveGraphResponseApiModel
	err := json.Unmarshal([]byte(jsonInput), &response)

	// Assert: Validate unmarshaled data
	assert.NoError(t, err, "Unmarshal should not return an error")
	assert.Equal(t, "https://graph.microsoft.com/v1.0/$metadata#driveItem", response.ODataContext, "ODataContext should match")
	assert.Equal(t, "Some tips", response.Tips, "Tips should match")
	assert.Len(t, response.Value, 1, "Response value should have one item")
	assert.Equal(t, "testfile.txt", response.Value[0].Name, "Name of first item should match")
}

func TestOneDriveWrapperModel(t *testing.T) {
	// Arrange: JSON input for the wrapper model
	jsonInput := `{
		"values": [{
			"createdDateTime": "2025-01-01T12:00:00Z",
			"id": "12345",
			"name": "testfile.txt"
		}]
	}`

	// Act: Unmarshal into a OneDriveWrapperModel
	var wrapper azure.OneDriveWrapperModel
	err := json.Unmarshal([]byte(jsonInput), &wrapper)

	// Assert: Validate unmarshaled data
	assert.NoError(t, err, "Unmarshal should not return an error")
	assert.Len(t, wrapper.Values, 1, "Wrapper values should have one item")
	assert.Equal(t, "12345", wrapper.Values[0].ID, "ID of first item should match")
	assert.Equal(t, "testfile.txt", wrapper.Values[0].Name, "Name of first item should match")
}

func TestOneDriveUploadSessionModel_UnmarshalJSON(t *testing.T) {
	// Sample JSON response
	jsonResponse := `{
        "@odata.context": "https://graph.microsoft.com/v1.0/$metadata#microsoft.graph.uploadSession",
        "expirationDateTime": "2025-01-23T17:37:05.007Z",
        "nextExpectedRanges": ["0-"],
        "uploadUrl": "https://example.com/upload"
    }`

	// Expected struct
	expected := azure.OneDriveUploadSessionModel{
		ODataContext:       "https://graph.microsoft.com/v1.0/$metadata#microsoft.graph.uploadSession",
		ExpirationDateTime: "2025-01-23T17:37:05.007Z",
		NextExpectedRanges: []string{"0-"},
		UploadUrl:          "https://example.com/upload",
	}

	// Unmarshal JSON into the struct
	var result azure.OneDriveUploadSessionModel
	err := json.Unmarshal([]byte(jsonResponse), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Compare the result with the expected struct
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unmarshalled struct does not match expected.\nGot: %+v\nExpected: %+v", result, expected)
	}
}
