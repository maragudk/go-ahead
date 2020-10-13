// +build integration

// Package storagetest makes integration testing easier.
package storagetest

import (
	"github.com/maragudk/go-ahead/storage"
)

func createStorer(user string) *storage.Storer {
	return storage.NewStorer(storage.NewStorerOptions{
		User:     user,
		Password: "123",
		Host:     "localhost",
		Port:     5432,
		Database: "ahead",
	})
}

// CreateStorer for testing.
// Usage:
// 	s, cleanup := CreateStorer()
// 	defer cleanup()
// 	…
func CreateStorer() (*storage.Storer, func()) {
	rootStorer := setupDB()
	s := createStorer("ahead")
	if err := s.Connect(); err != nil {
		panic(err)
	}

	return s, func() {
		if err := rootStorer.MigrateDown(); err != nil {
			panic(err)
		}
	}
}

// CreateRootStorer is like CreateStorer, but using the root user.
// This is primarily for migrations.
func CreateRootStorer() (*storage.Storer, func()) {
	s := setupDB()
	if err := s.Connect(); err != nil {
		panic(err)
	}

	return s, func() {
		if err := s.MigrateDown(); err != nil {
			panic(err)
		}
	}
}

// setupDB with root privileges.
func setupDB() *storage.Storer {
	s := createStorer("ahead")
	if err := s.Connect(); err != nil {
		panic(err)
	}

	if err := s.MigrateUp(); err != nil {
		panic(err)
	}

	return s
}
