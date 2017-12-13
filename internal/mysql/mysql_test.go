package mysql_test

import (
	"database/sql"
	"flag"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/mysql"
	_ "github.com/mattes/migrate/source/file"
)

var (
	scripts    = flag.String("scripts", "file://migrations", "The location of migration scripts.")
	dbUser     = flag.String("db-user", "jobexec", "Database username")
	dbPassword = flag.String("db-password", "jobexecsecret", "Database password")
	dbAddress  = flag.String("db-address", "localhost:3306", "Database address")
	dbName     = flag.String("db-name", "jobexec_test", "Database name")
)

const driverName = "mysql"

type Suite struct {
	T  *testing.T
	DB *sql.DB
}

func (s *Suite) TearDown() {
	if dbErr := s.DB.Close(); dbErr != nil {
		s.T.Error("failed closing db:", dbErr)
	}
}

func Setup(t *testing.T) *Suite {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?multiStatements=true&clientFoundRows=true&parseTime=true&loc=Local", *dbUser, *dbPassword, *dbAddress, *dbName)
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		t.Fatal("err:", err)
	}

	if err = db.Ping(); err != nil {
		t.Fatal("err:", err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		t.Fatal("err:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(*scripts, driverName, driver)
	if err != nil {
		t.Fatal("err:", err)
	}

	if err := m.Down(); err != nil {
		if err != migrate.ErrNoChange {
			t.Error("Failed execute migration down scripts:", err)
		}
	}

	if err := m.Drop(); err != nil {
		t.Error("Failed execute migration pre-drop:", err)
	}

	if err := m.Up(); err != nil {
		t.Error("Failed execute migration up scripts:", err)
	}

	return &Suite{
		T:  t,
		DB: db,
	}
}
