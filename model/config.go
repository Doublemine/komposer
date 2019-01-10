package model

type cluster struct {
	SkipTLS bool   `yaml:"insecure-skip-tls-verify"`
	Server  string `yaml:"server"`
	CaCert  string `yaml:"certificate-authority-data"`
}

type clusters struct {
	Cluster []cluster `yaml:"cluster"`
	Name    string    `yaml:"name"`
}

type context struct {
	Cluster   string `yaml:"cluster"`
	Namespace string `yaml:"namespace"`
	User      string `yaml:"user"`
}

type contexts struct {
	Context context `yaml:"context"`
	Name    string  `yaml:"name"`
}

type config struct {
	ApiVersion     string      `yaml:"apiVersion"`
	Kind           string      `yaml:"kind"`
	CurrentContext string      `yaml:"current-context"`
	Preferences    interface{} `yaml:"preferences"`
	Contexts       contexts    `yaml:"contexts"`
}
