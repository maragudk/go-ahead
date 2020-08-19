package storage

import (
	"fmt"
)

// Storer is the storage abstraction.
type Storer struct {
	User         string
	Host         string
	Port         string
	DatabaseName string
	SSLCert      string
	SSLKey       string
	SSLRootCert  string
}

// CreateDataSourceName for connecting with sql.Open. Also used during migrations.
func (s *Storer) CreateDataSourceName(cockroachSchema bool) string {
	var schema string
	if cockroachSchema {
		schema = "cockroachdb"
	} else {
		schema = "postgresql"
	}

	dataSourceName := fmt.Sprintf("%v://%v@%v:%v/%v?", schema, s.User, s.Host, s.Port, s.DatabaseName)
	if s.SSLCert != "" && s.SSLKey != "" && s.SSLRootCert != "" {
		dataSourceName += fmt.Sprintf("sslmode=verify-full&sslcert=%v&sslkey=%v&sslrootcert=%v", s.SSLCert, s.SSLKey, s.SSLRootCert)
	} else {
		dataSourceName += "sslmode=disable"
	}
	return dataSourceName
}
