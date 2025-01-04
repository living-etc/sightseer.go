package kubernetes

type KubernetesVersionInfo struct {
	Full  string
	Major string `json:"major"`
	Minor string `json:"minor"`
}
