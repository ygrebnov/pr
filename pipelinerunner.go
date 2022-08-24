package pipelinerunner

import (
	"bytes"
)

var (
	databaseService        DatabaseService
	actionService          ActionService
	pipelineService        PipelineService
	pipelineRunService     PipelineRunService
	actionExecutionService ActionExecutionService
)

// DatabaseService represents type capable of initializing a database, inserting data into it, and selecting data from it
type DatabaseService interface {
	InitializeDatabase() error
	InsertIntoDatabase(q string) error
	SelectFromDatabase(q string, result any) error
}

// PipelineService represents type capable of creating, viewing, modifying, running, and deleting pipelines
type PipelineService interface {
	// Print Pipelines data to console
	PrintPipelines() error
	// Print Pipeline data to console
	PrintPipeline(id string) error
	// Create a new Pipeline from Pfile and return Pipeline id
	CreatePipeline(pfile string, name string) (string, error)
	// Run Pipeline
	RunPipeline(id string) error
	// Delete Pipeline data
	DeletePipeline(id string) error
}

// ActionService represents type capable of creating new Actions and executing them
type ActionService interface {
	CreateAction(s string) (*Action, error)
	ExecuteAction(a *Action, aOut *bytes.Buffer, aErr *bytes.Buffer) error
}

// PipelineRunService represents type capable of creating, finalizing, viewing, and deleting PipelineRuns
type PipelineRunService interface {
	CreatePipelineRun(pipelineRunId string, pipelineId string) error
	FinalizePipelineRun(id string, status RunStatus) error
	PrintPipelineRuns(id string) error
	PrintAllPipelineRuns() error
	GetPipelineRun(id string) (PipelineRun, error)
	DeletePipelineRun(id string) error
}

// ActionExecutionService represents type capable of creating, finalizing, viewing, and deleting ActionExecutions
type ActionExecutionService interface {
	// Create an ActionExecution and return it's id
	CreateActionExecution(actionId string, pipelineRunId string) (string, error)
	// Set the given ActionExecution final status
	FinalizeActionExecution(id string, status RunStatus, stdout string, stderr string) error
}
