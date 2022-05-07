package kubernetes

type KubeConfig struct {
	ApiVersion     string `json:"apiVersion,omitempty"`
	Kind           string `json:"kind,omitempty"`
	CurrentContext string `json:"current-context,omitempty"`
	Clusters       []struct {
		Name    string `json:"name,omitempty"`
		Cluster struct {
			Server                   string `json:"server,omitempty"`
			CertificateAuthorityData string `json:"certificate-authority-data,omitempty"`
		} `json:"cluster,omitempty"`
	} `json:"clusters,omitempty"`
	Users []struct {
		Name string `json:"name,omitempty"`
		User struct {
			Token string `json:"token,omitempty"`
		} `json:"user,omitempty"`
	} `json:"users,omitempty"`
	Contexts []struct {
		Name    string `json:"name,omitempty"`
		Context struct {
			User    string `json:"user,omitempty"`
			Cluster string `json:"cluster,omitempty"`
		} `json:"context,omitempty"`
	} `json:"contexts,omitempty"`
}
