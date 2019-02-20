package model

import "strings"

type Verifier interface {
	Valid() bool
}

type Cluster struct {
	SkipTLS bool   `yaml:"insecure-skip-tls-verify,omitempty"`
	Server  string `yaml:"server"`
	CaCert  string `yaml:"certificate-authority-data,omitempty"`
}

func (cluster *Cluster) Valid() bool {
	if cluster == nil {
		return false
	}
	if len(cluster.Server) <= 0 || !strings.HasPrefix(cluster.Server, "https") {
		return false
	} else if len(cluster.CaCert) <= 0 && !cluster.SkipTLS {
		return false
	}
	return true
}

type Clusters struct {
	Cluster Cluster `yaml:"cluster"`
	Name    string  `yaml:"name"`
}

func (clusters *Clusters) Valid() bool {
	if clusters == nil {
		return false
	}
	if len(clusters.Name) <= 0 {
		return false
	}
	return clusters.Cluster.Valid()
}

type Context struct {
	Cluster   string `yaml:"cluster"`
	Namespace string `yaml:"namespace"`
	User      string `yaml:"user"`
}

func (context *Context) Valid() bool {
	if context == nil {
		return false
	}
	if len(context.Cluster) <= 0 || len(context.User) <= 0 || len(context.Namespace) <= 0 {
		return false
	}
	return true
}

type Contexts struct {
	Context Context `yaml:"context"`
	Name    string  `yaml:"name"`
}

func (contests *Contexts) Valid() bool {
	if contests == nil {
		return false
	}
	if len(contests.Name) <= 0 {
		return false
	}
	return contests.Context.Valid()
}

type User struct {
	ClientCert string `yaml:"client-certificate-data,omitempty"`
	ClientKey  string `yaml:"client-key-data,omitempty"`
	Token      string `yaml:"token,omitempty"`
}

func (user *User) Valid() bool {
	if user == nil {
		return false
	}
	if len(user.Token) > 0 {
		return true
	} else if len(user.ClientKey) > 0 && len(user.ClientCert) > 0 {
		return true
	}
	return false

}

type Users struct {
	Name string `yaml:"name"`
	User User   `yaml:"user"`
}

func (user *Users) Valid() bool {
	if user == nil {
		return false
	}
	if len(user.Name) <= 0 {
		return false
	}
	return user.User.Valid()
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

type MidConfigWare struct {
	Cluster Clusters
	Context Contexts
	User    Users
}

func (ware *MidConfigWare) Valid() bool {
	if ware == nil {
		return false
	} else if ware.User.Valid() && ware.Context.Valid() && ware.Cluster.Valid() {
		return true
	}
	return false
}
