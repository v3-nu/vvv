package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/scheiblingco/gofn/cfgtools"
	"github.com/scheiblingco/gofn/typetools"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

var DefaultSettingsLocation = "~/.clycli/settings.yaml"
var DefaultAliasConfigLocation = "~/.clycli/aliases.yaml"
var DefaultAliasRcLocation = "~/.clycli/.aliasrc"
var DefaultCacheLocation = "~/.cache/clycli"

type ConfigKey string

func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to get home directory: %v", err)
		}
		return homeDir + path[1:]
	}

	return path
}

func GetSettingsLocation() string {
	if os.Getenv("CLYCLI_SETTINGS_LOCATION") != "" {
		return os.Getenv("CLYCLI_SETTINGS_LOCATION")
	}

	return ExpandPath(DefaultSettingsLocation)
}

func GetAliasConfigLocation() string {
	if os.Getenv("CLYCLI_ALIAS_CONFIG_LOCATION") != "" {
		return os.Getenv("CLYCLI_ALIAS_CONFIG_LOCATION")
	}

	return ExpandPath(DefaultAliasConfigLocation)
}

func GetAliasRcLocation() string {
	if os.Getenv("CLYCLI_ALIAS_RC_LOCATION") != "" {
		return os.Getenv("CLYCLI_ALIAS_RC_LOCATION")
	}

	return ExpandPath(DefaultAliasRcLocation)
}

func GetCacheLocation() string {
	if os.Getenv("CLYCLI_CACHE_LOCATION") != "" {
		return os.Getenv("CLYCLI_CACHE_LOCATION")
	}

	return ExpandPath(DefaultCacheLocation)
}

func EnsureFileExists(path string, defaultContent []byte) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		folderName := filepath.Dir(path)
		if err := os.MkdirAll(folderName, 0o700); err != nil && !os.IsExist(err) {
			return err
		}

		err := os.WriteFile(path, defaultContent, 0o600)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadConfig() (*Config, error) {
	settingsLocation := GetSettingsLocation()
	aliasesLocation := GetAliasConfigLocation()

	EnsureFileExists(settingsLocation, []byte(""))
	EnsureFileExists(aliasesLocation, []byte(""))

	settings := Config{
		Aliases:  map[string]string{},
		Settings: Settings{},
	}

	err := cfgtools.LoadYamlConfig(settingsLocation, &settings.Settings)
	if err != nil {
		return nil, err
	}

	err = cfgtools.LoadYamlConfig(aliasesLocation, &settings.Aliases)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func (c *Config) SaveConfig() error {
	settingsLocation := GetSettingsLocation()
	aliasesLocation := GetAliasConfigLocation()

	settings := c.Settings
	aliases := c.Aliases

	yamlSettings, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}

	err = os.WriteFile(settingsLocation, yamlSettings, 0o600)
	if err != nil {
		return err
	}

	yamlAliases, err := yaml.Marshal(aliases)
	if err != nil {
		return err
	}

	err = os.WriteFile(aliasesLocation, yamlAliases, 0o600)
	if err != nil {
		return err
	}

	aliasRcLocation := GetAliasRcLocation()
	EnsureFileExists(aliasRcLocation, []byte(""))

	// Check if __clycli_aliases_registered is already registered as a function
	aliasRcContent := "if [ -z \"$(type -t __clycli_aliases_registered)\" ]; then \r\necho \"function already registered\"; \r\nreturn;\r\nfi\r\n"

	for alias, command := range aliases {
		aliasRcContent += fmt.Sprintf("alias %s='%s'\r\n", alias, command)
	}

	aliasRcContent += "\r\nfunction __clycli_aliases_registered() { return }\r\n"

	err = os.WriteFile(aliasRcLocation, []byte(aliasRcContent), 0o600)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) GetAvailableSettings() map[string][]string {
	availableSettings := map[string][]string{}

	ref := reflect.ValueOf(&c.Settings).Elem()
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		tag := ref.Type().Field(i).Tag.Get("yaml")
		availableSettings[tag] = []string{}

		for si := 0; si < field.NumField(); si++ {
			subField := field.Type().Field(si)
			availableSettings[tag] = append(availableSettings[tag], subField.Tag.Get("yaml"))
		}
	}

	return availableSettings
}

func (c *Config) GetCurrentSettingsMap() map[string]map[string]string {
	availableSettings := map[string]map[string]string{}

	ref := reflect.ValueOf(&c.Settings).Elem()
	for i := 0; i < ref.NumField(); i++ {
		field := ref.Field(i)
		tag := ref.Type().Field(i).Tag.Get("yaml")
		availableSettings[tag] = map[string]string{}

		for si := 0; si < field.NumField(); si++ {
			subField := field.Type().Field(si)
			value := field.Field(si).String()
			availableSettings[tag][subField.Tag.Get("yaml")] = value
		}
	}

	return availableSettings
}

func FindField(ref reflect.Value, field string) (string, error) {
	for i := 0; i < ref.NumField(); i++ {
		fieldName := ref.Type().Field(i).Name
		if strings.EqualFold(fieldName, field) {
			return fieldName, nil
		}

		fieldName = ref.Type().Field(i).Tag.Get("yaml")
		if strings.EqualFold(fieldName, field) {
			return fieldName, nil
		}
	}

	return "", fmt.Errorf("field %s not found", field)
}

func (c *Config) SetConfigValue(group string, key string, value string) error {
	ref := reflect.ValueOf(&c.Settings).Elem()

	groupFieldName, err := FindField(ref, group)
	if err != nil {
		return err
	}

	groupField := ref.FieldByName(groupFieldName)

	if !groupField.IsValid() {
		return fmt.Errorf("group %s does not exist", cases.Title(language.English).String(group))
	}

	keyFieldName, err := FindField(groupField, key)
	if err != nil {
		return err
	}

	keyField := groupField.FieldByName(keyFieldName)

	if !keyField.IsValid() {
		return fmt.Errorf("key %s does not exist", cases.Title(language.English).String(key))
	}

	if !keyField.CanSet() {
		return fmt.Errorf("key %s cannot be set", cases.Title(language.English).String(key))
	}

	fmt.Println("Setting", keyFieldName, "to", value)

	switch keyField.Kind() {
	case reflect.String:
		keyField.SetString(typetools.EnsureString(value))
	case reflect.Bool:
		keyField.SetBool(typetools.EnsureBool(value))
	case reflect.Float32:
	case reflect.Float64:
		keyField.SetFloat(typetools.EnsureFloat64(value))
	case reflect.Int:
	case reflect.Int64:
		keyField.SetInt(typetools.EnsureInt64(value))
	case reflect.Uint32:
	case reflect.Uint64:
		keyField.SetUint(typetools.EnsureUint64(value))
	default:
		keyField.Set(reflect.ValueOf(value))
	}

	return c.SaveConfig()
}
