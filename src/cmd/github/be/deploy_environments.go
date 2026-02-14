package be

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	githubHelper "github.com/eltiocaballoloco/sinaloa-cli/src/helpers/github"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models/github"
	"gopkg.in/yaml.v3"
)

// ParseEnvironmentName parses an environment name into prefix and region
func ParseEnvironmentName(envName string) (prefix string, region string) {
	parts := strings.Split(envName, "-")

	if len(parts) >= 2 {
		// e.g., "prod-euc1" -> prefix="prod-", region="euc1"
		prefix = parts[0] + "-"
		region = parts[1]
	} else {
		// e.g., "qa" -> prefix="qa-", region="default"
		prefix = parts[0] + "-"
		region = "default"
	}

	return prefix, region
}

// BuildDeployKey creates a deploy key from prefix and region
func BuildDeployKey(prefix, region string) string {
	return prefix + "|" + region
}

// MatchesEnvFilter checks if an environment name matches any of the filter patterns
func MatchesEnvFilter(envName string, envFilters []string) bool {
	if len(envFilters) == 0 {
		return true
	}

	prefix, _ := ParseEnvironmentName(envName)

	for _, filter := range envFilters {
		if strings.HasPrefix(prefix, filter) || prefix == filter {
			return true
		}
	}

	return false
}

// ProcessRepository processes a single repository and extracts deployment information
func ProcessRepository(
	repo github.GitHubAPIRepository,
	folders []string,
	manifestName string,
	envFilters []string,
) (*github.RepoData, error) {

	// Get the latest commit SHA
	commitSHA, err := githubHelper.GetDefaultBranchCommit(repo.FullName, repo.DefaultBranch)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit SHA: %w", err)
	}

	// Get repository tree
	tree, err := githubHelper.GetRepoTree(repo.FullName, commitSHA)
	if err != nil {
		return nil, fmt.Errorf("failed to get repo tree: %w", err)
	}

	// Find manifest files in target folders
	manifestPaths := findManifestPaths(tree, folders, manifestName)

	if len(manifestPaths) == 0 {
		return nil, nil // No manifests found, skip this repo
	}

	// Process each manifest
	subprojects := make(map[string]github.Subproject)

	for _, manifestPath := range manifestPaths {
		subproject, err := processManifest(repo, manifestPath, commitSHA, envFilters)
		if err != nil {
			fmt.Printf("Warning: failed to process manifest %s: %v\n", manifestPath, err)
			continue
		}

		if subproject != nil {
			subprojects[subproject.ProjectID] = *subproject
		}
	}

	if len(subprojects) == 0 {
		return nil, nil // No valid subprojects found
	}

	repoData := &github.RepoData{
		RepoID:        repo.FullName,
		RepoName:      repo.Name,
		RepoURL:       repo.HTMLURL,
		DefaultBranch: repo.DefaultBranch,
		LastCommitSHA: commitSHA,
		Subprojects:   subprojects,
	}

	return repoData, nil
}

// findManifestPaths finds all manifest.yaml files in the specified folders
func findManifestPaths(tree *github.GitTree, folders []string, manifestName string) []string {
	var paths []string

	for _, item := range tree.Tree {
		if item.Type != "blob" {
			continue
		}

		// Check if the file is a manifest in one of the target folders
		for _, folder := range folders {
			if strings.HasPrefix(item.Path, folder+"/") && strings.HasSuffix(item.Path, "/"+manifestName) {
				paths = append(paths, item.Path)
				break
			}
		}
	}

	return paths
}

// processManifest processes a single manifest file
func processManifest(
	repo github.GitHubAPIRepository,
	manifestPath string,
	commitSHA string,
	envFilters []string,
) (*github.Subproject, error) {

	// Fetch manifest content
	content, err := githubHelper.GetFileContent(repo.FullName, manifestPath, commitSHA)
	if err != nil {
		return nil, fmt.Errorf("failed to get manifest content: %w", err)
	}

	// Parse YAML
	var manifest github.Manifest
	if err := yaml.Unmarshal(content, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Extract project info
	folderPath := strings.TrimSuffix(manifestPath, "/manifest.yaml")
	projectID := manifest.Name
	if projectID == "" {
		// Use folder name as fallback
		parts := strings.Split(folderPath, "/")
		projectID = parts[len(parts)-1]
	}

	// Process environments
	deployments := make(map[string]github.DeploymentDetail)
	summary := make(map[string][]string)

	for envName, env := range manifest.Environments {
		// Skip if doesn't match env filters
		if !MatchesEnvFilter(envName, envFilters) {
			continue
		}

		// Skip if no clusters (Pattern 3: micro-frontends)
		if len(env.Clusters) == 0 {
			continue
		}

		// Parse environment name
		prefix, region := ParseEnvironmentName(envName)
		deployKey := BuildDeployKey(prefix, region)

		// Build deployment detail
		deployment := github.DeploymentDetail{
			EnvPrefix: prefix,
			Region:    region,
			Cluster:   env.Clusters[0], // Take first cluster
			Namespace: buildNamespace(projectID, envName, region),
		}

		// Extract hosts and URLs if available
		if env.Expose != nil && env.Expose.OutsideCluster != nil {
			deployment.Hosts, deployment.Paths, deployment.URLs = extractExposureInfo(env.Expose.OutsideCluster, envName)
		}

		deployments[deployKey] = deployment

		// Add to summary
		if summary[prefix] == nil {
			summary[prefix] = []string{}
		}
		if !contains(summary[prefix], region) {
			summary[prefix] = append(summary[prefix], region)
		}
	}

	if len(deployments) == 0 {
		return nil, nil // No matching deployments
	}

	subproject := &github.Subproject{
		ProjectID:    projectID,
		DisplayName:  manifest.Name,
		FolderPath:   folderPath,
		ManifestPath: manifestPath,
		Summary:      summary,
		Deployments:  deployments,
		SearchBlob:   buildSearchBlob(repo.FullName, projectID, manifest.Name, summary, deployments),
	}

	return subproject, nil
}

// buildNamespace constructs the namespace from project and environment info
func buildNamespace(projectID, envName, region string) string {
	// Pattern: area-gaming-<project>-<env>-<region>
	// or area-gaming-<project>-<env> for default region

	// Extract service name from project ID (e.g., "gpi-pragmaticplay" -> "gpi")
	parts := strings.Split(projectID, "-")
	service := parts[0]

	if region == "default" {
		return fmt.Sprintf("area-gaming-%s-%s", service, envName)
	}

	return fmt.Sprintf("area-gaming-%s-%s", service, envName)
}

// extractExposureInfo extracts hosts, paths, and URLs from exposure configuration
func extractExposureInfo(outsideCluster *github.OutsideCluster, envName string) (
	hosts map[string]string,
	paths map[string]string,
	urls map[string]string,
) {
	hosts = make(map[string]string)
	paths = make(map[string]string)
	urls = make(map[string]string)

	// Extract from ambassadorInternal
	if len(outsideCluster.AmbassadorInternal) > 0 {
		for i, mapping := range outsideCluster.AmbassadorInternal {
			if mapping.Hostname != "" {
				if i == 0 {
					hosts["ambassador_internal"] = mapping.Hostname
				}

				if mapping.Prefix != "" {
					key := "main_prefix"
					if strings.Contains(mapping.Prefix, "ping") {
						key = "ping_prefix"
					}
					paths[key] = mapping.Prefix

					// Build URL
					urlKey := "internal_main"
					if strings.Contains(mapping.Prefix, "ping") {
						urlKey = "internal_ping"
					}
					urls[urlKey] = fmt.Sprintf("https://%s%s", mapping.Hostname, mapping.Prefix)
				}
			}
		}
	}

	// Extract from ambassadorExternal
	if len(outsideCluster.AmbassadorExternal) > 0 {
		for i, mapping := range outsideCluster.AmbassadorExternal {
			if mapping.Hostname != "" {
				if i == 0 {
					hosts["ambassador_external"] = mapping.Hostname
				}

				if mapping.Prefix != "" {
					key := "external_main"
					if strings.Contains(mapping.Prefix, "ping") {
						key = "external_ping"
					}
					urls[key] = fmt.Sprintf("https://%s%s", mapping.Hostname, mapping.Prefix)
				}
			}
		}
	}

	return hosts, paths, urls
}

// buildSearchBlob creates a searchable string from all relevant fields
func buildSearchBlob(repoID, projectID, displayName string, summary map[string][]string, deployments map[string]github.DeploymentDetail) string {
	var parts []string

	parts = append(parts, repoID)
	parts = append(parts, projectID)
	parts = append(parts, displayName)

	// Add env prefixes
	for prefix := range summary {
		parts = append(parts, prefix)
	}

	// Add regions
	for _, regions := range summary {
		parts = append(parts, regions...)
	}

	// Add clusters and namespaces
	for _, deployment := range deployments {
		parts = append(parts, deployment.Cluster)
		parts = append(parts, deployment.Namespace)
	}

	return strings.ToLower(strings.Join(parts, " "))
}

// contains checks if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// ProcessRepositoriesConcurrently processes multiple repositories concurrently
func ProcessRepositoriesConcurrently(
	repos []github.GitHubAPIRepository,
	folders []string,
	manifestName string,
	envFilters []string,
	maxWorkers int,
) map[string]github.RepoData {

	jobs := make(chan github.GitHubAPIRepository, len(repos))
	results := make(chan *github.RepoData, len(repos))

	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < maxWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for repo := range jobs {
				fmt.Printf("Processing repository: %s\n", repo.FullName)
				repoData, err := ProcessRepository(repo, folders, manifestName, envFilters)
				if err != nil {
					fmt.Printf("Error processing %s: %v\n", repo.FullName, err)
					continue
				}
				if repoData != nil {
					results <- repoData
				}
			}
		}()
	}

	// Send jobs
	for _, repo := range repos {
		jobs <- repo
	}
	close(jobs)

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	repoDataMap := make(map[string]github.RepoData)
	for repoData := range results {
		repoDataMap[repoData.RepoID] = *repoData
	}

	return repoDataMap
}

// BuildDeploymentMatrix builds the complete deployment matrix from repo data
func BuildDeploymentMatrix(
	repoDataMap map[string]github.RepoData,
	organization string,
	query string,
	envFilters []string,
	folders []string,
) *github.DeploymentMatrix {

	// Collect all unique deploy keys
	deployKeySet := make(map[string]bool)
	envPrefixSet := make(map[string]bool)
	regionSet := make(map[string]bool)

	for _, repoData := range repoDataMap {
		for _, subproject := range repoData.Subprojects {
			for deployKey := range subproject.Deployments {
				deployKeySet[deployKey] = true

				parts := strings.Split(deployKey, "|")
				if len(parts) == 2 {
					envPrefixSet[parts[0]] = true
					regionSet[parts[1]] = true
				}
			}
		}
	}

	// Convert sets to sorted slices
	deployKeys := make([]string, 0, len(deployKeySet))
	for key := range deployKeySet {
		deployKeys = append(deployKeys, key)
	}
	sort.Slice(deployKeys, func(i, j int) bool {
		return compareDeployKeys(deployKeys[i], deployKeys[j])
	})

	envPrefixes := make([]string, 0, len(envPrefixSet))
	for prefix := range envPrefixSet {
		envPrefixes = append(envPrefixes, prefix)
	}
	sort.Strings(envPrefixes)

	regions := make([]string, 0, len(regionSet))
	for region := range regionSet {
		regions = append(regions, region)
	}
	sort.Strings(regions)

	// Build deploy key labels and groups
	deployKeyLabels := make(map[string]string)
	deployKeyGroups := make(map[string]string)
	for _, deployKey := range deployKeys {
		parts := strings.Split(deployKey, "|")
		if len(parts) == 2 {
			prefix := parts[0]
			region := parts[1]

			if region == "default" {
				deployKeyLabels[deployKey] = strings.TrimSuffix(prefix, "-")
			} else {
				deployKeyLabels[deployKey] = prefix + region
			}

			deployKeyGroups[deployKey] = strings.TrimSuffix(prefix, "-")
		}
	}

	// Build matrix table
	matrixTable := buildMatrixTable(repoDataMap, deployKeys)

	// Build indices
	rowIndex := buildRowIndex(matrixTable.Rows)

	// Calculate stats
	totalProjects := 0
	totalDeployments := 0
	for _, repoData := range repoDataMap {
		totalProjects += len(repoData.Subprojects)
		for _, subproject := range repoData.Subprojects {
			totalDeployments += len(subproject.Deployments)
		}
	}

	matrix := &github.DeploymentMatrix{
		Meta: github.Meta{
			SchemaVersion: 1,
			GeneratedAt:   time.Now().Format(time.RFC3339),
			Source: github.Source{
				Type:         "repo-scan",
				Organization: organization,
				Query:        query,
				EnvFilters:   envFilters,
				Folders:      folders,
			},
			Stats: github.Stats{
				TotalReposScanned:         len(repoDataMap),
				TotalReposWithDeployments: len(repoDataMap),
				TotalProjects:             totalProjects,
				TotalDeployments:          totalDeployments,
				TotalDeployKeys:           len(deployKeys),
			},
		},
		Dimensions: github.Dimensions{
			EnvPrefixes:     envPrefixes,
			Regions:         regions,
			DeployKeys:      deployKeys,
			DeployKeyLabels: deployKeyLabels,
			DeployKeyGroups: deployKeyGroups,
		},
		Repos: repoDataMap,
		Tables: github.Tables{
			MatrixGlobal: matrixTable,
			RowIndex:     rowIndex,
		},
	}

	return matrix
}

// compareDeployKeys compares two deploy keys for sorting
func compareDeployKeys(a, b string) bool {
	partsA := strings.Split(a, "|")
	partsB := strings.Split(b, "|")

	if len(partsA) != 2 || len(partsB) != 2 {
		return a < b
	}

	// First compare by group (prod, qa, test)
	groupA := strings.TrimSuffix(partsA[0], "-")
	groupB := strings.TrimSuffix(partsB[0], "-")

	if groupA != groupB {
		// Define order: prod, qa, test, dev
		order := map[string]int{"prod": 0, "qa": 1, "test": 2, "dev": 3}
		orderA, okA := order[groupA]
		orderB, okB := order[groupB]

		if okA && okB {
			return orderA < orderB
		}
		return groupA < groupB
	}

	// Then compare by region
	return partsA[1] < partsB[1]
}

// buildMatrixTable builds the matrix table from repo data
func buildMatrixTable(repoDataMap map[string]github.RepoData, deployKeys []string) github.MatrixGlobal {
	var rows []github.MatrixRow

	// Build column definitions
	columns := make([]github.ColumnDef, len(deployKeys))
	for i, deployKey := range deployKeys {
		parts := strings.Split(deployKey, "|")
		label := deployKey
		group := ""
		region := ""

		if len(parts) == 2 {
			group = strings.TrimSuffix(parts[0], "-")
			region = parts[1]

			if region == "default" {
				label = group
			} else {
				label = parts[0] + region
			}
		}

		columns[i] = github.ColumnDef{
			DeployKey: deployKey,
			Label:     label,
			Group:     group,
			Region:    region,
		}
	}

	// Build rows
	for _, repoData := range repoDataMap {
		for projectID, subproject := range repoData.Subprojects {
			rowID := repoData.RepoID + "::" + projectID

			cells := make(map[string]bool)
			for _, deployKey := range deployKeys {
				_, exists := subproject.Deployments[deployKey]
				cells[deployKey] = exists
			}

			row := github.MatrixRow{
				RowID:       rowID,
				RepoID:      repoData.RepoID,
				RepoName:    repoData.RepoName,
				ProjectID:   projectID,
				ProjectName: subproject.DisplayName,
				Cells:       cells,
				SearchBlob:  subproject.SearchBlob,
			}

			rows = append(rows, row)
		}
	}

	// Sort rows by repo name, then project name
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].RepoName != rows[j].RepoName {
			return rows[i].RepoName < rows[j].RepoName
		}
		return rows[i].ProjectName < rows[j].ProjectName
	})

	return github.MatrixGlobal{
		Columns: columns,
		Rows:    rows,
	}
}

// buildRowIndex builds indices for fast lookup
func buildRowIndex(rows []github.MatrixRow) github.RowIndex {
	byRepo := make(map[string][]string)
	byProjectID := make(map[string][]string)
	byDeployKey := make(map[string][]string)

	for _, row := range rows {
		// Index by repo
		byRepo[row.RepoID] = append(byRepo[row.RepoID], row.RowID)

		// Index by project ID
		byProjectID[row.ProjectID] = append(byProjectID[row.ProjectID], row.RowID)

		// Index by deploy key
		for deployKey, hasDeployment := range row.Cells {
			if hasDeployment {
				byDeployKey[deployKey] = append(byDeployKey[deployKey], row.RowID)
			}
		}
	}

	return github.RowIndex{
		ByRepo:      byRepo,
		ByProjectID: byProjectID,
		ByDeployKey: byDeployKey,
	}
}
