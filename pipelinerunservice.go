package pipelinerunner

import (
	"fmt"
	"os"
	"text/tabwriter"
	"text/template"
)

func NewPipelineRunService() PipelineRunService {
	databaseService = NewDatabaseService()
	databaseService.InitializeDatabase()
	return defaultPipelineRunService{}
}

type defaultPipelineRunService struct{}

const allPipelineRunsTemplate = "PipelineRunID\tStatus\tStart\tEnd\tPipelineID\n{{range .}}{{.ID}}\t{{.Status}}\t{{.Start}}\t{{.End}}\t{{.Pipelineref}}\n{{end}}"
const pipelineRunsTemplate = "PipelineRunID\tStatus\tStart\tEnd\n{{range .}}{{.ID}}\t{{.Status}}\t{{.Start}}\t{{.End}}\n{{end}}"

func (defaultPipelineRunService) PrintAllPipelineRuns() error {
	pipelineRuns, err := getAllPipelineRuns()
	if err != nil {
		return fmt.Errorf("error fetching Pipelines data from database: %v", err)
	}
	tmpl := template.Must(template.New("").Parse(allPipelineRunsTemplate))
	w := tabwriter.NewWriter(os.Stdout, 5, 0, 3, ' ', 0)
	if err := tmpl.Execute(w, pipelineRuns); err != nil {
		return fmt.Errorf("error printing PipelineRuns data")
	}
	w.Flush()
	return nil
}

func (defaultPipelineRunService) PrintPipelineRuns(id string) error {
	pipelineRuns, err := getPipelineRuns(id)
	if err != nil {
		return fmt.Errorf("error fetching Pipelines data from database: %v", err)
	}
	tmpl := template.Must(template.New("").Parse(pipelineRunsTemplate))
	w := tabwriter.NewWriter(os.Stdout, 5, 0, 3, ' ', 0)
	if err := tmpl.Execute(w, pipelineRuns); err != nil {
		return fmt.Errorf("error printing PipelineRuns data")
	}
	w.Flush()
	return nil
}

func (defaultPipelineRunService) CreatePipelineRun(pipelineRunId string, pipelineId string) error {
	prSql := fmt.Sprintf(
		"INSERT INTO pipeline_run VALUES ('%s', '%s', CURRENT_TIMESTAMP, NULL, '%d');",
		pipelineRunId,
		pipelineId,
		RUNNING,
	)
	if err := databaseService.InsertIntoDatabase(prSql); err != nil {
		return fmt.Errorf("error creating a PipelineRun record in the database: %v", err)
	}
	return nil
}

func (defaultPipelineRunService) FinalizePipelineRun(id string, status RunStatus) error {
	prSql := fmt.Sprintf(
		"UPDATE pipeline_run SET end = CURRENT_TIMESTAMP, status = '%d' WHERE id = '%s';",
		status,
		id,
	)
	if err := databaseService.InsertIntoDatabase(prSql); err != nil {
		return fmt.Errorf("error updating PipelineRun record in the database: %v", err)
	}
	return nil
}

func (defaultPipelineRunService) GetPipelineRun(id string) (PipelineRun, error) {
	p := PipelineRun{}
	prSql := fmt.Sprintf("SELECT * FROM pipeline_run WHERE id = '%s';", id)
	if err := databaseService.SelectFromDatabase(prSql, &p); err != nil {
		return p, fmt.Errorf("error fetching PipelineRun data from the database: %v", err)
	}
	return p, nil
}

func (defaultPipelineRunService) DeletePipelineRun(id string) error {
	deleteSql := fmt.Sprintf("DELETE FROM pipeline_run WHERE id = '%s';", id)
	if err := databaseService.InsertIntoDatabase(deleteSql); err != nil {
		return fmt.Errorf("error deleting PipelineRun data from the database: %v", err)
	}
	return nil
}

func getAllPipelineRuns() ([]PipelineRun, error) {
	pr := []PipelineRun{}
	prSql := "SELECT * FROM pipeline_run ORDER BY start DESC;"
	if err := databaseService.SelectFromDatabase(prSql, &pr); err != nil {
		return pr, fmt.Errorf("error fetching PipelineRun data from the database: %v", err)
	}
	return pr, nil
}

func getPipelineRuns(id string) ([]PipelineRun, error) {
	pr := []PipelineRun{}
	prSql := fmt.Sprintf("SELECT * FROM pipeline_run WHERE pipelineref = '%s' ORDER BY start DESC;", id)
	if err := databaseService.SelectFromDatabase(prSql, &pr); err != nil {
		return pr, fmt.Errorf("error fetching PipelineRun data from the database: %v", err)
	}
	return pr, nil
}
