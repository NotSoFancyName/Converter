package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

const (
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
	maxHeaderBytes = 1 << 20

	apiURL = "/api/v1/currency/convert"
)

type Server struct {
	server *http.Server
	cc     *grpc.ClientConn
	stop   chan struct{}
}

func NewServer(stop chan struct{}, listenPort int, fetcherAddress string) (*Server, error) {
	conn, err := grpc.Dial(fetcherAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to fetcher service: %v", err)
	}

	s := &Server{
		server: &http.Server{
			Addr:           ":" + strconv.FormatInt(int64(listenPort), 10),
			ReadTimeout:    readTimeout,
			WriteTimeout:   writeTimeout,
			MaxHeaderBytes: maxHeaderBytes,
		},
		cc:   conn,
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
		if err := s.cc.Close(); err != nil {
			log.Printf("Failed to close gRPC connection: %v", err)
		}
		log.Println("Server is shut")
		s.stop <- struct{}{}
	}()
	errs <- s.server.ListenAndServe()
}
