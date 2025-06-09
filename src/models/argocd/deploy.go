package argocd

type ArgoCDDeployParams struct {
	AppName    string `json:"ARGOCD_APP_NAME"`
	Namespace  string `json:"ARGOCD_APP_NAMESPACE"`
	SourcePath string `json:"ARGOCD_APP_SOURCE_PATH"`
	RepoURL    string `json:"ARGOCD_APP_SOURCE_REPO_URL"`
	Revision   string `json:"ARGOCD_APP_REVISION"`

	Profile      string `json:"ARGOCD_ENV_PROFILE"`
	Module       string `json:"ARGOCD_ENV_MODULE"`
	Tag          string `json:"ARGOCD_ENV_TAG"`
	ExtraSecrets string `json:"ARGOCD_ENV_EXTRA_SECRETS"`
	ChartName    string `json:"ARGOCD_ENV_CHART_NAME"`
	ChartRepo    string `json:"ARGOCD_ENV_CHART_REPO"`
	ChartParams  string `json:"ARGOCD_ENV_CHART_PARAMS"`
	ReleaseName  string `json:"ARGOCD_ENV_RELEASE_NAME"`

	DockerRepo string `json:"ARGOCD_EXTRA_DOCKER_REPO"`
}
