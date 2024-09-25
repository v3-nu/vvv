package config_test

import (
	"fmt"
	"testing"

	"github.com/clysec/clycli/config"
)

func TestGetAvailableSettings(t *testing.T) {
	config := config.Config{
		Settings: config.Settings{
			Packages: struct {
				DefaultPackageManager string "yaml:\"defaultPackageManager\""
			}{
				DefaultPackageManager: "apt",
			},
		},
	}
	availableSettings := config.GetCurrentSettingsMap()
	fmt.Println(availableSettings)
}
