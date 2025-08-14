package types

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type KubeConfig struct {
	CurrentContext string                       `json:"current-context" yaml:"current-context"`
	Contexts       []KubeConfigContext          `json:"contexts" yaml:"contexts"`
	ContextsMap    map[string]KubeConfigContext `json:"-" yaml:"-"`
	Clusters       []KubeConfigCluster          `json:"clusters" yaml:"clusters"`
	ClustersMap    map[string]KubeConfigCluster `json:"-" yaml:"-"`
	Users          []KubeConfigUser             `json:"users" yaml:"users"`
	UsersMap       map[string]KubeConfigUser    `json:"-" yaml:"-"`
	Kind           string                       `json:"kind" yaml:"kind"`
	Preferences    map[string]string            `json:"preferences" yaml:"preferences"`
}

type KubeConfigContext struct {
	Name    string `json:"name" yaml:"name"`
	Context struct {
		Cluster   string `json:"cluster" yaml:"cluster"`
		User      string `json:"user" yaml:"user"`
		Namespace string `json:"namespace" yaml:"namespace"`
	} `json:"context" yaml:"context"`
}

type KubeConfigCluster struct {
	Name    string `json:"name" yaml:"name"`
	Cluster struct {
		Server   string `json:"server" yaml:"server"`
		Cert     string `json:"certificate-authority-data" yaml:"certificate-authority-data"`
		Insecure bool   `json:"insecure-skip-tls-verify" yaml:"insecure-skip-tls-verify"`
	} `json:"cluster" yaml:"cluster"`
}

type KubeConfigUser struct {
	Name string `json:"name" yaml:"name"`
	Auth struct {
		Token string `json:"token" yaml:"token"`
		Cert  string `json:"client-certificate" yaml:"client-certificate"`
		Key   string `json:"client-key" yaml:"client-key"`
	} `json:"user" yaml:"user"`
}

func TryGetKubeConfig() (*KubeConfig, error) {
	osvar := os.Getenv("KUBECONFIG")
	if osvar == "" {
		osvar = os.Getenv("HOME") + "/.kube/config"
	}
	if _, err := os.Stat(osvar); os.IsNotExist(err) {
		return nil, fmt.Errorf("kubeconfig file not found at %s", osvar)
	}
	data, err := os.ReadFile(osvar)
	if err != nil {
		return nil, fmt.Errorf("failed to read kubeconfig file: %w", err)
	}
	var kubeConfig KubeConfig
	if err := yaml.Unmarshal(data, &kubeConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kubeconfig: %w", err)
	}

	kubeConfig.ContextsMap = make(map[string]KubeConfigContext)
	for _, context := range kubeConfig.Contexts {
		kubeConfig.ContextsMap[context.Name] = context
	}

	kubeConfig.ClustersMap = make(map[string]KubeConfigCluster)
	for _, cluster := range kubeConfig.Clusters {
		kubeConfig.ClustersMap[cluster.Name] = cluster
	}

	kubeConfig.UsersMap = make(map[string]KubeConfigUser)
	for _, user := range kubeConfig.Users {
		kubeConfig.UsersMap[user.Name] = user
	}

	return &kubeConfig, nil
}
