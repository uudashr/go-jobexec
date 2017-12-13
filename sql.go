package jobexec

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jonboulle/clockwork"
)

// SQLContext is sql based context.
type SQLContext struct {
	db       *sql.DB
	interval time.Duration
	name     string // context name

	clock clockwork.Clock
}

// Init will create tracker entries on database if not exists yet.
func (c *SQLContext) Init() error {
	exists, err := c.trackerExists()
	if err != nil {
		return err
	}

	if !exists {
		return c.createTracker()
	}

	return nil
}

func (c *SQLContext) trackerExists() (bool, error) {
	var trackerCount int
	sqlCount := "SELECT COUNT(*) FROM execution_trackers WHERE name = ?"
	if err := c.db.QueryRow(sqlCount, c.name).Scan(&trackerCount); err != nil {
		return false, err
	}

	return trackerCount == 1, nil
}

func (c *SQLContext) createTracker() error {
	sqlInsert := "INSERT INTO execution_trackers (name, last_executor) VALUES (?, ?)"
	res, err := c.db.Exec(sqlInsert, c.name, "")
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("jobexec: fail to create execution tracker")
	}

	return nil
}

// Acquire an execution permission.
func (c *SQLContext) Acquire(execID ExecutorID) error {
	now := c.clock.Now()
	nextExecTime := now.Add(c.interval)
	sqlUpdate := "UPDATE execution_trackers SET last_executor = ?, next_execution_time = ? WHERE name = ? AND (last_executor IS NULL OR last_executor = ? OR next_execution_time < ?)"
	res, err := c.db.Exec(sqlUpdate, execID, nextExecTime, c.name, execID, now)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return ErrNotPermitted
	}

	return nil
}

// LastExecutor permitted to do an execution.
func (c *SQLContext) LastExecutor() (ExecutorID, error) {
	var lastExecutor ExecutorID
	if err := c.db.QueryRow("SELECT last_executor FROM execution_trackers WHERE name = ?", c.name).Scan(&lastExecutor); err != nil {
		return EmptyExecutorID, err
	}

	return lastExecutor, nil
}

// NewSQLContext constructs new SQLContext.
func NewSQLContext(db *sql.DB, name string, interval time.Duration) (*SQLContext, error) {
	return NewSQLContextWithClock(db, name, interval, clockwork.NewRealClock())
}

// NewSQLContextWithClock construcrs new SQLContext with clock.
func NewSQLContextWithClock(db *sql.DB, name string, interval time.Duration, clock clockwork.Clock) (*SQLContext, error) {
	if db == nil {
		return nil, errors.New("jobcontext: nil db")
	}

	if name == "" {
		return nil, errors.New("jobcontext: empty name")
	}

	if interval == 0 {
		return nil, errors.New("jobcontext: 0 interval")
	}

	if clock == nil {
		return nil, errors.New("jobcontext: nil clock")
	}

	return &SQLContext{
		db:       db,
		name:     name,
		interval: interval,
		clock:    clock,
	}, nil
}
