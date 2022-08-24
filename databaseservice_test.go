package pipelinerunner

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestInitializeDatabase(t *testing.T) {
	setupTestEnv(t, "db")
	defer teardownTestEnv(t)

	testDatabaseService := NewDatabaseService()
	testDatabaseService.InitializeDatabase()

	if _, err := os.Stat(fmt.Sprintf("%s/%s", testWorkDir, defaultDbFile)); err != nil {
		t.Error("database file does not exist")
	}

	ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	var tables_exist int
	query := "SELECT COUNT(*) as tables_exist FROM sqlite_master WHERE type='table' AND name IN ('pipeline', 'action');"
	if err := db.GetContext(ctx, &tables_exist, query); err != nil {
		t.Errorf("error selecting data from the database: %v", err)
	}
	if tables_exist == 0 {
		t.Error("'pipeline' and 'action' tables were not created")
	}
}
