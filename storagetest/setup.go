// +build integration

package storagetest

import (
	"bufio"
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"go-ahead/storage"
)

const port = 26258

// HandleTestMain should be called in TestMain like so:
// func TestMain(m *testing.M) {
//   storagetest.HandleTestMain(m)
// }
func HandleTestMain(m *testing.M) {
	s := setupDB()
	code := m.Run()
	defer func(code int) {
		dropDB(s)
		os.Exit(code)
	}(code)
}

func createStorer(user string) *storage.Storer {
	return storage.NewStorer(storage.NewStorerOptions{
		AppName:  "ahead_test",
		User:     user,
		Host:     "localhost",
		Port:     port,
		Database: "ahead",
	})
}

func CreateStorer() (*storage.Storer, func()) {
	s := createStorer("ahead")
	if err := s.Connect(); err != nil {
		panic(err)
	}

	return s, func() {
		// Do cleanup, nothing yet
	}
}

func CreateRootStorer() (*storage.Storer, func()) {
	s := createStorer("root")
	if err := s.Connect(); err != nil {
		panic(err)
	}

	return s, func() {
		// Do cleanup, nothing yet
	}
}

// setupDB with root privileges.
func setupDB() *storage.Storer {
	s := createStorer("root")
	if err := s.Connect(); err != nil {
		panic(err)
	}

	executeSQLFromFile(s.DB.DB, "../storagetest/testdata/drop-database.sql")
	executeSQLFromFile(s.DB.DB, "../storagetest/testdata/create-database.sql")

	if err := s.MigrateUp(); err != nil {
		panic(err)
	}

	return s
}

func dropDB(s *storage.Storer) {
	executeSQLFromFile(s.DB.DB, "../storagetest/testdata/drop-database.sql")
}

func executeSQLFromFile(db *sql.DB, path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := ""
	for scanner.Scan() {
		line := scanner.Text()
		// Skip comments
		if strings.HasPrefix(line, "--") {
			continue
		}
		query += line + " "
		if !strings.HasSuffix(query, "; ") {
			continue
		}
		_, err := db.ExecContext(ctx, query)
		query = ""
		if err != nil {
			panic(err)
		}
	}
}
