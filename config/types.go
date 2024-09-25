package config

type Config struct {
	Settings Settings
	Aliases  map[string]string
}

type Settings struct {
	Packages struct {
		DefaultPackageManager string `yaml:"defaultPackageManager"`
	} `yaml:"packages"`
}
