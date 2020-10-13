// Package storage provides the Storer, which has all the methods to query the underlying storage database.
package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
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
	user     string
	password string
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
	User     string
	Password string
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
		user:     options.User,
		password: options.Password,
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

	dataSourceName := s.createDataSourceName()

	var err error
	s.DB, err = sqlx.ConnectContext(ctx, "postgres", dataSourceName)
	if err != nil {
		return errors2.Wrap(err, "could not connect to db")
	}

	return nil
}

// createDataSourceName for connecting with sql.Open. Also used during migrations.
func (s *Storer) createDataSourceName() string {
	dataSourceName := "postgresql://" + s.user
	if s.password != "" {
		dataSourceName += ":" + s.password
	}
	dataSourceName += "@" + url.PathEscape(s.host)
	if s.port != 0 {
		dataSourceName += ":" + strconv.Itoa(s.port)
	}
	dataSourceName += fmt.Sprintf("/%v?", s.database)

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
