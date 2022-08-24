package pipelinerunner

import (
	"fmt"

	"github.com/google/uuid"
)

type RunStatus uint8 // To identify ActionExecution and PipelineRun statuses

const (
	RUNNING RunStatus = iota
	SUCCESS
	FAILURE
	INVALIDRUNSTATUS
)

func (s RunStatus) String() string {
	return []string{"RUNNING", "SUCCESS", "FAILURE"}[s]
}

func runStatusFromString(s string) (RunStatus, error) {
	switch s {
	case "RUNNING":
		return RUNNING, nil
	case "SUCCESS":
		return SUCCESS, nil
	case "FAILUTE":
		return FAILURE, nil
	default:
		return INVALIDRUNSTATUS, fmt.Errorf("unknown run status: %s", s)
	}
}

func getUUID() string {
	return uuid.New().String()
}
