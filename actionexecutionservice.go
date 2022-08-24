package pipelinerunner

import (
	"fmt"
)

func NewActionExecutionService() ActionExecutionService {
	databaseService = NewDatabaseService()
	databaseService.InitializeDatabase()
	return defaultActionExecutionService{}
}

type defaultActionExecutionService struct{}

// Create an ActionExecution record in the database. Return ActionExecution id
func (defaultActionExecutionService) CreateActionExecution(actionId string, pipelineRunId string) (string, error) {
	aeId := getUUID()

	aeSql := fmt.Sprintf(
		"INSERT INTO action_execution VALUES ('%s', '%s', '%s', CURRENT_TIMESTAMP, '', '', '', '%d');",
		aeId,
		actionId,
		pipelineRunId,
		RUNNING,
	)
	if err := databaseService.InsertIntoDatabase(aeSql); err != nil {
		return "", fmt.Errorf("error creating an ActionExecution record in the database: %v", err)
	}
	return aeId, nil
}

func (defaultActionExecutionService) FinalizeActionExecution(
	id string,
	status RunStatus,
	stdout string,
	stderr string,
) error {
	aeSql := fmt.Sprintf(
		"UPDATE action_execution SET end = CURRENT_TIMESTAMP, status = '%d', stdout = '%s', stderr = '%s' WHERE id = '%s';",
		status,
		stdout,
		stderr,
		id,
	)
	if err := databaseService.InsertIntoDatabase(aeSql); err != nil {
		return fmt.Errorf("error updating ActionExecution record in the database: %v", err)
	}
	return nil
}
