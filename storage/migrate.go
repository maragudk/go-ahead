package storage

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"

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
	// From https://github.com/golang-migrate/migrate/tree/master/source/go_bindata
	source := bindata.Resource(AssetNames(), func(name string) ([]byte, error) {
		return Asset(name)
	})
	driver, err := bindata.WithInstance(source)
	if err != nil {
		return nil, errors2.Wrap(err, "could not create bindata source driver")
	}
	dataSourceName := s.createDataSourceName()
	m, err := migrate.NewWithSourceInstance("go-bindata", driver, dataSourceName)
	if err != nil {
		return nil, errors2.Wrap(err, "could not create migrator")
	}
	return m, nil
}
