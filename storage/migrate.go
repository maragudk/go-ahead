package storage

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/maragudk/go-ahead/errors2"
)

// MigrateTo a specific version.
// No database connection should be present already, MigrateTo handles connecting.
func (s *Storer) MigrateTo(version uint) error {
	m, err := s.getMigrator()
	if err != nil {
		return err
	}
	return errors2.Wrap(m.Migrate(version), "could not migrate to version %v", version)
}

// MigrateUp all versions.
func (s *Storer) MigrateUp() error {
	m, err := s.getMigrator()
	if err != nil {
		return err
	}
	return errors2.Wrap(m.Up(), "could not migrate up")
}

// MigrateDown all versions.
func (s *Storer) MigrateDown() error {
	m, err := s.getMigrator()
	if err != nil {
		return err
	}
	return errors2.Wrap(m.Down(), "could not migrate down")
}

func (s *Storer) getMigrator() (*migrate.Migrate, error) {
	dataSourceName := s.createDataSourceName()
	m, err := migrate.New("file://migrations", dataSourceName)
	if err != nil {
		return nil, errors2.Wrap(err, "could not create migrator")
	}
	return m, nil
}
