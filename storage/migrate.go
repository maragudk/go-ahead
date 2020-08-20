package storage

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"

	"go-ahead/errors2"
)

// Migrate database to a specific version.
// No database connection should be present already, Migrate handles connecting.
func (s *Storer) Migrate(version uint) error {
	// From https://github.com/golang-migrate/migrate/tree/master/source/go_bindata
	source := bindata.Resource(AssetNames(), func(name string) ([]byte, error) {
		return Asset(name)
	})
	driver, err := bindata.WithInstance(source)
	if err != nil {
		return errors2.Wrap(err, "could not create bindata source driver")
	}
	dataSourceName := s.createDataSourceName(true)
	log.Println("Connecting on", dataSourceName)
	m, err := migrate.NewWithSourceInstance("go-bindata", driver, dataSourceName)
	if err != nil {
		return errors2.Wrap(err, "could not create migrator")
	}
	return errors2.Wrap(m.Migrate(version), "could not migrate to version %v", version)
}
