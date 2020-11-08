// Package storage provides the Storer, which has all the methods to query the underlying storage database.
package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
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
	socket   string
	name     string
	log      *log.Logger
}

// Options for New.
type Options struct {
	User     string
	Password string
	Host     string
	Port     int
	Socket   string
	Name     string
	Logger   *log.Logger
}

// New returns a new Storer with the given options.
// If no logger is provided, the logs are discarded.
func New(options Options) *Storer {
	logger := options.Logger
	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}
	return &Storer{
		user:     options.User,
		password: options.Password,
		host:     options.Host,
		port:     options.Port,
		socket:   options.Socket,
		name:     options.Name,
		log:      logger,
	}
}

// Connect to the database and ping it to test that it works.
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
	if s.socket != "" {
		return fmt.Sprintf("user=%v password=%v host=%v dbname=%v sslmode=disable", s.user, s.password, s.socket, s.name)
	}
	dataSourceName := "postgresql://" + s.user
	if s.password != "" {
		dataSourceName += ":" + s.password
	}
	dataSourceName += "@" + s.host
	if s.port != 0 {
		dataSourceName += ":" + strconv.Itoa(s.port)
	}
	dataSourceName += fmt.Sprintf("/%v?sslmode=disable", s.name)

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
