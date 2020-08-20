package ahead

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-systemd/daemon"

	"go-ahead/errors2"
	"go-ahead/storage"
)

// Server takes requests and responds. 😎
type Server struct {
	Storer       *storage.Storer
	ExternalPort int
	InternalPort int
	Version      string
	externalMux  *http.ServeMux
	internalMux  *http.ServeMux
}

type NewServerOptions struct {
	Storer       *storage.Storer
	ExternalPort int
	InternalPort int
	Version      string
}

func NewServer(options NewServerOptions) *Server {
	return &Server{
		Storer:       options.Storer,
		ExternalPort: options.ExternalPort,
		InternalPort: options.InternalPort,
		externalMux:  http.NewServeMux(),
		internalMux:  http.NewServeMux(),
		Version:      options.Version,
	}
}

// Start the server, setting up listening for HTTP externally and internally, and notify systemd of readiness.
func (s *Server) Start() error {
	hostname := getLocalhostOrEmpty()
	errs := make(chan error)

	if err := s.Storer.Connect(); err != nil {
		return err
	}

	s.setupRoutes()

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

func (s *Server) setupRoutes() {
	s.setupHealthHandler()
	s.setupVersionHandler()
	s.setupMetricsHandler()
}

// listenAndServeExternal on the external port. Note that all routes should be defined on externalMux before calling this.
func (s *Server) listenAndServeExternal(hostname string) error {
	addr := fmt.Sprintf("%v:%v", hostname, s.ExternalPort)
	log.Println("Listening for external HTTP on", addr)
	if err := http.ListenAndServe(addr, s.externalMux); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors2.Wrap(err, "could not start external http listener")
	}
	return nil
}

// listenAndServeInternal on the internal port. Note that all routes should be defined on internalMux before calling this.
func (s *Server) listenAndServeInternal(hostname string) error {
	addr := fmt.Sprintf("%v:%v", hostname, s.InternalPort)
	log.Println("Listening for internal HTTP on", addr)
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
