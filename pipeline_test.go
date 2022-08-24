package pipelinerunner

import (
	"testing"
)

func TestNewPipeline(t *testing.T) {
	a_action_type := RUN
	a_async := false
	a_attr := map[string]string{"Base": "echo \"What a lovely day\" > whatalovelyday"}
	a := NewAction(a_action_type, a_async, a_attr)

	p_name := "Test pipeline"
	p_actions := []*Action{a}
	p := NewPipeline(p_name, p_actions)

	if p.ID == "" {
		t.Error("ID must be non-empty")
	}

	if p.Name != p_name {
		t.Errorf("Name must be equal to: %s", p_name)
	}

	if len(p.Actions) != len(p_actions) || p.Actions[0] != a {
		t.Errorf("Actions must be equal to: %v", p_actions)
	}
}
