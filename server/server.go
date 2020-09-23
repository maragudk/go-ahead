package server

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"

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
	}
}

// Start the server, setting up listening for HTTP externally and internally, and notify systemd of readiness.
func (s *Server) Start() error {

	if err := s.Storer.Connect(); err != nil {
		return err
	}

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
	s.log.Printf("Listening for external HTTP on http://%v\n", s.externalAddress)
	if err := http.ListenAndServe(s.externalAddress, s.externalMux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors2.Wrap(err, "could not start external http listener")
	}
	return nil
}

// listenAndServeInternal on the internal address. Note that all routes should be defined on internalMux before calling this.
func (s *Server) listenAndServeInternal() error {
	s.log.Printf("Listening for internal HTTP on http://%v\n", s.internalAddress)
	if err := http.ListenAndServe(s.internalAddress, s.internalMux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors2.Wrap(err, "could not start internal http listener")
	}
	return nil
}
