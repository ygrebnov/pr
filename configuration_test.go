package pipelinerunner

import (
	"testing"
)

func TestConfiguration(t *testing.T) {
	setupTestEnv(t, "config")
	defer teardownTestEnv(t)

	if Config != testConfig {
		t.Errorf("got: %v, want:%v", Config, testConfig)
	}
}
