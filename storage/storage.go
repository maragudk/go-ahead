// Package storage provides the Storer, which has all the methods to query the underlying storage database.
package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/maragudk/go-ahead/errors2"
)

const (
	connectionTimeout = 10 * time.Second
)

// Storer is the storage abstraction.
type Storer struct {
	DB       *sqlx.DB
	appName  string
	user     string
	host     string
	port     int
	database string
	cert     string
	key      string
	rootCert string
	log      *log.Logger
}

// NewStorerOptions are options for NewStorer.
type NewStorerOptions struct {
	AppName  string
	User     string
	Host     string
	Port     int
	Database string
	Cert     string
	Key      string
	RootCert string
	Logger   *log.Logger
}

// NewStorer returns a new Storer with the given options.
// If no logger is provided, the logs are discarded.
func NewStorer(options NewStorerOptions) *Storer {
	logger := options.Logger
	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}
	return &Storer{
		appName:  options.AppName,
		user:     options.User,
		host:     options.Host,
		port:     options.Port,
		database: options.Database,
		cert:     options.Cert,
		key:      options.Key,
		rootCert: options.RootCert,
		log:      logger,
	}
}

// Connect to the database and ping it to test that it works.
// Also sets the application name.
func (s *Storer) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	dataSourceName := s.createDataSourceName(false)

	s.log.Println("Connecting to db on", dataSourceName)
	var err error
	s.DB, err = sqlx.ConnectContext(ctx, "postgres", dataSourceName)
	if err != nil {
		return errors2.Wrap(err, "could not connect to db")
	}

	if _, err := s.DB.ExecContext(ctx, "set APPLICATION_NAME = $1", s.appName); err != nil {
		return errors2.Wrap(err, "could not set application name")
	}

	return nil
}

// createDataSourceName for connecting with sql.Open. Also used during migrations.
func (s *Storer) createDataSourceName(cockroachSchema bool) string {
	var schema string
	if cockroachSchema {
		schema = "cockroachdb"
	} else {
		schema = "postgresql"
	}

	dataSourceName := fmt.Sprintf("%v://%v@%v:%v/%v?", schema, s.user, s.host, s.port, s.database)
	if s.cert != "" && s.key != "" && s.rootCert != "" {
		dataSourceName += fmt.Sprintf("sslmode=verify-full&sslcert=%v&sslkey=%v&sslrootcert=%v", s.cert, s.key, s.rootCert)
	} else {
		dataSourceName += "sslmode=disable"
	}

	return dataSourceName
}

// Ping the db with the given context and runs a "select 1".
func (s *Storer) Ping(ctx context.Context) error {
	if err := s.DB.PingContext(ctx); err != nil {
		return errors2.Wrap(err, "could not ping")
	}
	if _, err := s.DB.ExecContext(ctx, "select 1"); err != nil {
		return errors2.Wrap(err, "could not select 1")
	}
	return nil
}
