package pipelinerunner

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const schema string = `
DROP TABLE IF EXISTS pipeline;

CREATE TABLE pipeline (
    id TEXT PRIMARY KEY, 
    name TEXT NOT NULL UNIQUE, 
    modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS action;

CREATE TABLE action (
    id TEXT PRIMARY KEY,
    pipelineref TEXT NOT NULL, 
    actiontype TEXT NOT NULL,
    async BOOLEAN NOT NULL,
    command TEXT NOT NULL,
    modified TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(pipelineref) REFERENCES pipeline(id) ON DELETE CASCADE
);

DROP TABLE IF EXISTS pipeline_run;

CREATE TABLE pipeline_run (
    id TEXT PRIMARY KEY,
    pipelineref TEXT NOT NULL,
    start TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end TIMESTAMP,
    status TEXT NOT NULL,
    FOREIGN KEY(pipelineref) REFERENCES pipeline(id) ON DELETE CASCADE
);


DROP TABLE IF EXISTS action_execution;

CREATE TABLE action_execution (
    id TEXT PRIMARY KEY,
    actionref TEXT NOT NULL,
    pipelinerunref TEXT NOT NULL,
    start TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end TIMESTAMP,
    stdout TEXT,
    stderr TEXT,
    status TEXT NOT NULL,
    FOREIGN KEY(actionref) REFERENCES action(id) ON DELETE CASCADE
    FOREIGN KEY(pipelinerunref) REFERENCES pipeline_run(id) ON DELETE CASCADE
);
`

var db *sqlx.DB

type defaultDatabaseService struct{}

func NewDatabaseService() DatabaseService {
	return defaultDatabaseService{}
}

func (defaultDatabaseService) InitializeDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbcxn, err := sqlx.Connect("sqlite", fmt.Sprintf("%s/%s", Config.WorkDir, Config.Database.FileName))
	if err != nil {
		return fmt.Errorf("cannot connect to database: %s/%s", Config.WorkDir, Config.Database.FileName)
	}
	db = dbcxn

	var tables_exist int
	query := "SELECT COUNT(*) as tables_exist FROM sqlite_master WHERE type='table' AND name IN ('pipeline', 'action');"
	if err := db.GetContext(ctx, &tables_exist, query); err != nil {
		return fmt.Errorf("error selecting data from the database: %v", err)
	}

	if tables_exist == 0 {
		if _, err := db.ExecContext(ctx, schema); err != nil {
			return fmt.Errorf("error initializing database: %v", err)
		}
	}
	return nil
}

func (defaultDatabaseService) InsertIntoDatabase(q string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := db.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("error inserting data into the database: %v", err)
	}
	return nil
}

func (defaultDatabaseService) SelectFromDatabase(q string, result any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	if strings.Contains(fmt.Sprintf("%T", result), "[]") {
		err = db.SelectContext(ctx, result, q)
	} else {
		err = db.GetContext(ctx, result, q)
	}
	if err != nil {
		return fmt.Errorf("error selecting data from the database: %v", err)
	}
	return nil
}
