// Package tests contains supporting code for running tests.
package tests

import (
	"context"
	"github.com/google/uuid"
	"github.com/jean-pasqualini/go-service/business/auth"
	"github.com/jean-pasqualini/go-service/business/data/schema"
	"github.com/jean-pasqualini/go-service/foundation/database"
	"github.com/jean-pasqualini/go-service/foundation/web"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"testing"
	"time"
)

// Success and failure markers
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// Configuration for running tests.
var (
	dbImage = "postgres:13-alpine"
	dbPort = "5432"
	envs = []string{"POSTGRES_PASSWORD=postgres"}
	AdminID = "5cf37266-3473-4006-984f-9325122678b7"
	UserID  = "45b5fbd3-755f-4379-8f07-a58d4a30fa2f"
)

// NewUnit creates a test database inside a Docker container. It creates the required table structure but
// the database is otherwise empty. It returns the database to use as well as a a function to call at the end of the test.
func NewUnit(t *testing.T) (*log.Logger, *sqlx.DB, func()) {
	c := startContainer(t, dbImage, dbPort, envs)

	cfg := database.Config{
		User: "postgres",
		Password: "postgres",
		Host: c.Host,
		Name: "postgres",
		DisableTLS: true,
	}
	db, err := database.Open(cfg)
	if err != nil {
		t.Fatalf("opening database connection: %v", err)
	}

	t.Log("waiting for database to be ready ...")

	// Wait for the database to be ready. Wait 100ms longer between each attempt.
	// Do not try more than 20 times.
	var pingError error
	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if pingError != nil {
		dumpContainerLogs(t, c.ID)
		stopContainer(t, c.ID)
		t.Fatalf("database never ready: %v", pingError)
	}

	if err := schema.Migrate(db); err != nil {
		stopContainer(t, c.ID)
		t.Fatalf("migrating error: %s", err)
	}

	// teardown is the function that should be invoked when the called is done with the database
	teardown := func() {
		t.Helper()
		db.Close()
		stopContainer(t, c.ID)
	}

	log := log.New(os.Stdout, "TEST : ", log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

	return log, db, teardown
}

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func IntPointer(i int) *int {
	return &i
}

// Test owns state for running and shutting down tests.
type Test struct {
	TraceID  string
	DB       *sqlx.DB
	Log      *log.Logger
	Auth     *auth.Auth
	KID      string
	Teardown func()

	t *testing.T
}

// Context returns an app level context for testing.
func Context() context.Context {
	values := web.Values{
		TraceId: uuid.New().String(),
		Now: time.Now(),
	}

	return context.WithValue(context.Background(), web.KeyValues, &values)
}