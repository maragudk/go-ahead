// +build integration

package storagetest

import (
	"bufio"
	"context"
	"database/sql"
	"os"
	"strings"
	"time"

	"go-ahead/storage"
)

const port = 26258

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
	rootStorer := setupDB()
	s := createStorer("ahead")
	if err := s.Connect(); err != nil {
		panic(err)
	}

	return s, func() {
		dropDB(rootStorer)
	}
}

func CreateRootStorer() (*storage.Storer, func()) {
	s := setupDB()
	if err := s.Connect(); err != nil {
		panic(err)
	}

	return s, func() {
		dropDB(s)
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
