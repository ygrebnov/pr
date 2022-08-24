package pipelinerunner

import (
	"fmt"
	"testing"
)

func TestCreatePipelineRun(t *testing.T) {
	setupTestEnv(t, "pipelinerunservice")
	defer teardownTestEnv(t)

	testPipelineService := NewPipelineService()
	testPipelineRunService := NewPipelineRunService()

	pipelineRunId := "testPipelineRunId"
	pipelineName := "Example1"

	pipelineId, err := testPipelineService.CreatePipeline("./pfiles/example1.pfile", pipelineName)
	if err != nil {
		t.Errorf("error creating Pipeline: %v", err)
	}
	testPipelineRunService.CreatePipelineRun(pipelineRunId, pipelineId)

	var prNum int
	pSql := fmt.Sprintf("SELECT COUNT(*) FROM pipeline_run WHERE pipelineref = '%s';", pipelineId)
	if err := databaseService.SelectFromDatabase(pSql, &prNum); err != nil {
		t.Errorf("error fetching PipelineRun data from the database: %v", err)
	}

	if prNum != 1 {
		t.Errorf("there must be one PipelineRun in the database for '%s' Pipeline", pipelineName)
	}

}
