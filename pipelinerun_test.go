package pipelinerunner

import (
	"testing"
)

func TestNewPipelineRun(t *testing.T) {
	pipelineref, start, end, status := "testPipelineRef", "testStart", "testEnd", SUCCESS
	want := PipelineRun{Pipelineref: pipelineref, Start: start, End: end, Status: status}

	if p := NewPipelineRun(pipelineref, start, end, status); !p.Equal(want) {
		t.Errorf("got: %v, want: %v", p, want)
	}
}
