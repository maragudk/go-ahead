// Package server holds the main Server that serves HTTP requests to clients,
// as well as the code that constructs the routes.
package server

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/coreos/go-systemd/daemon"
	"github.com/go-chi/chi"

	"github.com/maragudk/go-ahead/comms"
	"github.com/maragudk/go-ahead/errors2"
	"github.com/maragudk/go-ahead/storage"
)

// Server takes requests and responds. 😎
type Server struct {
	Emailer *comms.Emailer
	Storer  *storage.Storer
	address string
	mux     *chi.Mux
	log     *log.Logger
	cert    string
	key     string
	sm      *scs.SessionManager
}

// Options for New.
type Options struct {
	Emailer *comms.Emailer
	Storer  *storage.Storer
	Logger  *log.Logger
	Host    string
	Port    int
	Cert    string
	Key     string
}

// New creates a new Server.
func New(options Options) *Server {
	logger := options.Logger
	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}
	return &Server{
		Emailer: options.Emailer,
		Storer:  options.Storer,
		address: net.JoinHostPort(options.Host, strconv.Itoa(options.Port)),
		mux:     chi.NewRouter(),
		log:     logger,
		cert:    options.Cert,
		key:     options.Key,
		sm:      scs.New(),
	}
}

// Start the server, set up listening for HTTP, and notify systemd of readiness.
func (s *Server) Start() error {

	if err := s.Storer.Connect(); err != nil {
		return err
	}

	s.sm.Store = postgresstore.NewWithCleanupInterval(s.Storer.DB.DB, 0)
	s.sm.Cookie.Secure = true

	s.setupRoutes()

	errs := make(chan error)
	go func() {
		if err := s.listenAndServe(); err != nil {
			errs <- err
		}
	}()

	if _, err := daemon.SdNotify(false, daemon.SdNotifyReady); err != nil {
		return errors2.Wrap(err, "could not notify daemon of readiness")
	}

	err := <-errs
	return errors2.Wrap(err, "could not listen and serve")
}

// listenAndServe for HTTP.
// Note that all routes should be defined on mux before calling this.
func (s *Server) listenAndServe() error {
	server := http.Server{
		Addr:              s.address,
		Handler:           s.mux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		ErrorLog:          s.log,
	}

	if s.key != "" && s.cert != "" {
		s.log.Printf("Listening for HTTPS on https://%v\n", s.address)
		if err := server.ListenAndServeTLS(s.cert, s.key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return errors2.Wrap(err, "could not start https listener")
		}
		return nil
	}
	s.log.Printf("Listening for HTTP on http://%v\n", s.address)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors2.Wrap(err, "could not start http listener")
	}
	return nil
}
