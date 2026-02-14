package github

// Manifest represents the parsed manifest.yaml structure
type Manifest struct {
	Name         string                 `yaml:"name"`
	Environments map[string]Environment `yaml:"environments"`
	Asset        *Asset                 `yaml:"asset,omitempty"`
}

// Environment represents an environment configuration in the manifest
type Environment struct {
	Clusters []string `yaml:"clusters"`
	Replicas int      `yaml:"replicas,omitempty"`
	Expose   *Expose  `yaml:"expose,omitempty"`
}

// Asset represents the asset section for micro-frontends
type Asset struct {
	Environments map[string]AssetEnvironment `yaml:"environments"`
}

// AssetEnvironment represents an asset environment (micro-frontends)
type AssetEnvironment struct {
	SpecificJurisdictions []string `yaml:"specificJurisdictions,omitempty"`
}

// Expose represents the expose configuration
type Expose struct {
	OutsideCluster *OutsideCluster `yaml:"outsideCluster,omitempty"`
}

// OutsideCluster represents the outside cluster exposure configuration
type OutsideCluster struct {
	AmbassadorInternal []AmbassadorMapping `yaml:"ambassadorInternal,omitempty"`
	AmbassadorExternal []AmbassadorMapping `yaml:"ambassadorExternal,omitempty"`
}

// AmbassadorMapping represents an Ambassador mapping configuration
type AmbassadorMapping struct {
	Hostname string `yaml:"hostname,omitempty"`
	Prefix   string `yaml:"prefix,omitempty"`
	Rewrite  string `yaml:"rewrite,omitempty"`
}

