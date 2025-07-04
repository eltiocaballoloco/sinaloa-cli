package docker

// Structs to unmarshal JSON response from dockerhub api
type Image struct {
	Architecture string  `json:"architecture"`
	Features     string  `json:"features"`
	Variant      *string `json:"variant"`
	Digest       string  `json:"digest"`
	OS           string  `json:"os"`
	OSFeatures   string  `json:"os_features"`
	OSVersion    *string `json:"os_version"`
	Size         int64   `json:"size"`
	Status       string  `json:"status"`
	LastPulled   string  `json:"last_pulled"`
	LastPushed   string  `json:"last_pushed"`
}

type TagResult struct {
	Creator             int     `json:"creator"`
	ID                  int64   `json:"id"`
	Images              []Image `json:"images"`
	LastUpdated         string  `json:"last_updated"`
	LastUpdater         int     `json:"last_updater"`
	LastUpdaterUsername string  `json:"last_updater_username"`
	Name                string  `json:"name"`
	Repository          int     `json:"repository"`
	FullSize            int64   `json:"full_size"`
	V2                  bool    `json:"v2"`
	TagStatus           string  `json:"tag_status"`
	TagLastPulled       string  `json:"tag_last_pulled"`
	TagLastPushed       string  `json:"tag_last_pushed"`
	MediaType           string  `json:"media_type"`
	ContentType         string  `json:"content_type"`
	Digest              string  `json:"digest"`
}

type TagsResponse struct {
	Count    int         `json:"count"`
	Next     *string     `json:"next"`
	Previous *string     `json:"previous"`
	Results  []TagResult `json:"results"`
}

// TagInfo struct to hold detailed information about a tag (internal logic)
type TagInfoInternal struct {
	Name                string `json:"name"`
	IDImage             int64  `json:"id_image"`
	IDRepository        int    `json:"id_repository"`
	IDCreator           int    `json:"id_creator"`
	LastUpdaterUsername string `json:"last_updater_username"`
	LastUpdated         string `json:"last_updated"`
	TagLastPulled       string `json:"tag_last_pulled"`
	TagLastPushed       string `json:"tag_last_pushed"`
	Digest              string `json:"digest"`
}

type TagResponseInternal struct {
	TagList []TagInfoInternal `json:"tags_list"`
	Count   int               `json:"count"`
}

// Model used for argocd deploy
type DockerHubResponse struct {
	Response bool                `json:"response"`
	Code     string              `json:"code"`
	Message  string              `json:"message"`
	Data     TagResponseInternal `json:"data"`
}
