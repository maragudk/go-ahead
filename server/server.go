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

	"go-ahead/comms"
	"go-ahead/errors2"
	"go-ahead/storage"
)

// Server takes requests and responds. 😎
type Server struct {
	Emailer         *comms.Emailer
	Storer          *storage.Storer
	externalAddress string
	internalAddress string
	Name            string
	Version         string
	externalMux     *chi.Mux
	internalMux     *chi.Mux
	log             *log.Logger
	cert            string
	key             string
	sm              *scs.SessionManager
}

type Options struct {
	Emailer      *comms.Emailer
	Storer       *storage.Storer
	Logger       *log.Logger
	ExternalHost string
	ExternalPort int
	InternalHost string
	InternalPort int
	Name         string
	Version      string
	Cert         string
	Key          string
}

// New creates a new Server.
func New(options Options) *Server {
	logger := options.Logger
	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}
	return &Server{
		Emailer:         options.Emailer,
		Storer:          options.Storer,
		externalAddress: net.JoinHostPort(options.ExternalHost, strconv.Itoa(options.ExternalPort)),
		internalAddress: net.JoinHostPort(options.InternalHost, strconv.Itoa(options.InternalPort)),
		Name:            options.Name,
		Version:         options.Version,
		externalMux:     chi.NewRouter(),
		internalMux:     chi.NewRouter(),
		log:             logger,
		cert:            options.Cert,
		key:             options.Key,
		sm:              scs.New(),
	}
}

// Start the server, setting up listening for HTTP externally and internally, and notify systemd of readiness.
func (s *Server) Start() error {

	if err := s.Storer.Connect(); err != nil {
		return err
	}

	s.sm.Store = postgresstore.NewWithCleanupInterval(s.Storer.DB.DB, 0)
	s.sm.Cookie.Secure = true

	s.setupInternalRoutes()
	s.setupExternalRoutes()

	errs := make(chan error)
	go func() {
		if err := s.listenAndServeInternal(); err != nil {
			errs <- err
		}
	}()

	go func() {
		if err := s.listenAndServeExternal(); err != nil {
			errs <- err
		}
	}()

	if _, err := daemon.SdNotify(false, daemon.SdNotifyReady); err != nil {
		return errors2.Wrap(err, "could not notify daemon of readiness")
	}

	err := <-errs
	return errors2.Wrap(err, "could not listen and serve")
}

// listenAndServeExternal on the external address. Note that all routes should be defined on externalMux before calling this.
func (s *Server) listenAndServeExternal() error {
	server := http.Server{
		Addr:              s.externalAddress,
		Handler:           s.externalMux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		ErrorLog:          s.log,
	}

	if s.key != "" && s.cert != "" {
		s.log.Printf("Listening for external HTTPS on https://%v\n", s.externalAddress)
		if err := server.ListenAndServeTLS(s.cert, s.key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return errors2.Wrap(err, "could not start external https listener")
		}
		return nil
	}
	s.log.Printf("Listening for external HTTP on http://%v\n", s.externalAddress)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors2.Wrap(err, "could not start external http listener")
	}
	return nil
}

// listenAndServeInternal on the internal address. Note that all routes should be defined on internalMux before calling this.
func (s *Server) listenAndServeInternal() error {
	server := http.Server{
		Addr:              s.internalAddress,
		Handler:           s.internalMux,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		ErrorLog:          s.log,
	}

	if s.key != "" && s.cert != "" {
		s.log.Printf("Listening for internal HTTPS on https://%v\n", s.internalAddress)
		if err := server.ListenAndServeTLS(s.cert, s.key); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return errors2.Wrap(err, "could not start internal http listener")
		}
		return nil
	}
	s.log.Printf("Listening for internal HTTP on http://%v\n", s.internalAddress)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors2.Wrap(err, "could not start internal http listener")
	}
	return nil
}
