package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-systemd/daemon"
	"github.com/go-chi/chi"

	"go-ahead/comms"
	"go-ahead/errors2"
	"go-ahead/storage"
)

// Server takes requests and responds. 😎
type Server struct {
	Emailer      *comms.Emailer
	Storer       *storage.Storer
	ExternalPort int
	InternalPort int
	Name         string
	Version      string
	externalMux  *chi.Mux
	internalMux  *chi.Mux
	log          *log.Logger
}

type Options struct {
	Emailer      *comms.Emailer
	Storer       *storage.Storer
	Logger       *log.Logger
	ExternalPort int
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
		Emailer:      options.Emailer,
		Storer:       options.Storer,
		ExternalPort: options.ExternalPort,
		InternalPort: options.InternalPort,
		Name:         options.Name,
		Version:      options.Version,
		externalMux:  chi.NewRouter(),
		internalMux:  chi.NewRouter(),
		log:          logger,
	}
}

// Start the server, setting up listening for HTTP externally and internally, and notify systemd of readiness.
func (s *Server) Start() error {
	hostname := getLocalhostOrEmpty()
	errs := make(chan error)

	if err := s.Storer.Connect(); err != nil {
		return err
	}

	s.setupInternalRoutes()
	s.setupExternalRoutes()

	go func() {
		if err := s.listenAndServeInternal(hostname); err != nil {
			errs <- err
		}
	}()

	go func() {
		if err := s.listenAndServeExternal(hostname); err != nil {
			errs <- err
		}
	}()

	if _, err := daemon.SdNotify(false, daemon.SdNotifyReady); err != nil {
		return errors2.Wrap(err, "could not notify daemon of readiness")
	}

	err := <-errs
	return errors2.Wrap(err, "could not listen and serve")
}

// listenAndServeExternal on the external port. Note that all routes should be defined on externalMux before calling this.
func (s *Server) listenAndServeExternal(hostname string) error {
	addr := fmt.Sprintf("%v:%v", hostname, s.ExternalPort)
	s.log.Println("Listening for external HTTP on", addr)
	if err := http.ListenAndServe(addr, s.externalMux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors2.Wrap(err, "could not start external http listener")
	}
	return nil
}

// listenAndServeInternal on the internal port. Note that all routes should be defined on internalMux before calling this.
func (s *Server) listenAndServeInternal(hostname string) error {
	addr := fmt.Sprintf("%v:%v", hostname, s.InternalPort)
	s.log.Println("Listening for internal HTTP on", addr)
	if err := http.ListenAndServe(addr, s.internalMux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors2.Wrap(err, "could not start internal http listener")
	}
	return nil
}

// getLocalhostOrEmpty tries to figure out whether we're on a development machine, in which case we listen on localhost only.
// Otherwise, return the empty string.
func getLocalhostOrEmpty() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	if strings.Contains(hostname, ".local") {
		return "localhost"
	}
	return ""
}
