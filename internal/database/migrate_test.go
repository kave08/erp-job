package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"sync"
	"testing"
)

const recordingDriverName = "erp-job-recording-migration-driver"

var registerRecordingDriver sync.Once
var recordingStates sync.Map

func TestApplyMigrationExecutesWithoutTransaction(t *testing.T) {
	t.Parallel()

	state := &recordingState{}
	db := openRecordingDB(t, state)

	if err := applyMigration(context.Background(), db, "000003_delivery_state_entity_keys.up.sql", "SELECT 1"); err != nil {
		t.Fatalf("applyMigration returned error: %v", err)
	}

	if state.beginCalled {
		t.Fatal("expected applyMigration to avoid explicit transactions")
	}
	if len(state.queries) != 2 {
		t.Fatalf("expected 2 statements, got %d (%#v)", len(state.queries), state.queries)
	}
	if state.queries[0] != "SELECT 1" {
		t.Fatalf("unexpected migration query: %q", state.queries[0])
	}
	if !strings.Contains(state.queries[1], "INSERT INTO schema_migrations") {
		t.Fatalf("unexpected schema_migrations query: %q", state.queries[1])
	}
}

func TestApplyMigrationStopsBeforeRecordingVersionOnFailure(t *testing.T) {
	t.Parallel()

	state := &recordingState{failAt: 1}
	db := openRecordingDB(t, state)

	if err := applyMigration(context.Background(), db, "000003_delivery_state_entity_keys.up.sql", "BROKEN SQL"); err == nil {
		t.Fatal("expected migration failure")
	}

	if len(state.queries) != 1 {
		t.Fatalf("expected only migration statement to run, got %d (%#v)", len(state.queries), state.queries)
	}
	if state.beginCalled {
		t.Fatal("expected applyMigration to avoid explicit transactions")
	}
}

func openRecordingDB(t *testing.T, state *recordingState) *sql.DB {
	t.Helper()

	registerRecordingDriver.Do(func() {
		sql.Register(recordingDriverName, &recordingDriver{})
	})

	dsn := t.Name()
	recordingStates.Store(dsn, state)

	db, err := sql.Open(recordingDriverName, dsn)
	if err != nil {
		t.Fatalf("open recording db: %v", err)
	}

	t.Cleanup(func() {
		recordingStates.Delete(dsn)
		_ = db.Close()
	})

	return db
}

type recordingState struct {
	queries     []string
	beginCalled bool
	failAt      int
	execCount   int
}

type recordingDriver struct{}

func (d *recordingDriver) Open(name string) (driver.Conn, error) {
	stateValue, ok := recordingStates.Load(name)
	if !ok {
		return nil, errors.New("missing recording state")
	}

	return &recordingConn{state: stateValue.(*recordingState)}, nil
}

type recordingConn struct {
	state *recordingState
}

func (c *recordingConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare not supported")
}

func (c *recordingConn) Close() error {
	return nil
}

func (c *recordingConn) Begin() (driver.Tx, error) {
	c.state.beginCalled = true
	return &recordingTx{}, nil
}

func (c *recordingConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	c.state.beginCalled = true
	return &recordingTx{}, nil
}

func (c *recordingConn) ExecContext(_ context.Context, query string, _ []driver.NamedValue) (driver.Result, error) {
	c.state.queries = append(c.state.queries, query)
	c.state.execCount++
	if c.state.failAt > 0 && c.state.execCount == c.state.failAt {
		return nil, errors.New("exec failed")
	}

	return driver.RowsAffected(1), nil
}

type recordingTx struct{}

func (t *recordingTx) Commit() error {
	return nil
}

func (t *recordingTx) Rollback() error {
	return nil
}
