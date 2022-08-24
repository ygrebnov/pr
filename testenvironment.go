package pipelinerunner

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

var (
	testWorkDir string
	testConfig  Configuration
)

func setupTestEnv(t *testing.T, testDir string) {
	testWorkDir = fmt.Sprintf("test/%s", testDir)
	testConfigFilePath := fmt.Sprintf("%s/%s", testWorkDir, defaultConfigFile)

	// create test working directory
	os.MkdirAll(testWorkDir, 0777)

	// create test configuration
	testConfig = Configuration{
		WorkDir:  testWorkDir,
		Database: DatabaseConfiguration{FileName: defaultDbFile},
	}

	// save test configuration in working directory
	cbyte, _ := json.Marshal(testConfig)
	if err := os.WriteFile(testConfigFilePath, cbyte, 0644); err != nil {
		t.Errorf("Cannot create test configuration file: %v", testConfigFilePath)
	}

	// parse test configuration file (in order to populate Config struct, like it is done by the application)
	if err := ParseConfiguration(testConfigFilePath); err != nil {
		t.Errorf("error parsing test configuration file: %v", err)
	}

}

func teardownTestEnv(t *testing.T) {
	os.RemoveAll(testWorkDir)
}
