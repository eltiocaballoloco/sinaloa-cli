package github

// DeploymentMatrix represents the complete output JSON structure
type DeploymentMatrix struct {
	Meta       Meta                  `json:"meta"`
	Dimensions Dimensions            `json:"dimensions"`
	Repos      map[string]RepoData   `json:"repos"`
	Tables     Tables                `json:"tables"`
}

// Meta contains metadata about the scan
type Meta struct {
	SchemaVersion int    `json:"schema_version"`
	GeneratedAt   string `json:"generated_at"`
	Source        Source `json:"source"`
	Stats         Stats  `json:"stats"`
}

// Source contains information about the scan source
type Source struct {
	Type         string   `json:"type"`
	Organization string   `json:"organization"`
	Query        string   `json:"query"`
	EnvFilters   []string `json:"env_filters"`
	Folders      []string `json:"folders"`
}

// Stats contains statistics about the scan
type Stats struct {
	TotalReposScanned          int `json:"total_repos_scanned"`
	TotalReposWithDeployments  int `json:"total_repos_with_deployments"`
	TotalProjects              int `json:"total_projects"`
	TotalDeployments           int `json:"total_deployments"`
	TotalDeployKeys            int `json:"total_deploy_keys"`
}

// Dimensions contains all unique dimensions found
type Dimensions struct {
	EnvPrefixes      []string          `json:"env_prefixes"`
	Regions          []string          `json:"regions"`
	DeployKeys       []string          `json:"deploy_keys"`
	DeployKeyLabels  map[string]string `json:"deploy_key_labels"`
	DeployKeyGroups  map[string]string `json:"deploy_key_groups"`
}

// RepoData contains all data for a single repository
type RepoData struct {
	RepoID        string                `json:"repo_id"`
	RepoName      string                `json:"repo_name"`
	RepoURL       string                `json:"repo_url"`
	DefaultBranch string                `json:"default_branch"`
	LastCommitSHA string                `json:"last_commit_sha"`
	Subprojects   map[string]Subproject `json:"subprojects"`
}

// Subproject represents a project within a repository
type Subproject struct {
	ProjectID    string                     `json:"project_id"`
	DisplayName  string                     `json:"display_name"`
	FolderPath   string                     `json:"folder_path"`
	ManifestPath string                     `json:"manifest_path"`
	Summary      map[string][]string        `json:"summary"`
	Deployments  map[string]DeploymentDetail `json:"deployments"`
	SearchBlob   string                     `json:"search_blob"`
}

// DeploymentDetail contains detailed deployment information
type DeploymentDetail struct {
	EnvPrefix string            `json:"env_prefix"`
	Region    string            `json:"region"`
	Cluster   string            `json:"cluster"`
	Namespace string            `json:"namespace"`
	Hosts     map[string]string `json:"hosts,omitempty"`
	Paths     map[string]string `json:"paths,omitempty"`
	URLs      map[string]string `json:"urls,omitempty"`
}

// Tables contains pre-computed tables for frontend
type Tables struct {
	MatrixGlobal MatrixGlobal       `json:"matrix_global"`
	RowIndex     RowIndex           `json:"row_index"`
}

// MatrixGlobal represents the global matrix table
type MatrixGlobal struct {
	Columns []ColumnDef `json:"columns"`
	Rows    []MatrixRow `json:"rows"`
}

// ColumnDef defines a column in the matrix
type ColumnDef struct {
	DeployKey string `json:"deploy_key"`
	Label     string `json:"label"`
	Group     string `json:"group"`
	Region    string `json:"region"`
}

// MatrixRow represents a row in the matrix table
type MatrixRow struct {
	RowID       string          `json:"row_id"`
	RepoID      string          `json:"repo_id"`
	RepoName    string          `json:"repo_name"`
	ProjectID   string          `json:"project_id"`
	ProjectName string          `json:"project_name"`
	Cells       map[string]bool `json:"cells"`
	SearchBlob  string          `json:"search_blob"`
}

// RowIndex contains indices for fast lookup
type RowIndex struct {
	ByRepo       map[string][]string `json:"by_repo"`
	ByProjectID  map[string][]string `json:"by_project_id"`
	ByDeployKey  map[string][]string `json:"by_deploy_key"`
}

