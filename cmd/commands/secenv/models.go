package secenv

import "github.com/99designs/keyring"

type SecureEnvironment struct {
	Name string `json:"name" yaml:"name"`
	Desc string `json:"desc" yaml:"desc"`

	Vars  map[string]SecureItem `json:"vars" yaml:"vars"`
	Files map[string]SecureItem `json:"files" yaml:"files"`

	config keyring.Config  `json:"-" yaml:"-"`
	ring   keyring.Keyring `json:"-" yaml:"-"`
}

type ItemType string

const (
	// ItemTypeFile represents a file item type
	ItemTypeFile ItemType = "file"

	// ItemTypeEnv represents an environment variable item type
	ItemTypeEnv ItemType = "env"
)

type SecureItem struct {
	Name string   `json:"name" yaml:"name"`
	Kind ItemType `json:"kind" yaml:"kind"`
	Desc string   `json:"desc" yaml:"desc"`
	Data []byte   `json:"data" yaml:"data"`
}
