package pipelinerunner

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type DatabaseConfiguration struct {
	FileName string
}

type Configuration struct {
	WorkDir  string
	Database DatabaseConfiguration
}

const (
	defaultConfigFile string = "conf.json"
	defaultDbFile     string = "pipelinerunner.db"
)

var (
	Config         Configuration
	defaultWorkDir string = fmt.Sprintf("%s/.pr", os.Getenv("HOME"))
)

// Configure the application
func Configure() error {
	defaultConfigFilePath := fmt.Sprintf("%s/%s", defaultWorkDir, defaultConfigFile)

	if _, err := os.Stat(defaultConfigFilePath); err != nil {
		// there is no configuration file at default location

		createWorkDirectory(defaultWorkDir)

		// create default configuration
		Config = Configuration{
			WorkDir:  defaultWorkDir,
			Database: DatabaseConfiguration{FileName: defaultDbFile},
		}

		saveConfiguration() // save configuration to the default location
	} else {
		// there is a configuration file at default location
		if err := ParseConfiguration(defaultConfigFilePath); err != nil {
			return fmt.Errorf("error parsing configuration file: %v", err)
		}
	}

	return nil
}

// Parse application configuration file
func ParseConfiguration(f string) error {
	file, err := os.Open(f)
	if err != nil {
		return fmt.Errorf("error opening configuration file: %v", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&Config); err != nil {
		return fmt.Errorf("error fetching configuration from the given file: %v", err)
	}

	// working directory is always equal to the default location, except for tests
	if !strings.HasPrefix(Config.WorkDir, "test") || strings.ContainsRune(Config.WorkDir, '.') {
		Config.WorkDir = defaultWorkDir
	}

	// save configuration file to the working directory if it does not exist there
	if _, err := os.Stat(fmt.Sprintf("%s/%s", Config.WorkDir, defaultConfigFile)); err != nil {
		createWorkDirectory(Config.WorkDir)
		saveConfiguration()
	}
	return nil
}

// Create application work directory at the given location
func createWorkDirectory(d string) error {
	if err := os.MkdirAll(d, 0777); err != nil {
		return fmt.Errorf("error creating working directory: %v", err)
	}
	return nil
}

// Save configuration to file in the work directory
func saveConfiguration() error {
	cbyte, _ := json.Marshal(Config)
	if err := os.WriteFile(fmt.Sprintf("%s/%s", Config.WorkDir, defaultConfigFile), cbyte, 0644); err != nil {
		return fmt.Errorf("error saving configuration to default location: %v", err)
	}
	return nil
}
