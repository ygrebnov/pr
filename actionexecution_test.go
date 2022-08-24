package pipelinerunner

import (
	"testing"
)

func TestNewActionExecution(t *testing.T) {
	actionref, pipelinerunref, start, end := "testActionRef", "testPipelineRunRef", "testStart", "testEnd"
	stdout, stderr, status := "testStdout", "testStderr", SUCCESS
	want := ActionExecution{Actionref: actionref, Pipelinerunref: pipelinerunref, Start: start, End: end,
		Stdout: stdout, Stderr: stderr, Status: status}

	if p := NewActionExecution(actionref, pipelinerunref, start, end, stdout, stderr, status); !p.Equal(want) {
		t.Errorf("got: %v, want: %v", p, want)
	}
}
