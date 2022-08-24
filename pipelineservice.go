package pipelinerunner

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"text/template"
)

const pipelinesTemplate = "PipelineID\tName\n{{range .}}{{.ID}}\t{{.Name}}\n{{end}}"
const pipelineTemplate = "PipelineID: {{.ID}}\nName: {{.Name}}\nActions:\n"

func NewPipelineService() PipelineService {
	actionService = NewActionService()
	databaseService = NewDatabaseService()
	databaseService.InitializeDatabase()
	pipelineRunService = NewPipelineRunService()
	actionExecutionService = NewActionExecutionService()
	return defaultPipelineService{}
}

type defaultPipelineService struct{}

func (defaultPipelineService) PrintPipelines() error {
	pipelines, err := getPipelines()
	if err != nil {
		return fmt.Errorf("error fetching Pipelines data from database: %v", err)
	}
	tmpl := template.Must(template.New("").Parse(pipelinesTemplate))
	w := tabwriter.NewWriter(os.Stdout, 5, 0, 3, ' ', 0)
	if err := tmpl.Execute(w, pipelines); err != nil {
		return fmt.Errorf("error printing Pipelines data")
	}
	w.Flush()
	return nil
}

func (defaultPipelineService) PrintPipeline(id string) error {
	pipeline, err := getPipeline(id)
	if err != nil {
		return fmt.Errorf("error fetching Pipeline data from database: %v", err)
	}
	tmpl := template.Must(template.New("").Parse(pipelineTemplate))
	if err := tmpl.Execute(os.Stdout, pipeline); err != nil {
		return fmt.Errorf("error printing Pipeline data")
	}

	tmpl = template.Must(template.New("").Parse(actionsTemplate))
	w := tabwriter.NewWriter(os.Stdout, 5, 0, 3, ' ', 0)
	if err := tmpl.Execute(w, &pipeline.Actions); err != nil {
		return fmt.Errorf("error printing Pipeline data")
	}
	w.Flush()
	return nil
}

// Create a new Pipeline from Pfile. Add Pipeline data to the database. Return Pipeline id
func (defaultPipelineService) CreatePipeline(pfilePath string, name string) (string, error) {

	pipelineId := getUUID()

	// Insert Pipeline data into the database
	pipelineSql := fmt.Sprintf("INSERT INTO pipeline VALUES ('%s', '%s', CURRENT_TIMESTAMP);", pipelineId, name)
	if err := databaseService.InsertIntoDatabase(pipelineSql); err != nil {
		return "", fmt.Errorf("error inserting Pipeline data into the database: %v", err)
	}

	// Open Pfile
	pfile, err := os.OpenFile(pfilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("pfile opening error: %v", err)
	}
	defer pfile.Close()

	// Parse Pfile
	pfileScanner := bufio.NewScanner(pfile)
	var fullLine string
	for pfileScanner.Scan() {
		line := pfileScanner.Text()
		fullLine += strings.TrimSpace(line)
		// Skip empty lines
		if len(line) == 0 {
			continue
		}
		// Concatenate lines ending by "\"
		if strings.HasSuffix(fullLine, "\\") {
			fullLine = fullLine[:len(fullLine)-1]
			continue
		}

		// Create an Action from a Pfile line
		action, err := actionService.CreateAction(fullLine)
		if err != nil {
			return "", fmt.Errorf("error creating an Action from string: %v", err)
		}
		fullLine = ""
		// Add Action to Pipeline
		if err := addActionToPipeline(action, pipelineId); err != nil {
			return "", fmt.Errorf("error adding Action to Pipeline")
		}
	}
	return pipelineId, nil
}

// Run Pipeline
func (defaultPipelineService) RunPipeline(pipelineId string) error {
	// Get Pipeline data
	var (
		p   Pipeline
		err error
	)
	p, err = getPipelineData(pipelineId)
	if err != nil {
		return fmt.Errorf("error running Pipeline: %v", err)
	}
	pipelineRunId := getUUID()
	if err := pipelineRunService.CreatePipelineRun(pipelineRunId, pipelineId); err != nil {
		return fmt.Errorf("error running Pipeline: %v", err)
	}

	var w sync.WaitGroup
	var errs []error
	errChan := make(chan error)
	go func() {
		for err := range errChan {
			errs = append(errs, err)
		}
	}()

	for _, a := range p.Actions {

		w.Add(1)
		go func(a *Action) {
			var aOut bytes.Buffer
			var aErr bytes.Buffer
			actionExecutionId, err := actionExecutionService.CreateActionExecution(a.ID, pipelineRunId)
			if err != nil {
				errChan <- err
				w.Done()
			}
			if err := actionService.ExecuteAction(a, &aOut, &aErr); err != nil {
				actionExecutionService.FinalizeActionExecution(actionExecutionId, FAILURE, strings.TrimSpace(aOut.String()), strings.TrimSpace(aErr.String()))
				errChan <- err
				w.Done()
			}
			if err := actionExecutionService.FinalizeActionExecution(actionExecutionId, SUCCESS, strings.TrimSpace(aOut.String()), strings.TrimSpace(aErr.String())); err != nil {
				errChan <- err
				w.Done()
			}
			w.Done()
		}(a)

		if !a.Async {
			w.Wait() // Non-async Action waits until all previous Actions finish executing
		}
	}
	w.Wait()

	close(errChan)
	if len(errs) > 0 {
		pipelineRunService.FinalizePipelineRun(pipelineRunId, FAILURE)
		return errs[0] // Return the first error
	} else {
		pipelineRunService.FinalizePipelineRun(pipelineRunId, SUCCESS)
	}

	return nil
}

// Delete Pipeline data from the Database
func (defaultPipelineService) DeletePipeline(id string) error {
	deleteSql := fmt.Sprintf("DELETE FROM pipeline WHERE id = '%s';", id)
	if err := databaseService.InsertIntoDatabase(deleteSql); err != nil {
		return fmt.Errorf("error deleting Pipeline data from the database: %v", err)
	}
	return nil
}

// Select Pipelines data from the database
func getPipelines() ([]Pipeline, error) {
	p := []Pipeline{}
	pSql := "SELECT * FROM pipeline;"
	if err := databaseService.SelectFromDatabase(pSql, &p); err != nil {
		return p, fmt.Errorf("error fetching Pipeline data from the database: %v", err)
	}
	return p, nil
}

// Get Pipeline data
func getPipeline(id string) (Pipeline, error) {
	var (
		p   Pipeline
		err error
	)
	p, err = getPipelineData(id)
	if err != nil {
		return p, fmt.Errorf("error fetching Pipeline data from the database: %v", err)
	}
	return p, nil
}

// Insert Action data into the database. Add Action to Pipeline
func addActionToPipeline(a *Action, pipelineId string) error {
	actionSql := fmt.Sprintf(
		"INSERT INTO action VALUES ('%s', '%s', '%d', '%s', '%s', CURRENT_TIMESTAMP);",
		a.ID,
		pipelineId,
		a.Actiontype,
		strings.ToUpper(strconv.FormatBool(a.Async)),
		a.Command,
	)
	if err := databaseService.InsertIntoDatabase(actionSql); err != nil {
		return fmt.Errorf("error inserting Action data into the database: %v", err)
	}
	return nil
}

func getPipelineData(id string) (Pipeline, error) {
	var p Pipeline
	pSql := fmt.Sprintf("SELECT * FROM pipeline WHERE id = '%s';", id)
	if err := databaseService.SelectFromDatabase(pSql, &p); err != nil {
		return p, fmt.Errorf("error fetching Pipeline data from the database: %v", err)
	}

	aSql := fmt.Sprintf("SELECT * FROM action WHERE pipelineref = '%s';", p.ID)
	if err := databaseService.SelectFromDatabase(aSql, &p.Actions); err != nil {
		return p, fmt.Errorf("error fetching Actions data from the database: %v", err)
	}

	return p, nil
}
