package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/NotSoFancyName/conversion_service/service/currency_manager"
)

const (
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
	maxHeaderBytes = 1 << 20

	apiURL = "/api/v1/currency/convert"
)

type Server struct {
	server *http.Server
	cm     currency_manager.Manager
	stop   chan struct{}
}

func NewServer(stop chan struct{}, port int) (*Server, error) {
	cm, err := currency_manager.NewManagerOfType(currency_manager.PostgresManager)
	if err != nil {
		return nil, fmt.Errorf("failed to create currency manager: %v", err)
	}
	s := &Server{
		server: &http.Server{
			Addr:           ":" + strconv.FormatInt(int64(port), 10),
			ReadTimeout:    readTimeout,
			WriteTimeout:   writeTimeout,
			MaxHeaderBytes: maxHeaderBytes,
		},
		cm:   cm,
		stop: stop,
	}
	s.server.Handler = http.HandlerFunc(s.handleGetCurrenciesExchangeRate)
	return s, nil
}

func (s *Server) Run(errs chan<- error) {
	log.Printf("Running server on port %v", s.server.Addr)
	go func() {
		<-s.stop
		if err := s.server.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to shutdown server properly: %v", err)
		}
		if err := s.cm.Shutdown(); err != nil {
			log.Printf("Failed to shutdown DB properly: %v", err)
		}
		log.Println("Server is shut")
		s.stop <- struct{}{}
	}()
	errs <- s.server.ListenAndServe()
}
