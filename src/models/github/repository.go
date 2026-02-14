package github

// Repository represents a GitHub repository
type Repository struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	HTMLURL       string `json:"html_url"`
	Description   string `json:"description"`
	DefaultBranch string `json:"default_branch"`
	Private       bool   `json:"private"`
	Fork          bool   `json:"fork"`
	Archived      bool   `json:"archived"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	PushedAt      string `json:"pushed_at"`
	Size          int    `json:"size"`
	Language      string `json:"language"`
	ForksCount    int    `json:"forks_count"`
	StargazersCount int  `json:"stargazers_count"`
	WatchersCount int    `json:"watchers_count"`
	OpenIssuesCount int  `json:"open_issues_count"`
}

// GitHubAPIRepository represents the GitHub API response for a repository
type GitHubAPIRepository struct {
	ID          int    `json:"id"`
	NodeID      string `json:"node_id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Owner       struct {
		Login string `json:"login"`
		ID    int    `json:"id"`
	} `json:"owner"`
	Private       bool   `json:"private"`
	HTMLURL       string `json:"html_url"`
	Description   string `json:"description"`
	Fork          bool   `json:"fork"`
	URL           string `json:"url"`
	DefaultBranch string `json:"default_branch"`
	Archived      bool   `json:"archived"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	PushedAt      string `json:"pushed_at"`
	Size          int    `json:"size"`
	Language      string `json:"language"`
	ForksCount    int    `json:"forks_count"`
	StargazersCount int  `json:"stargazers_count"`
	WatchersCount int    `json:"watchers_count"`
	OpenIssuesCount int  `json:"open_issues_count"`
}

// GitTree represents a GitHub tree API response
type GitTree struct {
	SHA  string `json:"sha"`
	URL  string `json:"url"`
	Tree []struct {
		Path string `json:"path"`
		Mode string `json:"mode"`
		Type string `json:"type"`
		SHA  string `json:"sha"`
		Size int    `json:"size,omitempty"`
		URL  string `json:"url"`
	} `json:"tree"`
	Truncated bool `json:"truncated"`
}

// GitHubCommit represents a GitHub commit
type GitHubCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Date  string `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

// GitHubContent represents a file or directory content from GitHub API
type GitHubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	SHA         string `json:"sha"`
	Size        int    `json:"size"`
	URL         string `json:"url"`
	HTMLURL     string `json:"html_url"`
	GitURL      string `json:"git_url"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"` // "file" or "dir"
	Content     string `json:"content,omitempty"`
	Encoding    string `json:"encoding,omitempty"`
}

