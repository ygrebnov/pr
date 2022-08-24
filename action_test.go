package pipelinerunner

import (
	"testing"
)

func TestNewAction(t *testing.T) {
	actionType := RUN
	async := false
	attr := map[string]string{"Base": "echo \"What a lovely day\" > whatalovelyday"}
	a := NewAction(actionType, async, attr)

	if a.ID == "" {
		t.Error("ID must be non-empty")
	}

	if a.Actiontype != actionType {
		t.Errorf("Identifier must be equal to: %q", actionType)
	}

	if a.Async != async {
		t.Errorf("Async must be equal to: %v", async)
	}

	if a.Command != "echo \"What a lovely day\" > whatalovelyday" {
		t.Error("Command must be equal to: echo \"What a lovely day\" > whatalovelyday")
	}
}
