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
		PastebinUrl        string `yaml:"pastebinUrl"`
		TransfershUrl      string `yaml:"transfershUrl"`
		TransfershUsername string `yaml:"transfershUsername"`
		TransfershPassword string `yaml:"transfershPassword"`
	} `yaml:"uploads"`
	Secenv struct {
		LocalInit    bool   `yaml:"localInit"`
		ServiceName  string `yaml:"serviceName"`
		KeychainName string `yaml:"keychainName"`
		SecurePath   string `yaml:"securePath"`
	} `yaml:"secenv"`
}
