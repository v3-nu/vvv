package config

type Config struct {
	Settings Settings
	Aliases  map[string]string
}

type Settings struct {
	Packages struct {
		DefaultPackageManager string `yaml:"defaultPackageManager"`
	} `yaml:"packages"`
	Uploads struct {
		TransfershUrl      string `yaml:"transfershUrl"`
		TransfershUsername string `yaml:"transfershUsername"`
		TransfershPassword string `yaml:"transfershPassword"`
	} `yaml:"uploads"`
}
