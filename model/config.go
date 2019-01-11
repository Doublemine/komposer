package model

type Cluster struct {
	SkipTLS bool   `yaml:"insecure-skip-tls-verify,omitempty"`
	Server  string `yaml:"server"`
	CaCert  string `yaml:"certificate-authority-data,omitempty"`
}

type Clusters struct {
	Cluster Cluster `yaml:"cluster"`
	Name    string  `yaml:"name"`
}

type Context struct {
	Cluster   string `yaml:"cluster"`
	Namespace string `yaml:"namespace"`
	User      string `yaml:"user"`
}

type Contexts struct {
	Context Context `yaml:"context"`
	Name    string  `yaml:"name"`
}

type User struct {
	ClientCert string `yaml:"client-certificate-data,omitempty"`
	ClientKey  string `yaml:"client-key-data,omitempty"`
	Token      string `yaml:"token,omitempty"`
}

type Users struct {
	Name string `yaml:"name"`
	User User   `yaml:"user"`
}

type Config struct {
	ApiVersion     string      `yaml:"apiVersion"`
	Kind           string      `yaml:"kind"`
	CurrentContext string      `yaml:"current-context"`
	Preferences    interface{} `yaml:"preferences"`
	Clusters       []Clusters  `yaml:"clusters"`
	Contexts       []Contexts  `yaml:"contexts"`
	Users          []Users     `yaml:"users"`
}
