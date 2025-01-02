package azure

import (
	"encoding/json"
)

// ApiResponse encapsulates the common properties of an API response.
type OneDriveGraphResponseApiModel struct {
	ODataContext string              `json:"@odata.context"`
	Tips         string              `json:"@microsoft.graph.tips"`
	Value        []OneDriveItemModel `json:"value"`
}

type OneDriveWrapperModel struct {
	Values []OneDriveItemModel `json:"values"`
}

// Item represents a single item, which could be a file or folder.
type OneDriveItemModel struct {
	CreatedDateTime           string          `json:"createdDateTime"`
	LastModifiedDateTime      string          `json:"lastModifiedDateTime"`
	ID                        string          `json:"id"`
	Name                      string          `json:"name"`
	WebUrl                    string          `json:"webUrl"`
	Size                      int             `json:"size"`
	CTag                      string          `json:"cTag"`
	ParentReference           ParentReference `json:"parentReference"`
	FileSystemInfo            FileSystemInfo  `json:"fileSystemInfo"`
	Folder                    *Folder         `json:"folder"`
	File                      *File           `json:"file"`
	Shared                    *Shared         `json:"shared"`
	MicrosoftGraphDownloadUrl string          `json:"@microsoft.graph.downloadUrl,omitempty"`
	DownloadUrl               string          `json:"downloadUrl,omitempty"`
	Type                      string          `json:"type"`
}

// ParentReference holds the reference to the parent item.
type ParentReference struct {
	DriveType string `json:"driveType"`
	DriveID   string `json:"driveId"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	SiteID    string `json:"siteId"`
}

// FileSystemInfo holds timestamps of creation and modification.
type FileSystemInfo struct {
	CreatedDateTime      string `json:"createdDateTime"`
	LastModifiedDateTime string `json:"lastModifiedDateTime"`
}

// Folder details if the item is a folder.
type Folder struct {
	ChildCount int `json:"childCount"`
}

// File details if the item is a file.
type File struct {
	MimeType string `json:"mimeType"`
	Hashes   Hashes `json:"hashes"`
}

// Hashes contains file hash information.
type Hashes struct {
	QuickXorHash string `json:"quickXorHash"`
}

// Shared details about the sharing status of the item.
type Shared struct {
	Scope string `json:"scope"`
}

// Custom UnmarshalJSON method for OneDriveItemModel
func (o *OneDriveItemModel) UnmarshalJSON(data []byte) error {
	// Create an alias to avoid recursion
	type Alias OneDriveItemModel
	aux := &struct {
		MicrosoftGraphDownloadUrl string `json:"@microsoft.graph.downloadUrl"` // Original field name from API
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	// Unmarshal the JSON into the auxiliary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Map the original field to the renamed field
	o.DownloadUrl = aux.MicrosoftGraphDownloadUrl

	return nil
}
