package pipelinerunner

import (
	"fmt"
	"testing"
)

func TestCreatePipeline(t *testing.T) {
	setupTestEnv(t, "pipelineservice")
	defer teardownTestEnv(t)

	testPipelineService := NewPipelineService()

	var tests = []struct {
		pipelineName string
		pipelineFile string
		want         []Action
	}{
		{
			"Example1",
			"./pfiles/example1.pfile",
			[]Action{
				{Actiontype: RUN, Async: false, Command: "echo \"What a lovely day\" > whatalovelyday"},
			},
		},
		{
			"Example2",
			"./pfiles/example2.pfile",
			[]Action{

				{Actiontype: RUN, Async: false, Command: "touch step1"},
				{Actiontype: RUN, Async: true, Command: "touch step2; sleep 2; echo \"step2 async, after sleeping for 2s\" > step2"},
				{Actiontype: RUN, Async: true, Command: "echo \"Parallel running action is sleeping. It's current output:\n$(cat sleep2)\" > step3"},
				{Actiontype: RUN, Async: false, Command: "echo \"Two previous parallel actions are completed. Final output of the job which was sleeping:\n$(cat sleep2)\" > step4"},
			},
		},
		{
			"Example3",
			"./pfiles/example3.pfile",
			[]Action{

				{Actiontype: GIT, Async: true, Command: "git clone -b main https://github.com/google/re2.git /tmp"},
				{Actiontype: RUN, Async: false, Command: "echo hello"},
			},
		},
	}

	for _, tt := range tests {
		testname := tt.pipelineName
		t.Run(testname, func(t *testing.T) {
			testPipelineService.CreatePipeline(tt.pipelineFile, tt.pipelineName)
			p := Pipeline{}
			pSql := fmt.Sprintf("SELECT * FROM pipeline WHERE name = '%s';", tt.pipelineName)
			if err := databaseService.SelectFromDatabase(pSql, &p); err != nil {
				t.Errorf("error fetching Pipeline data from the database: %v", err)
			}

			if p.Name != tt.pipelineName {
				t.Errorf("incorrect Pipeline name: got: %s, want: %s", p.Name, tt.pipelineName)
			}

			aa := []Action{}
			aSql := fmt.Sprintf("SELECT * FROM action WHERE pipelineref = '%s' ORDER BY modified;", p.ID)
			if err := databaseService.SelectFromDatabase(aSql, &aa); err != nil {
				t.Errorf("error fetching Action data from the database: %v", err)
			}

			for i, a := range aa {
				if !a.Equal(tt.want[i]) {
					t.Errorf("\ngot: %v, \nwant: %v", aa, tt.want)
				}
			}

		})
	}
}

func TestGetPipeline(t *testing.T) {
	setupTestEnv(t, "pipelineservice")
	defer teardownTestEnv(t)

	pipelineName := "Example1"

	pipeline := NewPipeline(pipelineName, []*Action{})

	databaseService = NewDatabaseService()
	databaseService.InitializeDatabase()

	// Insert Pipeline data into the database
	pipelineSql := fmt.Sprintf("INSERT INTO pipeline VALUES ('%s', '%s', CURRENT_TIMESTAMP);", pipeline.ID, pipeline.Name)
	if err := databaseService.InsertIntoDatabase(pipelineSql); err != nil {
		t.Errorf("error inserting Pipeline data into the database: %v", err)
	}

	// Get Pipeline data using PipelineService
	pipelineDb, err := getPipeline(pipeline.ID)
	if err != nil {
		t.Errorf("error getting Pipeline data from the database: %v", err)
	}

	if pipelineDb.Name != pipeline.Name {
		t.Errorf("incorrect Pipeline name: got: %s, want: %s", pipelineDb.Name, pipeline.Name)
	}
}

func TestDeletePipeline(t *testing.T) {
	setupTestEnv(t, "pipelineservice")
	defer teardownTestEnv(t)

	testPipelineService := NewPipelineService()

	pipelineName := "Example1"

	pipeline := NewPipeline(pipelineName, []*Action{})

	// Insert Pipeline data into the database
	pipelineSql := fmt.Sprintf("INSERT INTO pipeline VALUES ('%s', '%s', CURRENT_TIMESTAMP);", pipeline.ID, pipeline.Name)
	if err := databaseService.InsertIntoDatabase(pipelineSql); err != nil {
		t.Errorf("error inserting Pipeline data into the database: %v", err)
	}

	// Delete Pipeline data using PipelineService
	if err := testPipelineService.DeletePipeline(pipeline.ID); err != nil {
		t.Errorf("error deleting Pipeline data from the database: %v", err)
	}

	var pNum int
	pSql := fmt.Sprintf("SELECT COUNT(*) FROM pipeline WHERE name = '%s';", pipelineName)
	if err := databaseService.SelectFromDatabase(pSql, &pNum); err != nil {
		t.Errorf("error fetching Pipeline data from the database: %v", err)
	}

	if pNum != 0 {
		t.Errorf("Pipeline data should have beeen deleted from the database")
	}

	var aNum int
	aSql := fmt.Sprintf("SELECT COUNT(*) FROM action WHERE pipelineref = '%s';", pipeline.ID)
	if err := databaseService.SelectFromDatabase(aSql, &aNum); err != nil {
		t.Errorf("error fetching Action data from the database: %v", err)
	}

	if aNum != 0 {
		t.Errorf("actions data should have beeen deleted from the database")
	}
}

func TestRunPipeline(t *testing.T) {
	setupTestEnv(t, "pipelineservice")
	defer teardownTestEnv(t)

	testPipelineService := NewPipelineService()
	pipelineId, err := testPipelineService.CreatePipeline("./pfiles/example0.pfile", "Example0")
	if err != nil {
		t.Errorf("error creating Pipeline")
	}
	testPipelineService.RunPipeline(pipelineId)

	pr := PipelineRun{}
	prSql := fmt.Sprintf("SELECT * FROM pipeline_run WHERE pipelineref = '%s';", pipelineId)
	if err := databaseService.SelectFromDatabase(prSql, &pr); err != nil {
		t.Errorf("error fetching PipelineRun data from the database: %v", err)
	}

	if pr == (PipelineRun{}) || pr.Status == FAILURE {
		t.Errorf("incorrect PipelineRun data: %q", pr)
	}

	ae := ActionExecution{}
	aeSql := fmt.Sprintf("SELECT * FROM action_execution WHERE pipelinerunref = '%s';", pr.ID)
	if err := databaseService.SelectFromDatabase(aeSql, &ae); err != nil {
		t.Errorf("error fetching ActionExecution data from the database: %v", err)
	}

	if ae == (ActionExecution{}) || ae.Status == FAILURE || ae.Stdout != "hello world" {
		t.Errorf("incorrect ActionExecution data: %q", ae)
	}
}
